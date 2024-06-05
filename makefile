GITHUB_TOKEN =  $(shell cat ghp_sercet)

window:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o prst.exe ./
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o prst ./

.PHONY: window linux


dev:
	go  run ./ -c config.dev.yml
release: 
	go run  ./ -c config.yml

.PHONY: dev release

swag:
	swag init --parseDependency --parseDepth=6 --instanceName prst -o ./docs/   
.PHONY:


build:
	@docker login ghcr.io -u seedlings-calm -p $(GITHUB_TOKEN)
	@docker build -t ghcr.io/seedlings-calm/prst:latest . 
	docker push ghcr.io/seedlings-calm/prst:latest

.PHONY: build

