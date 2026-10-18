.PHONY: test doc wire run gen build tidy proto push

MODULE   		:= $(shell go list -m)
APP_NAME 		:= $(notdir $(MODULE))
ORG_NAME 		:= $(notdir $(patsubst %/$(APP_NAME),%,$(MODULE)))

IMAGE_BASE 		:= ${ORG_NAME}/${APP_NAME}
REGISTRY   		:= docker.io
GIT_SHORT_SHA 	:= $(shell git describe --always --dirty 2>/dev/null || echo "dirty")

wire:
	@cd ./cmd/api && go tool wire

doc:
	@go run ./cmd/openapi

gen:
	go generate ./...

PROTOC_GEN_GO := $(shell go tool -n protoc-gen-go)

proto:
	protoc --plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
	    --proto_path=. \
		--go_out=. \
		--go_opt=paths=source_relative \
		pkg/schema/*.proto

MIGRATE_CLI		:= go run -modfile=tools.mod -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate
DATABASE		:=
action			?= 

migrate:
	${MIGRATE_CLI} create -ext sql -dir ./db/migrations -format 20060102_150405 -tz "Asia/Shanghai" ${action}
migrate-up:
	${MIGRATE_CLI} -path ./db/migrations -database ${DATABASE} up
migrate-down:	
	${MIGRATE_CLI} -path ./db/migrations -database ${DATABASE} down

build := 1

compose-up: assert-name
	ORG_NAME=$(ORG_NAME) \
	APP_NAME=$(APP_NAME) \
	COMMIT_SHA=$(GIT_SHORT_SHA) \
	docker compose up $(if $(build),--build,--no-build)
compose-clean: assert-name
	ORG_NAME=$(ORG_NAME) \
	APP_NAME=$(APP_NAME) \
	docker compose down --volumes --remove-orphans
	docker image prune -f

assert-name:
	@if [ -z "$(ORG_NAME)" ] || [ -z "$(APP_NAME)" ]; then \
		echo "ORG_NAME or APP_NAME is empty"; \
		exit 1; \
	fi

env ?= local

run:
	go mod tidy
	cd ./cmd/api && APP_ENV=${env} go run -ldflags="-X main.commitSHA=${GIT_SHORT_SHA}" .

tag 		?= ${GIT_SHORT_SHA}
config		?= debug
IMAGE_FULL 	:= $(REGISTRY)/$(IMAGE_BASE):$(tag)
IMAGE_URL 	:= docker://$(IMAGE_FULL)

build:
	docker build \
		--build-arg COMMIT_SHA=${GIT_SHORT_SHA} \
		--target ${config} \
		--platform linux/amd64 \
		--tag ${IMAGE_FULL} \
		.

push:
	docker push $(IMAGE_FULL) $(IMAGE_URL)

tidy:
	go mod tidy
	go mod edit -modfile=tools.mod -module=$$(go list -m)
	go list -m -f '{{if and (not .Indirect) (not .Main) .Version}}{{.Path}}@{{.Version}}{{end}}' all | xargs -r go get -modfile=tools.mod
	go mod tidy -modfile=tools.mod

