# Blockbook Build Guide

## Setting up your development environment

Supported environment to develop Blockbook is Linux. Although it is possible build and run Blockbook on macOS
or Windows our build process is not prepared for it. But you can still build Blockbook [manually](#manual-build).

The only dependency required to build Blockbook is Docker. You can see how to install Docker [here](https://docs.docker.com/install/linux/docker-ce/debian/).
Manual build require additional dependencies that are described in appropriate section.

## Build in Docker environment

All build operations run in Docker container in order to keep build environment isolated. Makefile in root of repository
define few targets used for building, testing and packaging of Blockbook. With Docker image definitions and Debian
package templates in *build/docker* and *build/templates* respectively, they are only inputs that make build process.

Docker build images are created at first execution of Makefile and that information is persisted. (Actually there are
created two files in repository – .bin-image and .deb-image – that are used as tags.) Sometimes it is necessary to
rebuild Docker images, it is possible by executing `make build-images`.

### Building binary

Just run `make` and that is it. Output binary is stored in *build* directory. Note that although Blockbook is Go application
it is dynamically linked with RocksDB dependencies and ZeroMQ. Therefore the operating system where Blockbook will be
executed still need that dependencies installed. See [Manual build](#manual-build) instructions below or install
Blockbook via Debian packages.

### Building debug binary

Standard binary contains no debug symbols. Execute `make build-debug` to get binary for debugging.

### Testing

How to execute tests is described in separate document [here](/docs/testing.md).

### Building Debian packages

Blockbook and particular coin back-end are usually deployed together. They are defined in the same place as well.
So typical way to build Debian packages is build Blockbook and back-end deb packages by single command. But it is not
mandatory, of course.

> Early releases of Blockbook weren't so friendly for extending. One had to define back-end package, Blockbook package,
> back-end configuration and Blockbook configuration as well. There were many options that were duplicated across
> configuration files and therefore error prone.
>
> Actually all configuration options and also build options for both Blockbook and backend are defined in single JSON
> file and all stuff required during build is generated dynamically.

Makefile targets follow simple pattern, there are few prefixes that define what to build.

* *deb-blockbook-<coin>* – Build Blockbook package for given coin.

* *deb-backend-<coin>* – Build back-end package for given coin.

* *deb-<coin>* – Build both Blockbook and back-end packages for given coin.

* *all-<coin>* – Similar to deb-<coin> but clean repository and rebuild Docker image before package build. It is useful
  for production deployment.

* *all* – Build both Blockbook and back-end packages for all coins.

Which coins are possible to build is defined in *configs/coins*. Particular coin has to have JSON config file there.

For example we want to build some packages for Bitcoin and Bitcoin Testnet.

```bash
# make all-bitcoin deb-backend-bitcoin_testnet
...
# ls build/*.deb
build/backend-bitcoin_0.16.1-satoshilabs-1_amd64.deb  build/backend-bitcoin-testnet_0.16.1-satoshilabs-1_amd64.deb  build/blockbook-bitcoin_0.0.6_amd64.deb
```

We have built two backend packages – for Bitcoin and Testnet – and Blockbook package for Bitcoin. Before build have been
performed there was cleaned build directory and rebuilt Docker image.

### Common notes

There are few variables that can be passed to make in order to modify build process.

In general, build of Blockbook binary require some dependencies. They are downloaded automatically during build process
but if you need to build the binary repeatedly it consumes a lot of time. Here comes variable *UPDATE_VENDOR* that if is
unset says that build process uses *vendor* (i.e. dependencies) from your local repository. For example:
`make deb-bitcoin UPDATE_VENDOR=0`. But before the command is executed there must be *vendor* directory populated,
you can do it by calling `dep ensure --vendor-only`. See [Manual build](#manual-build) instructions below.

All build targets allow pass additional parameters to underlying command inside container. It is possible via ARGS
variable. For example if you want run only subset of unit-tests, you will perform it by calling:
`make test ARGS='-run TestBitcoinRPC' UPDATE_VENDOR=0`

Common behaviour of Docker image build is that build steps are cached and next time they are executed much faster.
Although this is a good idea, when something went wrong you will need to override this behaviour somehow. This is
the command: `make build-images NO_CACHE=true`.

## Manual build

Instructions below are focused on Debian 9 (Stretch). If you want to use another Linux distribution or operating system
like macOS or Windows, please read instructions specific for each project.

Setup go environment:

```
wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz && tar xf go1.10.3.linux-amd64.tar.gz
sudo mv go /opt/go
sudo ln -s /opt/go/bin/go /usr/bin/go
# see `go help gopath` for details
mkdir $HOME/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Install RocksDB: https://github.com/facebook/rocksdb/blob/master/INSTALL.md
and compile the static_lib and tools

```
sudo apt-get update && sudo apt-get install -y \
    build-essential git wget pkg-config libzmq3-dev libgflags-dev libsnappy-dev zlib1g-dev libbz2-dev liblz4-dev
git clone https://github.com/facebook/rocksdb.git
cd rocksdb
CFLAGS=-fPIC CXXFLAGS=-fPIC make release
```

Setup variables for gorocksdb: https://github.com/tecbot/gorocksdb

```
export CGO_CFLAGS="-I/path/to/rocksdb/include"
export CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4"
```

Install ZeroMQ: https://github.com/zeromq/libzmq

Install go-dep tool:
```
go get github.com/golang/dep/cmd/dep
```

Get blockbook sources, install dependencies, build:

```
cd $GOPATH/src
git clone https://github.com/trezor/blockbook.git
cd blockbook
dep ensure -vendor-only
go build
```

### Example command

Blockbook require full node daemon as its back-end. You are responsible for proper installation. Port numbers and
daemon configuration are defined in *configs/coins* and *build/templates/backend/config* directories. You should use
specific installation process for particular coin you want run (e.g. https://bitcoin.org/en/full-node#other-linux-distributions for Bitcoin).

When you have running back-end daemon you can start Blockbook. It is highly recomended use ports described in [ports.md](/docs/ports.md)
for both Blockbook and back-end daemon. You can use *contrib/scripts/build-blockchaincfg.sh* that will generate
Blockbook's blockchain configuration from our coin definition files.

Example for Bitcoin:
```
contrib/scripts/build-blockchaincfg.sh
./blockbook -sync -blockchaincfg=build/blockchaincfg.json -internal=:9030 -public=:9130 -certfile=server/testcert -logtostderr
```

This command starts Blockbook with parallel synchronization and providing HTTP and Socket.IO interface, with database
in local directory *data* and established ZeroMQ and RPC connections to back-end daemon specified in configuration
file passed to *-blockchaincfg* option.

Blockbook logs to stderr (option *-logtostderr*) or to directory specified by parameter *-log_dir* . Verbosity of logs can be tuned
by command line parameters *-v* and *-vmodule*, for details see https://godoc.org/github.com/golang/glog.

You can check that Blockbook is running by simple HTTP request: `curl https://localhost:9130`. Returned data is JSON with some
run-time information. If port is closed, Blockbook is syncing data.
