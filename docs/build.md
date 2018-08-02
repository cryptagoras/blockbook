# Blockbook Build Guide

## Setting up your development environment

Supported environment to develop Blockbook is Linux. Although it is possible build and run Blockbook on macOS
or Windows our build system is not prepared for it. But you can still build Blockbook [manually](#manual-build)

The only dependency required to build Blockbook is Docker. You can see how to install Docker [here](https://docs.docker.com/install/linux/docker-ce/debian/).
Manual build require additional dependencies that are described in appropriate section.

## Build in Docker environment

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
