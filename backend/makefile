APP_NAME=api-service

run-local-infra:
	docker-compose -f docker-compose.infra.yml up

stop-local-infra:
	docker-compose -f docker-compose.infra.yml down

run-app:
	$(MAKE) -C app/$(APP_NAME) run

test-app:
	$(MAKE) -C app/$(APP_NAME) test

test-app-report:
	$(MAKE) -C app/$(APP_NAME) test-report

test-core:
	$(MAKE) -C core test

test-core-report:
	$(MAKE) -C core test-report

test:
	$(MAKE) test-core
	$(MAKE) test-app
	$(MAKE) test-app APP_NAME=token-generator

generate-token:
	$(MAKE) -C app/token-generator run
