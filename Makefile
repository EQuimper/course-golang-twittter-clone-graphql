mock:
	mockery --all --keeptree

migrate:
	migrate -source file://postgres/migrations \
			-database postgres://postgres:postgres@127.0.0.1:5432/twitter_clone_development?sslmode=disable up

rollback:
	migrate -source file://postgres/migrations \
			-database postgres://postgres:postgres@127.0.0.1:5432/twitter_clone_development?sslmode=disable down 1

drop:
	migrate -source file://postgres/migrations \
			-database postgres://postgres:postgres@127.0.0.1:5432/twitter_clone_development?sslmode=disable drop

migration:
	@read -p "Enter migration name: " name; \
		migrate create -ext sql -dir postgres/migrations $$name

run:
	go run cmd/graphqlserver/*.go

generate: 
	go generate ./..