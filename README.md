# golvm

`golvm` provides both a library for dealing with LVM operations as well as a docker plugin (`lvmvol`).

## Plugin

### Plugin activation

Every volume managed by the plugin is mounted under `/mnt/lvmvol/volumes` by default (can be configured). Using the default configuration means that it's required that the directory `/mnt` exists before the plugin is enabled.


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

### Plugin usage

Once the plugin has been properly installed we can start using it.

1. (optional) Create a whitelist of volumegroups to be used by the plugin

```sh
echo "vgs1" >> /mnt/lvmvol/whitelist
```

2.      Install the plugin

```sh
docker plugin install \
        --grant-all-permissions \
        cirocosta/golvm
```

3.      Create a volume

```sh
docker volume create \
        --driver lvmvol \
        --opt size=10M \
        myvol
```

4.      List the volumes

```sh
docker volume ls
``` 

## `lvmctl` Usage

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

## Examples

1.	Create a volume named `vol` with size `10M`.

```sh
lvmctl create \
	--size=10M \
	vol
```


2. 	Create a thinly-provisioned lvm volume named `thin_vol` in `mythinpool`

```sh

sudo lvcreate \
        --size 20M \
        --thin \
        volgroup0/thinpool0

lvmctl create \
	--size=10M \
	--thinpool=mythinpool \
	thin_vol
```


3. 	Create a snapshot volume of `foobar` named `foobar_snap`. 

```sh
lvmctl create \
	--size=100M \
	--snapshot=foobar \
	foobar_snap
```


4.	Create a thin snapshot of `foobar` named `foobar_thin_snap` (same command as the normal snapshot but without `--size` option)

```sh
lvmctl create \
	--snapshot=foobar \
	foobar_snap
```

5.	Create a `LUKS` encrypted volume named `crypt_vol` with the contents of `/root/key.bin` as a binary passphrase. 

```sh
lvmctl \
	--size=0.2G \
	--keyfile=/root/key.bin \
	crypt_vol
```

ps.: Snapshots of encrypted volumes use the same key file. The key file must be present when the volume is created, and when it is mounted to a container.


## Docker Plugin

1.      Whitelist volume groups to be used

```
echo "myvg" >> /etc/docker/golvm
```


## TODO

- check if the LV will conflict prior to the creation

