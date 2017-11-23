# golvm

`golvm` provides both a library for dealing with LVM operations as well as a docker plugin (`lvmvol`).


<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Plugin](#plugin)
  - [Activation](#activation)
    - [Default configuration](#default-configuration)
    - [Custom configuration](#custom-configuration)
  - [Usage](#usage)
    - [Create regular volume](#create-regular-volume)
    - [Create thinly provisioned volume](#create-thinly-provisioned-volume)
    - [Create snapshot volume](#create-snapshot-volume)
    - [Create thin snapshot volume](#create-thin-snapshot-volume)
    - [List volumes](#list-volumes)
    - [Inspect volume](#inspect-volume)
- [lvmctl](#lvmctl)
  - [Create regular volume](#create-regular-volume-1)
  - [Dependencies](#dependencies)
  - [Usage](#usage-1)
    - [Create thin volume](#create-thin-volume)
    - [Create snapshot](#create-snapshot)
    - [Create thin snapshot](#create-thin-snapshot)
    - [Encrypted volume](#encrypted-volume)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## Plugin

### Activation

Every volume managed by the plugin is mounted under `/mnt/lvmvol/volumes` by default (can be configured). 

Using the default configuration means that it's required that the directory `/mnt` exists before the plugin is enabled.


#### Default configuration

Using the default configuration is the simples way to go:

0. Optionally create a whitelist of volumegroups
1. Install the plugin
2. Start using it


```sh
# Whitelist some volume groups that we want to make use of
echo "volumegroup0" | sudo tee "/mnt/lvmvol/whitelist.txt"

# Install the plugin but don't enable it
docker plugin install \
        --grant-all-permissions \
        cirocosta/golvm

# Create a volume
docker volume create \
        --driver lvmvol \
        --size 10M \
        myvol
```

By default the following parameters are used:

```
VOLUME_MOUNT_ROOT:      /mnt/lvmvol/volumes
WHITELIST_FILE:         /mnt/lvmvol/whitelist.txt
DEBUG:                  0
```


#### Custom configuration

If the default values are not feasible for your configuration it's possible to configure each of them. 

To do so, install the plugin but don't enable it (append a `--disable` to the install command), then set the variables and then finally enable it.


```sh
# Write some whitelisted volumegroups to a file
# in a custom location
echo "volumegroup0" | sudo tee "/mnt/somewhere/blabla.txt"

# Install the plugin but don't enable it
docker plugin install \
        --disable \
        --grant-all-permissions \
        cirocosta/golvm

# Enable 'debug' log level
docker plugin set \
        cirocosta/golvm \
        DEBUG=1

# Set the root of the volume mounts to be '/somewhere'
# instead of /mnt/lvmvol/volumes
docker plugin set \
        cirocosta/golvm \
        VOLUME_MOUNT_ROOT=/somewhere

# Set the path of the whitelist file
docker plugin set \
        cirocosta/golvm \
        WHITELIST_FILE=/mnt/somewhere/blabla.txt

# Enable the plugin
docker plugin enable \
        cirocosta/golvm

# Check whether everything went fine
docker plugin ls
ID                  NAME                DESCRIPTION                           ENABLED
84628b54dea6        lvmvol:latest       Docker plugin to manage LVM volumes   true
```

### Usage

#### Create regular volume

```sh
docker volume create \
        --driver lvmvol \
        --opt size=10M \
        myvol
```

#### Create thinly provisioned volume

#### Create snapshot volume

#### Create thin snapshot volume

#### List volumes

```sh
docker volume ls
``` 

#### Inspect volume

## lvmctl

`lvmctl` is a side utility that eases the process of managing the LVM volumes. It's only needed for performing actions that can't be covered by Docker's plugin semantics.

```
lvmctl --help
NAME:
   lvmctl - Controls the 'golvm' volume plugin

USAGE:
   lvmctl [global options] command [command options] [arguments...]

VERSION:
   master-dev

COMMANDS:
     check    checks verifies the environment
     create   create an LVM volume
     get      inspects existing LVM volumes
     ls       lists existing LVM volumes
     rm       removes existing LVM volumes
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
``` 

### Create regular volume

```sh
lvmctl create \
	--size=10M \
	vol
```

### Dependencies

Differently from the Docker plugin, `lvmctl` requires some dependencies. These are:

- lvm2 (for commands like `pvs`, `lvs`,`vgs`, `lvcreate`... )
- cryptsetup (for `luks` encryption)
- util-linux  (for `lsblk`)
- e2fsprogs  (for `mkfs.ext4`)
- xfsprogs (for `mkfs.xfs`)


### Usage

To make use of `lvmctl` make sure you have the right privileges. The same privileges needed for `pvs` are applicable for `lvmctl`. 

The following examples are all ran from a privileged user.

#### Create thin volume

```sh
lvcreate \
        --size 20M \
        --thin \
        volgroup0/thinpool0

lvmctl create \
	--size=10M \
	--thinpool=mythinpool \
	thin_vol
```

#### Create snapshot


```sh
lvmctl create \
	--size=100M \
	--snapshot=foobar \
	foobar_snap
```

#### Create thin snapshot

```sh
lvmctl create \
	--snapshot=foobar \
	foobar_snap
```

#### Encrypted volume

```sh
lvmctl \
	--size=0.2G \
	--keyfile=/root/key.bin \
	crypt_vol
```

ps.: Snapshots of encrypted volumes use the same key file. The key file must be present when the volume is created, and when it is mounted to a container.

