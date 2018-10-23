# Install Golang and build
FROM golang:stretch AS build
WORKDIR /go/src/github.com/Kunal-Shah-Bose/playlist-scraper
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Install Heroku using Alpine base image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=build /go/src/github.com/Kunal-Shah-Bose/playlist-scraper/app .
CMD ["./app"]
