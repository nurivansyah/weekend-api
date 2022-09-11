PROJECT_ID=weekend-api
IMAGE_TAG=weekend-api
GITHUB_SHA?=latest
LOCAL_TAG=$(IMAGE_TAG):$(GITHUB_SHA)
REMOTE_TAG=gcr.io/$(PROJECT_ID)/$(LOCAL_TAG)
VERSION?=latest
VERSION_LOCAL_TAG=$(IMAGE_TAG):$(VERSION)
VERSION_REMOTE_TAG=gcr.io/$(PROJECT_ID)/$(VERSION_LOCAL_TAG)

docker-images:
	docker images | grep $(IMAGE_TAG)

docker-build:
	docker build . --tag  $(LOCAL_TAG)

docker-tag:
	docker tag $(LOCAL_TAG) $(REMOTE_TAG)

docker-pull:
	docker pull $(REMOTE_TAG)

docker-push:
	docker push $(REMOTE_TAG)

docker-tag-version:
	docker tag $(REMOTE_TAG) $(VERSION_REMOTE_TAG)

docker-push-version:
	docker push $(VERSION_REMOTE_TAG)
