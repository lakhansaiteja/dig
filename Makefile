BINARY_NAME=dig

build:
	 GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin-amd64
	 GOARCH=arm64 GOOS=darwin go build -o ${BINARY_NAME}-darwin-arm64
	 GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux
	 GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows

all: build
	 @echo "done"

clean:
	 @go clean
	 @rm ${BINARY_NAME}-darwin-amd64
	 @rm ${BINARY_NAME}-darwin-arm64
	 @rm ${BINARY_NAME}-linux
	 @rm ${BINARY_NAME}-windows
	 @echo "cleaned"

docker-build:
	 docker build --tag lakhan-dig:v1 .

docker-run:
	docker run -it --rm -d -p 8080:8080 --name lakhan-dig lakhan-dig:v1

docker: docker-build docker-run
	@echo "dig app running in docker. Try the APIs!"

k8s: docker-build
	kubectl apply -f deployment.yaml