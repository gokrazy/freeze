Little helper to create tar balls of an executable together with its ELF shared
library dependencies. This is useful for prototyping with gokrazy:
https://gokrazy.org/prototyping/

## Installation

```shell
go install github.com/gokrazy/freeze/cmd/...@latest
```

## Usage

Let’s assume you want to try the upcoming Linux ksmbd feature.

On Linux, build `ksmbd-tools`:

```shell
$ git clone https://github.com/cifsd-team/ksmbd-tools
$ cd ksmbd-tools
$ ./autogen.sh
$ ./configure 
$ make -j8
$ freeze control/ksmbd.control
[…]
2021/10/24 15:29:33 Download freeze1373262977.tar to your gokrazy device and run:
	LD_LIBRARY_PATH=$PWD ./ld-linux-x86-64.so.2 ./ksmbd.control
```

Then, on your gokrazy device, e.g. via [breakglass](https://github.com/gokrazy/breakglass):

```shell
$ cd /tmp
$ wget http://10.0.0.76:4080/freeze1373262977.tar
$ tar xf freeze1373262977.tar 
$ cd freeze1373262977/
$ LD_LIBRARY_PATH=$PWD ./ld-linux-x86-64.so.2 ./ksmbd.control
Usage: ksmbd.control
	-s | --shutdown
	-d | --debug=all or [smb, auth, etc...]
	-c | --ksmbd-version
	-V | --version
```
