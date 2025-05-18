FROM golang:1.24.3-alpine AS build
RUN apk add upx
WORKDIR $GOPATH/src/wol_relay
COPY . .
RUN GOOS=linux GOARCH=$TARGETARCH go build -v -o /go/bin/wol_relay
RUN upx /go/bin/wol_relay

# final stage
FROM alpine

COPY --from=build /go/bin/wol_relay .
ENTRYPOINT ["./wol_relay" ]
