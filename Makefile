# Makefile

init:
	npm install
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/cosmtrek/air@latest
	go mod tidy

assets:
	tailwindcss -i ./internal/dist/main.css -o ./internal/dist/tailwind.css

generate:
	templ generate

run: assets generate
	go run ./cmd/server/...

test:
	go test ./...
