.PHONY: app
app:
	docker-compose -f ./docker-compose.app.yml up --build

.PHONY: resources
resources:
	docker-compose -f ./docker-compose.resources.yml up

.PHONY: test
test:
	go test -count=1 ./...

.PHONY: mocks
mocks:
	mockgen -destination=./internal/platform/mongodb/mock/client.go github.com/mazxaxz/donut-batcher/internal/platform/mongodb Clienter,SingleResulter
	mockgen -destination=./internal/platform/rabbitmq/mock/client.go github.com/mazxaxz/donut-batcher/internal/platform/rabbitmq Publisher
	mockgen -destination=./internal/batch/mock/service.go github.com/mazxaxz/donut-batcher/internal/batch Service
