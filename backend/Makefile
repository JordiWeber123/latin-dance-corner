.PHONY: run-postgres test build

run-postgres:
	docker run --name postgres-test --rm -d \
	-p 127.0.0.1:5433:5432 -e POSTGRES_USER=testuser \
	-e POSTGRES_PASSWORD=testpass -e POSTGRES_DB=app_test \
	postgres:18



 