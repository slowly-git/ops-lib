.PHONY:	build push

PREFIX = gcr.io/google-containers
IMAGE = fluentd-elasticsearch
TAG = v2.0.3

build:
	docker build --pull -t $(PREFIX)/$(IMAGE):$(TAG) .

push:
	gcloud docker -- push $(PREFIX)/$(IMAGE):$(TAG)
