
PROD_IMGNAME=br
PROD_VERSION=$(shell sed 1q .version)

all:
	docker rmi $(PROD_IMGNAME):$(PROD_VERSION) || true
	docker build --label $(PROD_IMGNAME):$(PROD_VERSION) -t $(PROD_IMGNAME):$(PROD_VERSION) .
	docker image prune --force

