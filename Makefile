# Makefile

# Variables
DOCKER_COMPOSE_PATH=./resources/local-setup/docker-compose.yml
KAFKA_TOPICS=heart-beat-events.json mapreduce-topic
CQL_SCRIPT_PATH=./create_cql_tables.cql

GO_MOD=go.mod

export APP_ENV=dev

all: docker-up create-kafka-topics create-cql-tables go-deps run-app

docker-up:
	@echo "Starting Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_PATH) up -d

create-kafka-topics:
	@echo "Creating Kafka topics..."
	for topic in $(KAFKA_TOPICS); do \
		docker exec kafka-container kafka-topics --create --topic $$topic --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1; \
	done

create-cql-tables:
	@echo "Initialising cassandra..."
	docker exec cassandra-container cqlsh -f $(CQL_SCRIPT_PATH)
	@echo "successfully initialised cassandra"

go-deps:
	@echo "Downloading Go dependencies..."
	go mod download

run-app:
	@echo "Running Go application..."
	go run cmd/main.go
