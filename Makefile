compile:
	go mod download
	go build -tags static -v
