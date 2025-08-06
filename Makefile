.PHONY: build

build:
	docker-compose build $(svc)

.PHONY: status logs start stop clean

ps:
	docker-compose ps $(svc)

logs:
	docker-compose logs -f $(svc)

up:
	docker-compose up -d $(svc)

start:
	docker-compose start $(svc)

stop:
	docker-compose stop $(svc)

down:stop
	docker-compose down -v --remove-orphans

attach:
	docker-compose exec $(svc) bash

prune:
	docker system prune --all --volumes

.PHONY: test

test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down --volumes

.PHONY: gen gtest swag

gen:
	protoc \
	--go_out=proto/pb \
	--go_opt=paths=source_relative \
	--go-grpc_out=proto/pb \
	--go-grpc_opt=paths=source_relative \
	--proto_path=proto/protofiles \
	proto/protofiles/*.proto

gtest:
	go test -v -cover -coverprofile coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

swag:
	swag init -o adapter/inbound/http/docs -g adapter/inbound/http/http.go
