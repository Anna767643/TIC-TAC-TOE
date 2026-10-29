APP_DIR = cmd/app/
PACKAGES := ./internal/datasource ./internal/domain ./internal/web

make: 
	go run $(APP_DIR)main.go 

test: 
	 go test $(PACKAGES) -v

test_coverage:
	go test $(PACKAGES) -cover

test_html:
	go test $(PACKAGES) -coverprofile=coverage.txt
	go tool cover -html coverage.txt -o index.html
	open index.html

clean:
	rm index.html
	rm coverage.txt

