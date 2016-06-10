How to run:
`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o testgateway .`
`docker build -f Dockerfile -t testgateway .`
`docker-compose up`
