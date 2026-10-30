APP_DIR = cmd/app/
PACKAGES := ./internal/...
run: 
	go run $(APP_DIR)main.go 

test: 
	 go test $(PACKAGES) -v

test_coverage:
	go test $(PACKAGES) -cover

test_html:
	go test $(PACKAGES) -coverprofile=coverage.txt
	go tool cover -html coverage.txt -o index.html
	open index.html

test_analiz:
# 	Потенциальные ошибки
	go vet ./...
# 	Ревью кода
	golangci-lint run ./...
# 	Актуальность mod
	go mod tidy -diff
# 	Проверка на гонки, долго обрабатывается
	go test -race -cover ./...

clean:
	rm -f index.html coverage.txt

