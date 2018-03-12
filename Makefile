build:
		env GOOS=linux go build crawler/*.go
		env GOOS=linux go build -o bin/seq_down *.go
