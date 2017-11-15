#!/bin/bash

set -o errexit

main() {
	setup_dependencies

	echo "INFO:
  Done! Finished setting up travis machine.
  "
}

setup_dependencies() {
	echo "INFO:
  Setting up dependencies.
  "

	sudo apt update -y
	sudo apt install realpath -y
	sudo apt install --only-upgrade docker-ce -y

	git --version
	git config --global user.name "WeDeploy CI"
	git config --global user.email "ci@wedeploy.com"

	docker info
}

main
