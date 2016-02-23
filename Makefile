run:
	go-bindata assets/...
	go build ./...
	./dev-console

compile-all: compile-darwin-amd64 compile-linux-amd64 compile-windows-amd64

compile-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin-amd64/dev-console -v .

compile-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/dev-console -v .

compile-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/dev-console.exe -v .
