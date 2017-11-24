#!/bin/bash

set -o errexit

readonly IMAGES_DIR="/images"

main() {
	detach_device 0
	remove_image 0
}

detach_device() {
	local number=$1
	local device=/dev/loop$number

	echo "INFO:
  Detaching loopback device.

  DEVICE:  $device
  "

	losetup --detach $device
}

remove_image() {
	local number=$1
	local image=$IMAGES_DIR/lvm$number.img

	echo "INFO:
  Removing image.

  IMAGE:  $image
  "

	rm -f $image
}

main
