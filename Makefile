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
