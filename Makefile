IMAGE_NAME=yevhenhrytsai/dice:v1.0.1


build:
	docker build . -t ${IMAGE_NAME}

push:
	docker push ${IMAGE_NAME}
