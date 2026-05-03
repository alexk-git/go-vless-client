BINARY=build

build:
    go build -o $(BINARY) .

run: build
    ./$(BINARY) $(ARGS)

clean:
    rm -f $(BINARY)
