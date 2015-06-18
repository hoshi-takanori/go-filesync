all: server client

server: main.go file.go sync.go
	go build -o server

client: client.go file.go sync.go
	go build -o client -tags client

clean:
	rm -f server client
