FROM golang:1.11
WORKDIR /go/src/github.com/Analyse4/digimon
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go install -v ./...

FROM alpine:latest
WORKDIR /root
COPY --from=0 /go/bin/digimon .
EXPOSE 2244
ENTRYPOINT ["./digimon"]