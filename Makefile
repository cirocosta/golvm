PLUGIN_NAME := lvmvol
ROOTFS_IMAGE := cirocosta/$(PLUGIN_NAME)-rootfs
ROOTFS_CONTAINER := rootfs
PLUGIN_FULL_NAME := cirocosta/$(PLUGIN_NAME)

all: install

fmt:
	go fmt
	cd ./lib && go fmt
	cd ./driver && go fmt

install:
	go install -v

test:
	cd ./lib && go test -v

rootfs-image:
	docker build -t $(ROOTFS_IMAGE) .


.PHONY: fmt install test
