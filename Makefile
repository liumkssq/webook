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
	mockgen -source=internal/repository/user.go -destination=internal/repository/mocks/user.mock.go -package=repomocks
	mockgen -source=internal/repository/code.go -destination=internal/repository/mocks/code.mock.go -package=repomocks
	mockgen -source=internal/repository/dao/user.go -destination=internal/repository/dao/mocks/user.mock.go -package=daomocks
	mockgen -source=internal/repository/cache/user.go -destination=internal/repository/cache/mocks/user.mock.go -package=cachemocks
	mockgen -package=redismocks -destination=internal/repository/cache/redismocks/cmd.mock.go github.com/redis/go-redis/v9 Cmdable
