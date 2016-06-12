default: testgateway
	docker build -f Dockerfile -t testgateway .

testgateway:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o testgateway

clean:
	rm testgateway
