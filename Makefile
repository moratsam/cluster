compile:
	for f in api/v1/*/*.proto; do					\
		protoc $$f										\
			--go_out=. 									\
			--go-grpc_out=. 							\
			--go_opt=paths=source_relative 		\
			--go-grpc_opt=paths=source_relative	\
			--proto_path=. ;							\
	done

test:
	go test -race ./...
