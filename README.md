# golvm

## Plugin usage

1.      (optional) Create a whitelist of volumegroups to be used by the plugin

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

