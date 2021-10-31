hello:
	echo "Hello"

build:
	go build -o bin/snick main.go

run:
	go run main.go

a:
	@bin/snick test --files examples/hello/hello.yaml --rego examples/hello/hello.rego

d:
	bin/snick test --files examples/hello/hello.yaml --rego examples/hello/hello.rego -d

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go

all: hello build