update_gogo_proto:
	go install github.com/gogo/protobuf/proto
	go install github.com/gogo/protobuf/protoc-gen-gogofast
	go install github.com/gogo/protobuf/protoc-gen-gogo
	go install github.com/gogo/protobuf/gogoproto
	go install github.com/gogo/protobuf/protoc-gen-gofast
	go install github.com/gogo/protobuf/protoc-gen-gogofaster
	go install github.com/gogo/protobuf/protoc-gen-gogoslick

gen_proto:
	#protoc --gogofaster_out=paths=source_relative:. ./core/actor/actor.proto
	#protoc --gogofaster_out=paths=source_relative:. ./core/remote/remote.proto
	protoc --go_out=. --go_opt=paths=source_relative \
           ./core/actor/actor.proto ./core/remote/remote.proto

update_protobuf:
	go install google.golang.org/protobuf/cmd/protoc-gen-go