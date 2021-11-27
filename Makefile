TAG ?= 0.0.1
build-docker:
	docker build -t github.com/moratsam/cluster:$(TAG) .
	kind load docker-image github.com/moratsam/cluster:$(TAG)

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
