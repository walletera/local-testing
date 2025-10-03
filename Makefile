.PHONY: help compose-up configure-barong start stop delete-payments-image delete-dinopay-gateway-image delete-payments-read-model-image build-test-runner run-tests clean delete-images

# Variables
DC = docker-compose
PAYMENTS_IMAGE = local-testing-payments
DINOPAY_GATEWAY_IMAGE = local-testing-dinopay-gateway
PAYMENTS_READ_MODEL_IMAGE = local-testing-payments-read-model
TESTS_DIR = tests
RUNNER_BIN = $(TESTS_DIR)/bin/runner

help:
	@echo "Available targets:"
	@echo "  compose-up                   Start docker-compose services in background"
	@echo "  configure-barong             Run db setup on barong service"
	@echo "  start                        Compose up & configure barong"
	@echo "  stop                         Stop all running containers"
	@echo "  delete-payments-image        Remove the payments docker image"
	@echo "  delete-dinopay-gateway-image Remove the dinopay-gateway docker image"
	@echo "  delete-payments-read-model-image Remove the payments read model docker image"
	@echo "  build-test-runner            Build the Go test runner"
	@echo "  run-tests                    Build and run tests"
	@echo "  clean                        Remove built binaries"
	@echo "  delete-images                Remove all custom docker images"

compose-up:
	$(DC) up -d

configure-barong: compose-up
	$(DC) exec barong bundle exec rake db:create db:migrate db:seed

start: configure-barong

stop:
	$(DC) down

delete-payments-image:
	docker rmi $(PAYMENTS_IMAGE)

delete-dinopay-gateway-image:
	docker rmi $(DINOPAY_GATEWAY_IMAGE)

delete-payments-read-model-image:
	docker rmi $(PAYMENTS_READ_MODEL_IMAGE)

delete-images: delete-payments-image delete-dinopay-gateway-image delete-payments-read-model-image

build-test-runner:
	cd $(TESTS_DIR) && go build -o bin/runner cmd/runner.go

run-tests: build-test-runner
	$(RUNNER_BIN) run constant dinopayOutboundSucceed --max-duration 10s --rate 1/s

clean:
	rm -rf $(TESTS_DIR)/bin/runner