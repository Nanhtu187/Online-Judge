migrate-up:
	go run cmd/migrate/main.go up

migrate-down-1:
	go run cmd/migrate/main.go down 1

migrate:
	echo \# make migrate name="$(name)"
	go run cmd/migrate/main.go create $(name)

start:
	go run cmd/server/main.go start