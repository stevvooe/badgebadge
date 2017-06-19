IMAGE_NAME := badgebadge:latest

default: serve

build:
	docker build -t $(IMAGE_NAME) .

serve: clean build
	docker run -p 8080:8080 --rm $(IMAGE_NAME)

clean:
	-docker image rm $(IMAGE_NAME)

push:
	docker tag $(IMAGE_NAME) appaws/badgebadge:0.2
	docker push appaws/badgebadge:0.2
