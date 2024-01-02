update_gogo_proto:
	go install github.com/gogo/protobuf/proto
	go install github.com/gogo/protobuf/protoc-gen-gogofast
	go install github.com/gogo/protobuf/protoc-gen-gogo
	go install github.com/gogo/protobuf/gogoproto
	go install github.com/gogo/protobuf/protoc-gen-gofast
	go install github.com/gogo/protobuf/protoc-gen-gogofaster
	go install github.com/gogo/protobuf/protoc-gen-gogoslick

gen_proto:
	protoc --gogofaster_out=paths=source_relative:. ./core/actor/actor.proto
	protoc --gogofaster_out=paths=source_relative:. ./core/remote/remote.proto