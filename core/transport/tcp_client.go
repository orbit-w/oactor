package transport

import (
	"context"
	"fmt"
	"github.com/orbit-w/golib/bases/misc/number_utils"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/golib/bases/wrappers/sender_wrapper"
	"github.com/orbit-w/oactor/core/mmrpcs"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

/*
   @Author: orbit-w
   @File: tcp_client
   @2023 11月 周日 16:32
*/

// TcpClient implements the IConn interface with TCP.
type TcpClient struct {
	mu            sync.Mutex
	state         atomic.Uint32
	lastAck       atomic.Int64
	remoteAddr    string
	remoteNodeId  string
	currentNodeId string

	ctx      context.Context
	cancel   context.CancelFunc
	codec    *TcpCodec
	msgCodec *Codec

	conn    net.Conn
	buf     *ControlBuffer
	sw      *sender_wrapper.SenderWrapper
	r       *receiver
	dHandle func(remoteNodeId string)
}

// DialWithOps Encapsulates asynchronous TCP connection establishment (with retries and backoff)
func DialWithOps(remoteAddr string, _ops ...*DialOption) IConn {
	dp := parseOptions(_ops...)
	ctx, cancel := context.WithCancel(context.Background())
	buf := new(ControlBuffer)
	BuildControlBuffer(buf, dp.MaxIncomingPacket)
	tc := &TcpClient{
		mu:            sync.Mutex{},
		remoteAddr:    remoteAddr,
		remoteNodeId:  dp.RemoteNodeId,
		currentNodeId: dp.CurrentNodeId,
		dHandle:       dp.DisconnectHandler,
		buf:           buf,
		ctx:           ctx,
		cancel:        cancel,
		codec:         NewTcpCodec(dp.MaxIncomingPacket, false),
		r:             newReceiver(),
	}

	go tc.handleDial(dp)
	return tc
}

// Write TcpClient obj does not implicitly call IPacket.Return to return the
// packet to the pool, and the user needs to explicitly call it.
func (tc *TcpClient) Write(out packet.IPacket) error {
	pack := tc.msgCodec.encode(out, TypeMessageRaw)
	defer pack.Return()
	err := tc.buf.Set(pack)
	return err
}

func (tc *TcpClient) Recv() (packet.IPacket, error) {
	return tc.r.read()
}

func (tc *TcpClient) Close() error {
	if tc.conn != nil {
		_ = tc.conn.Close()
	}
	return nil
}

func (tc *TcpClient) handleDial(_ *DialOption) {
	defer func() {
		if tc.dHandle != nil {
			tc.dHandle(tc.remoteNodeId)
		}
		tc.buf.OnClose()
	}()

	task := func() error {
		return tc.dial()
	}

	//When the number of failed connection attempts reaches the upper limit,
	//the conn state will be set to the 'disconnected' state,
	//and all virtual streams will be closed.
	if err := withRetry(task); err != nil {
		tc.mu.Lock()
		defer tc.mu.Unlock()
		tc.state.Store(StatusDisconnected)
		tc.r.onClose(mmrpcs.ErrCanceled)
		return
	}

	defer func() {
		if tc.state.CompareAndSwap(StatusConnected, StatusDisconnected) {
			tc.r.onClose(mmrpcs.ErrCanceled)
		}
	}()

	tc.state.Store(StatusConnected)
	tc.lastAck.Store(0)
	tc.sw = sender_wrapper.NewSender(tc.SendData)
	tc.buf.Run(tc.sw)
	go tc.keepalive()
	<-tc.ctx.Done()
}

func (tc *TcpClient) SendData(data packet.IPacket) error {
	err := tc.sendData(data)
	if err != nil {
		log.Println("[TcpClient] [func: SendData] exec failed: ", err.Error())
		if tc.conn != nil {
			_ = tc.conn.Close()
		}
	}
	return err
}

func (tc *TcpClient) sendData(data packet.IPacket) error {
	body := tc.codec.EncodeBody(data)
	if err := tc.conn.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
		body.Return()
		return err
	}
	_, err := tc.conn.Write(body.Data())
	body.Return()
	return err
}

