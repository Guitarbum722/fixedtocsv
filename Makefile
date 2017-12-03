BINARY = fixedtocsv

build: clean
	go build -o bin/$(BINARY)

test:
	go test -v .

clean:
	go clean
	rm -f $(BINARY)
	rm -f bin/*

install: clean build
	cp -f bin/$(BINARY) /usr/local/bin/$(BINARY)

uninstall: 
	rm -f /usr/local/bin/$(BINARY)*

darwin:
	GOOS=darwin GOARCH=amd64 go build -v -o bin/$(BINARY)-darwin
	tar -czvf bin/$(BINARY)-darwin.tar.gz bin/$(BINARY)-darwin
freebsd:
	GOOS=freebsd GOARCH=amd64 go build -v -o bin/$(BINARY)-freebsd
	tar -czvf bin/$(BINARY)-freebsd.tar.gz bin/$(BINARY)-freebsd

windows:
	GOOS=windows GOARCH=amd64 go build -v -o bin/$(BINARY)-windows
	tar -czvf bin/$(BINARY)-windows.tar.gz bin/$(BINARY)-windows

linux:
	GOOS=linux GOARCH=amd64 go build -v -o bin/$(BINARY)-linux
	tar -czvf bin/$(BINARY)-linux.tar.gz bin/$(BINARY)-linux

release: clean darwin freebsd windows linux