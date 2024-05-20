IMAGE_NAME=yevhenhrytsai/dice:v1.0.6


compile:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o dice

build:
	docker build . -t ${IMAGE_NAME}

push:
	docker push ${IMAGE_NAME}