func (tc *TcpClient) dial() error {
	conn, err := net.Dial("tcp", tc.remoteAddr)
	if err != nil {
		log.Println("[TcpClient] dial failed: ", err.Error())
		return err
	}

	tc.conn = conn
	go tc.reader()
	return nil
}

func (tc *TcpClient) reader() {
	header := make([]byte, HeadLen)
	body := make([]byte, tc.codec.maxIncomingSize)

	var (
		err   error
		in    packet.IPacket
		bytes []byte
	)

	defer func() {
		if tc.conn != nil {
			_ = tc.conn.Close()
		}
		if tc.cancel != nil {
			tc.cancel()
		}

		if err != nil {
			if !(err == io.EOF || mmrpcs.IsClosedConnError(err)) {
				log.Println(fmt.Errorf("tcp %s disconnected: %s", tc.remoteAddr, err.Error()))
			}
		}
	}()

	tc.ack()

	for {
		in, err = tc.recv(header, body)
		if err != nil {
			return
		}

		tc.ack()
		for len(in.Remain()) > 0 {
			bytes, err = in.ReadBytes32()
			if err != nil {
				break
			}
			reader := packet.Reader(bytes)
			_ = tc.decodeRspAndDispatch(reader)
		}
	}
}

func (tc *TcpClient) recv(header []byte, body []byte) (packet.IPacket, error) {
	in, err := tc.codec.BlockDecodeBody(tc.conn, header, body)
	if err != nil {
		return nil, err
	}
	return in, err
}

func (tc *TcpClient) decodeRspAndDispatch(body packet.IPacket) error {
	mt, data, _ := tc.msgCodec.decode(body)
	switch mt {
	case TypeMessageHeartbeat, TypeMessageHeartbeatAck:
		return nil
	default:
		if data != nil {
			tc.r.put(data, nil)
		}
	}
	return nil
}

func (tc *TcpClient) keepalive() {
	ticker := time.NewTicker(time.Second)
	ping := tc.msgCodec.encode(nil, TypeMessageHeartbeat)
	defer ping.Return()

	prev := time.Now().Unix()
	timeout := time.Duration(0)
	outstandingPing := false

	for {
		select {
		case <-ticker.C:
			la := tc.lastAck.Load()
			if la > prev {
				prev = la
				ticker.Reset(time.Duration(la-time.Now().Unix()) + AckInterval)
				outstandingPing = false
				continue
			}

			if outstandingPing && timeout <= 0 {
				log.Println("[TcpClient] no heartbeat: ", tc.remoteAddr)
				_ = tc.conn.Close()
				return
			}

			if !outstandingPing {
				_ = tc.buf.Set(ping)
				outstandingPing = true
				timeout = PingTimeOut
			}
			sd := number_utils.Min[time.Duration](AckInterval, timeout)
			timeout -= sd
			ticker.Reset(sd)
		case <-tc.ctx.Done():
			return
		}
	}
}

func (tc *TcpClient) ack() {
	tc.lastAck.Store(time.Now().Unix())
}

func (tc *TcpClient) StateCompareAndSwap(old, new uint32) bool {
	return tc.state.CompareAndSwap(old, new)
}

func withRetry(handle func() error) error {
	retried := int32(0)
	for {
		err := handle()
		if err == nil {
			return nil
		}
		//exponential backoff
		time.Sleep(time.Millisecond * time.Duration(100<<retried))
		if retried >= MaxRetried {
			return mmrpcs.ErrMaxOfRetry
		}
		retried++
	}
}

func parseOptions(ops ...*DialOption) (dp *DialOption) {
	dp = new(DialOption)
	if len(ops) > 0 {
		op := ops[0]
		if op.MaxIncomingPacket > 0 {
			dp.MaxIncomingPacket = op.MaxIncomingPacket
		}
		dp.RemoteNodeId = op.RemoteNodeId
		dp.CurrentNodeId = op.CurrentNodeId
		dp.IsBlock = op.IsBlock
		dp.IsGzip = op.IsGzip
		dp.DisconnectHandler = op.DisconnectHandler
	}
	if dp.MaxIncomingPacket <= 0 {
		dp.MaxIncomingPacket = RpcMaxIncomingPacket
	}
	return
}
