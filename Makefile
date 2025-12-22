compile: app-local app-ssh

app-local:
	go build -o bin/app-local ./cmd/app-local/

app-ssh:
	go build -o bin/app-ssh ./cmd/app-ssh/

test:
	go test ./...

clear:
	rm bin/* | echo "dir empty"

all: clear test compile
