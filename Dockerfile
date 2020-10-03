FROM golang:1.15 as build-go
WORKDIR /file-reader
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/file-reader-server .

FROM alpine:latest
RUN addgroup -S file-reader && file-reader -S file-reader -G file-reader
USER file-reader
WORKDIR /home/file-reader
COPY --from=build-go /bin/file-reader-server ./
EXPOSE 3000
ENTRYPOINT ["./file-reader-server"]
