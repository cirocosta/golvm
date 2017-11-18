#!/bin/bash

set -o errexit

readonly IMAGES_DIR="/images"

main() {
	setup_images_dir
	create_vg "0"
	create_vg "1"
	create_vg "2"
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

create_vg() {
	local number=$1

	if [[ -z "$number" ]]; then
		echo "ERROR:
    create_vg expects an argument (number).

  Aborting.
  "
		exit 1
	fi

	local device=/dev/loop$number
	local image=$IMAGES_DIR/lvm$number.img
	local vg_name=volgroup$number

	echo "INFO:
  Starting to set up vg.

  DEVICE:   $device
  IMAGE:    $image
  "

	dd if=/dev/zero of=$image bs=1M count=50
	losetup $device $image
	echo ",,8e,," | sfdisk $device
	partx --update $device
	pvcreate $device
	vgcreate $vg_name $device
}

main "$@"
