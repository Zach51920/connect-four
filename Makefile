run: templ
	@go run main.go

templ:
	@templ generate

dev:
	@templ generate -watch -proxy=http://localhost:8080 &
	@air