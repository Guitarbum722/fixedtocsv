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

