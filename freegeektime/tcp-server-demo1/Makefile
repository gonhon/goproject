all:server client

server:cmd/server/main.go
	   go build github.com/gonhon/tcp-server-demo1/cmd/server

client:cmd/client/main.go
	   go build github.com/gonhon/tcp-server-demo1/cmd/client

clean:
		rm -rf ./server
		rm -rf ./client