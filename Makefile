compose-up:
	docker-compose up -d

configure-barong: compose-up
	 docker-compose exec barong bundle exec rake db:create db:migrate db:seed

start: configure-barong

stop:
	docker-compose down

delete-payments-image:
	docker rmi local-testing-payments

delete-dinopay-gateway-image:
	docker rmi local-testing-dinopay-gateway

delete-payments-read-model-image:
	docker rmi local-testing-payments-read-model

build-test-runner:
	cd tests && go build -o bin/runner cmd/runner.go

run-tests: build-test-runner
	tests/bin/runner run constant dinopayOutboundSucceed --max-duration 10s rate 1/s

