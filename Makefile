BINARY=build/go-vless-client

build:
	mkdir -p build
	go build -o $(BINARY) .

run: build
	./$(BINARY) $(ARGS)

clean:
	rm -rf build
