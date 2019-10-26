all: bin/servethis

bin:
	mkdir -p bin

bin/servethis: bin $(shell find . -name '*.go') go.mod go.sum
	cd cmd/servethis && go build -o ../../$@

clean:
	rm -rf bin

install: bin/servethis
	install $< $(HOME)/bin/

test:
	go test ./... -v

.PHONY: clean all install test
