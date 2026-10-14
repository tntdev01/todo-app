portforwardingdb:
	kubectl port-forward -n todo-app-namespace service/todo-app-db-service 5432:5432

portforwardingapi:
	kubectl port-forward -n todo-app-namespace deployment/todo-api 8080:8080

migrateup:
	migrate -path migrations -database "postgresql://todo-app-db-user:todo-app-db-password@localhost:5432/todo-app-db?sslmode=disable" -verbose up 

migratedown:
	migrate -path migrations -database "postgresql://todo-app-db-user:todo-app-db-password@localhost:5432/todo-app-db?sslmode=disable" -verbose down 

test:
	go test -v -cover ./...

swagger:
	swag init -g main.go -d ./cmd/api,./internal/handler,./internal/model,./internal/dto --parseInternal
	
server:
	go run cmd/api/main.go
	
.PHONY: postgres createdb dropdb migrateup migratedown run test server mock