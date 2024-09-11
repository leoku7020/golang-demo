CMD_DIR = cmd
CMDS = $(patsubst $(CMD_DIR)/%,%,$(wildcard $(CMD_DIR)/*))
# Overrides command line arguments
PROJ :=
DEV := dev
CUR_COMMIT :=
COMMIT := master
DRY_RUN :=

# make
default: clean install gen test

init:
	cp .env.example .env

gen: generate fmt adjust

install:
	# for enum
	go install github.com/abice/go-enum@v0.5.5
	# for mocks
	go install github.com/vektra/mockery/v2@v2.20.0

	# for protobuf
	# - ref: https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/introduction/#prerequisites
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	go install github.com/favadi/protoc-go-inject-tag@v1.4.0
	go install github.com/envoyproxy/protoc-gen-validate@v0.9.1

	# for formatting
	go install golang.org/x/tools/cmd/goimports@latest

	# for SQL migration
	go install github.com/pressly/goose/v3/cmd/goose@latest

generate: generate-without-mockery generate-mockery

generate-without-mockery:
	# gen proto
	protoc -I ./proto \
		--go_out ./proto \
			--go_opt paths=source_relative \
		--go-grpc_out ./proto \
			--go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./proto \
			--grpc-gateway_opt paths=source_relative \
			--grpc-gateway_opt logtostderr=true \
			--grpc-gateway_opt generate_unbound_methods=true \
		--openapiv2_out ./proto \
			--openapiv2_opt logtostderr=true \
			--openapiv2_opt use_go_templates=true \
		--validate_out="lang=go,paths=source_relative:./proto" \
		./proto/*/*.proto

	# inject proto
	protoc-go-inject-tag -input=./proto/*/*.pb.go

	# gen enum
	go generate ./...

generate-mockery:
	# gen mocks
	mockery --all --inpackage --dir=./internal

fmt:
	gofmt -s -w -l .
	goimports --local gs-mono -w $(shell find . -type f -name '*.go' -not -path "./proto/*")

adjust:
	# adjust go.mod
	go mod tidy

clean:
	# remove protogen
	find . -type f -name '*.pb.go' -delete
	find . -type f -name '*.pb.*.go' -delete
	find . -type f -name '*.swagger.json' -delete
	find . -type f -name '*.pb.gw.go' -delete
	find . -type f -name '*.pb.validate.go' -delete

	# remove mocks
	find . -type f -name 'mock_*.go' -delete

	# remove enum
	find . -type f -name '*_enum.go' -delete

	# remove out
	find . -type f -name '*.out' -delete

	# remove dependency directories
	rm -rf bin/
	rm -rf sql/

test:
	$(eval @_TMP := $(shell git diff $(COMMIT) --name-only | grep ".go" | xargs -I {} dirname {} | uniq | sed 's/^/.\//' ))
	@[ "$(@_TMP)" ] && go test -v -race $(@_TMP) || ( echo "nothing to test!"; exit 0 )

ci-test-purge:
	# sql migration
	PURGE_DOCKER_AT_BEGIN=1 go test -v ./internal/adapter/repository -count=1

ci-test: ci-test-purge
	# TODO: go test ./pkg/...

	# run all unit tests
	go test -race -coverprofile=profile.out ./internal/...

	# calculate coverage rate
	cat profile.out | grep -v "mock_" | grep -v "_enum.go" > coverage.out
	@echo ~~~
	@go tool cover -func coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}' | xargs -I {} echo "total coverage: {}%"

ci-build:
	CUR_COMMIT="$(CUR_COMMIT)" COMMIT="$(COMMIT)" DRY_RUN=$(DRY_RUN) sh build.sh

check-test-coverage:
	bash ./scripts/testcoverage.sh

stop:
	docker-compose stop

check:
	@[ "$(PROJ)" ] && echo PROJ checked \
		|| ( echo "PROJ is not set. ex: PROJ=example"; exit 1 )

$(CMDS): check down
	# command=$@, project=$(PROJ)

	# build docker for $(PROJ)-$@

	DEV=$(DEV) PROJ=$(PROJ) docker-compose build $(PROJ)-$@
	DEV=$(DEV) PROJ=$(PROJ) docker-compose stop  $(PROJ)-$@
	DEV=$(DEV) PROJ=$(PROJ) docker-compose rm -f $(PROJ)-$@
	DEV=$(DEV) PROJ=$(PROJ) docker-compose up -d $(PROJ)-$@

.PHONY: install generate fmt adjust clean gen check default init
.PHONY: down stop $(CMDS) # docker related
.PHONY: test ci-test ci-build ci-test-purge # test, ci / cd
