ROOT_DIR := $(shell pwd)

go:
	go run ./cmd/main.go

down:
	docker compose down --remove-orphans

up: down build
	docker compose up

build:
	docker build -t rinha2023_app .

setup-gen:
	cd $(ROOT_DIR)/load-test/geradores/faker/ && npm i

gen: setup-gen
	node ./load-test/geradores/faker/gerar-pessoas > $(ROOT_DIR)/load-test/resources/pessoas-payloads.tsv
	node ./load-test/geradores/faker/gerar-termos-busca > $(ROOT_DIR)/load-test/resources/termos-busca.tsv

gatit:
	cd $(ROOT_DIR)/load-test && npm run computerdatabase
	echo "test over, sleeping 6 and then will count people..."
	sleep 6
	curl http://localhost:9999/contagem-pessoas

observe:
	chmod +x $(ROOT_DIR)/scripts/observe.sh
	$(ROOT_DIR)/scripts/observe.sh

pprof-cpu:
	mkdir -p $(ROOT_DIR)/logs
	curl -s -o "$(ROOT_DIR)/logs/cpu-$$(date +%Y%m%dT%H%M%S).pb.gz" \
		"http://localhost:6060/debug/pprof/profile?seconds=30"
	@echo "saved; analyze with: go tool pprof logs/cpu-*.pb.gz"
