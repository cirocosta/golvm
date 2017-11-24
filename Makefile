PLUGIN_NAME := lvmvol
ROOTFS_IMAGE := cirocosta/$(PLUGIN_NAME)-rootfs
ROOTFS_CONTAINER := rootfs
PLUGIN_FULL_NAME := cirocosta/$(PLUGIN_NAME)

all: install

fmt:
	go fmt
	cd ./lib && go fmt
	cd ./driver && go fmt
	cd ./lvmctl && go fmt
	cd ./lvmctl/commands && go fmt

install:
	go install -v
	cd ./lvmctl && go install -v

test:
	cd ./lib && go test
	cd ./driver && go test

rootfs-image:
	docker build -t $(ROOTFS_IMAGE) .


rootfs: rootfs-image
	docker rm -vf $(ROOTFS_CONTAINER) || true
	docker create --name $(ROOTFS_CONTAINER) $(ROOTFS_IMAGE) || true
	mkdir -p plugin/rootfs
	rm -rf plugin/rootfs/*
	docker export $(ROOTFS_CONTAINER) | tar -x -C plugin/rootfs
	docker rm -vf $(ROOTFS_CONTAINER)


plugin: rootfs
	docker plugin disable $(PLUGIN_NAME) || true
	docker plugin rm --force $(PLUGIN_NAME) || true
	docker plugin create $(PLUGIN_NAME) ./plugin
	docker plugin enable $(PLUGIN_NAME)


plugin-push: rootfs
	docker plugin rm --force $(PLUGIN_FULL_NAME) || true
	docker plugin create $(PLUGIN_FULL_NAME) ./plugin
	docker plugin push $(PLUGIN_FULL_NAME)


.PHONY: fmt install test
