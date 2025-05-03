.PHONY: rebuild
rebuild:
	docker-compose down --rmi local --volumes --remove-orphans
	docker-compose up --build

generate-docs:
	swag init --parseInternal -g ./cmd/main.go