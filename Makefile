pack:
	mkdir -p .output/awpark.alfredworkflow;
	cp -rf static/* .output/awpark.alfredworkflow/;
	go build -o .output/awpark.alfredworkflow/awpark main.go;