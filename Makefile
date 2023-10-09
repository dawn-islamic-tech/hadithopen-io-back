SEARCH_PATH := "&search_path=public"

go-gen:
	go generate ./...

gen: go-gen

lint:
	golangci-lint run ./...

migrate-create:
	migrate create -ext sql -dir migrations/postgres -seq $(name)

migrate-up:
	 migrate -database $(conn)$(SEARCH_PATH) -path migrations/postgres up

migrate-down:
	 migrate -database $(conn)$(SEARCH_PATH) -path migrations/postgres down

migrate-up-docker:
	 migrate -database "postgres://hadith:hadith@localhost:8787/hadith?sslmode=disable"$(SEARCH_PATH) -path migrations/postgres up

migrate-down-docker:
	migrate -database "postgres://hadith:hadith@localhost:8787/hadith?sslmode=disable"$(SEARCH_PATH) -path migrations/postgres down
