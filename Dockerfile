FROM golang:1.21.5-alpine3.19 as build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . ./

WORKDIR ./cmd/server

RUN go build -v -o server

FROM alpine:latest  

RUN apk add ca-certificates

COPY --from=build /app/cmd/server /app/server

EXPOSE 8000

CMD ["/app/server/server"]