FROM golang:1.17-alpine AS build
WORKDIR /go/src/cluster
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/cluster ./main.go
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.6 && \
	wget -qO /go/bin/grpc_health_probe https://github.com/grpc-ecosystem/\
grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/\
grpc_health_probe-linux-amd64 && \
chmod +x /go/bin/grpc_health_probe


FROM scratch
COPY --from=build /go/bin/cluster /bin/cluster
COPY --from=build /go/bin/grpc_health_probe /bin/grpc_health_probe
