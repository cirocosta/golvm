#!/bin/bash

set -o errexit

readonly IMAGES_DIR="/images"

main() {
	setup_images_dir

	setup_device 0
	create_pv 0
	create_vg 0
}

setup_images_dir() {
	echo "INFO:
  Setting images directory.

  IMAGES_DIR: $IMAGES_DIR
  "

	if [[ ! -d "$IMAGES_DIR" ]]; then
		mkdir -p $IMAGES_DIR
	fi
}

setup_device() {
	local number=$1
	local device=/dev/loop$number
	local image=$IMAGES_DIR/lvm$number.img

	echo "INFO: Setting up device
  NUMBER: $number
  DEVICE: $device
  IMAGE:  $image
  "

	test -z "$number" &&
		{
			echo "a device number must be specified"
			exit 1
		}

	dd if=/dev/zero of=$image bs=1M count=50
	losetup $device $image
}

create_pv() {
	local number=$1
	local device=/dev/loop$number

	echo "INFO: Preparing physical volume
  NUMBER: $number
  DEVICE: $device
  "

	test -z "$number" &&
		{
			echo "a device number must be specified"
			exit 1
		}

	echo ",,8e,," | sfdisk $device
	partx --update $device
	pvcreate $device
}

create_vg() {
	local number=$1
	local device=/dev/loop$number
	local vg_name=volgroup$number

	echo "INFO: Preparing volume group
  NUMBER: $number
  DEVICE: $device
  VG_NAME:$vg_name
  "

	test -z "$number" &&
		{
			echo "a device number must be specified"
			exit 1
		}

	vgcreate $vg_name $device
}

main "$@"
