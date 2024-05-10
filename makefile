.PHONY: run build docker-build docker-run db-run db-remove docker-remove test benchmark check-pprof
cur := $(shell pwd)

db-container-id := $(shell docker ps -a| grep mysql | awk '{print $$1}')
trading-engine-container-id := $(shell docker ps -a | grep trading-engine | awk '{print $$1}')

args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

test: 
	@go clean -testcache & go test -v ./test/order_book_test.go

benchmark:
	@go test -v -bench=. -run=none ./test/...

run:
	@go run main.go server -l true -c configs/dev.yaml

build:
	@go build -v -o bin/trading-engine ./main.go

check-pprof:
	@go tool pprof -http :8080 $(call args,defaultstring)

docker-build:
	@docker build --tag n960321/trading-engine:latest --file build/dockerfile .

docker-remove:
	@docker rm -f $(trading-engine-container-id)

docker-run:
	@docker run --name trading-engine \
	-p 8080:8080 \
	--link mysql:mysql \
	--volume $(cur)/configs:/app/configs \
	n960321/trading-engine:latest

db-remove:
	docker rm -f $(db-container-id)

db-run:
	docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=123456 mysql:8.3.0