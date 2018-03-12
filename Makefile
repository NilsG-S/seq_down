build:
	dep ensure
	env GOOS=linux go build -i crawler/*.go
	env GOOS=linux go build -o bin/seq_down *.go
