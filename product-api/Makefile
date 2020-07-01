check_install: 
	which swagger || G0111MODULE=on go get -u github.com/go-swagger/go-swagger/cmd/swagger
swagger: check_install
	G0111MODULE=on swagger generate spec -o ./swagger.yaml --scan-models