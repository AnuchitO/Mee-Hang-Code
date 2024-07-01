.PHONY:
run:
	@echo "Running the program..."
	go mod tidy
	go run .

.PHONY:
start-db:
	@echo "Running the database..."
	docker-compose up -d

.PHONY:
e2e:
	@echo "Running the end-to-end tests..."
	@curl -s http://localhost:8080/skills/go | jq '.data | has("key")'

.PHONY:
get-skill:
	@echo "Getting a skill..."
	@curl -s http://localhost:8080/skills/go