FROM golang:1.19-alpine as builder

ENV GO111MODULE=on

WORKDIR /app
COPY . .

RUN apk --no-cache add git alpine-sdk build-base gcc

RUN go get

RUN go build -o gree-havc-mqtt-bridge-go && chmod +x ./gree-havc-mqtt-bridge-go


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gree-havc-mqtt-bridge-go .
COPY ./config.yaml .
CMD ["./gree-havc-mqtt-bridge-go", "-c", "./config.yaml"]