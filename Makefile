build:
	GOOS=darwin GOARCH=amd64 go build -o bin/awpark.amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/awpark.arm64 main.go
	mkdir -p bin
	lipo -create -output ./bin/awpark bin/awpark.amd64 bin/awpark.arm64