FROM golang:1.13.1 as build-stage

WORKDIR /go/src/github.com/cyruzin/hexagony

COPY . .

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/hexagony

FROM alpine:latest  

RUN apk add ca-certificates

COPY --from=build-stage /go/bin/hexagony .

EXPOSE 8000

CMD ["./hexagony"]