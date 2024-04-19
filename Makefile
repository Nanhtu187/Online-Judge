gen-proto:
	for proto_file in $$(find . -name "*.proto" -not -path "./proto_definition/.third_party/*"); do \
		protoc -I proto_definition/proto -I ./proto_definition/.third_party/googleapis -I ./proto_definition/.third_party/envoyproxy \
		--go_out proto/rpc --go_opt paths=source_relative \
		--go-grpc_out proto/rpc --go-grpc_opt paths=source_relative \
		--grpc-gateway_out proto/rpc --grpc-gateway_opt paths=source_relative \
		$$proto_file;\
	done

