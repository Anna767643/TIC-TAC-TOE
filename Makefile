APP_DIR = cmd/app/
PACKAGES := ./internal/...
run: 
	go run $(APP_DIR)main.go 

test: 
	go test $(PACKAGES) -v
	go test -race  -cover $(PACKAGES) 

gcov_report:
	go test $(PACKAGES) -coverprofile=coverage.txt
	go tool cover -html coverage.txt -o index.html
	open index.html

fmt:
#	Формотирование
	gofmt -l . 
# 	Потенциальные ошибки
	go vet $(PACKAGES) 
# 	Ревью кода
	golangci-lint run $(PACKAGES) 
# 	Актуальность mod
	go mod tidy -diff

clean:
	rm -f index.html coverage.txt

