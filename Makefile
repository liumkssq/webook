.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f flycash/webook:v0.0.1
	@docker build -t flycash/webook:v0.0.1 .


.PHONY: etcdctl
etcdctl:
	powershell -File ps/etcdctl.ps1



.PHONY: docker-compose
docker-compose:
	docker-compose down
	docker-compose up -d

.PHONY: mock
mock:
