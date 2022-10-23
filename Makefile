all: getdeps test build

getdeps:
	@echo "Installing gocyclo" && go get  -u github.com/fzipp/gocyclo/cmd/gocyclo
	@echo "Installing deadcode" && go get -u github.com/remyoudompheng/go-misc/deadcode
	@echo "Installing misspell" && go get -u github.com/client9/misspell/cmd/misspell
	@echo "Installing golangci-lint" && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.50.1
	@wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.50.1


verifiers: fmt lint cyclo spelling staticcheck

fmt:
	@echo "Running $@"
	@gofmt -d .

lint:
	@echo "Running $@"
	@$(go env GOPATH)/bin/golangci-lint run

cyclo:
	@echo "Running $@"
	@$(go env GOPATH)/bin/gocyclo -over 100 .

deadcode:
	@echo "Running $@"
	@$(go env GOPATH)/bin/deadcode -test $(shell go list ./...) || true

spelling:
	@$(go env GOPATH)/bin/misspell -i monitord -error `find .`

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck -- ./...


test: verifiers build
	go test -v ./...

testrace: verifiers build
	go test -v -race ./...

build:
	go build -o bin/urls-checker-cli main.go 

clean:
	rm ./bin/data -rf