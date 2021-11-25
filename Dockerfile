FROM golang:1.17-alpine AS build
WORKDIR /go/src/cluster
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/cluster ./main.go


FROM scratch
COPY --from=build /go/bin/cluster /bin/cluster
