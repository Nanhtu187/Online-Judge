gen-proto:
	for proto_file in $$(find . -name "*.proto" -not -path "./proto_definition/.third_party/*"); do \
	echo "Generate sources for $${proto_file}"; \
	protoc -I /usr/local/include \
			-I ./proto_definition \
			-I proto_definition/proto \
			-I ./proto_definition/.third_party/googleapis \
			-I ./proto_definition/.third_party/envoyproxy \
			-I ./proto_definition/.third_party/protoc-gen-swagger \
			--go_out proto/rpc --go_opt=module=proto/rpc \
			--go-grpc_out proto/rpc --go-grpc_opt=module=proto/rpc \
			--grpc-gateway_out proto/rpc --grpc-gateway_opt=module=proto/rpc \
			--validate_out="lang=go,module=proto/rpc:proto/rpc"  \
			$$proto_file; \
    done

install-tools:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install github.com/envoyproxy/protoc-gen-validate