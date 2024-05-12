window:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o prst.exe ./
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o prst ./

.PHONY: window linux


swag:
	swag init --parseDependency --parseDepth=6 --instanceName prst -o ./docs/   
.PHONY: