# Blockbook Contributor Guide

Blockbook is a back-end service for Trezor wallet. Although it is open source, design and development of core packages
is done by Trezor developers in order to keep Blockbook compatible with Trezor. If you feel you could use Blockbook
for another purposes, we recommend you to make a fork.

However you can still help us find bugs or add support for new coins.

## Development environment

Instructions to set up your development environment and build Blockbook are described in separated
[document](/docs/build.md).

## How can I contribute?

### Reporting bugs

### Adding coin support

Trezor harware wallet supports over 500 coins, see https://trezor.io/coins/. You are free to add support for any of
them to Blockbook. Actually implemented coins are listed [here](/docs/ports.md).

You should follow few steps bellow to get smooth merge of your PR.

> Altough we are happy for support of new coins we have not enough capacity to run them all on our infrastructure.
> Actually we can run Blockbook instances only for coins supported by Trezor wallet. If you want to have Blockbook
> instance for your coin, you will have to deploy your own server.

#### Add coin definition

Coin definitions are stored in JSON files in *configs/coins* directory. They are single source of Blockbook
configuration, Blockbook and back-end package definition and build meta-data. Since Blockbook supports only single
coin index per running instance, every coin (including testnet) must have single definition file.

Because most of coins are fork of Bitcoin and they have similar way to install and configure their daemon, we use
templates to generate package definition and configuration files during the build process. It is similar to build
Blockbook package too. Templates are filled with data from the coin definition. Although the build process generate
packages automatically, there is sometimes necessary see an intermediate step. You can generate all files by calling
`go run build/templates/generate.go coin` where *coin* is name of definition file without .json extension. Files are
generated to *build/pkg-defs* directory.

Sections of the coin definition are described bellow. Good examples are
[*configs/coins/bitcoin.json*](configs/coins/bitcoin.json) and
[*configs/coins/ethereum.json*](configs/coins/ethereum.json) for Bitcoin-like coins and different coins, respectively.

* `coin` – Base information about coin.
    * `name` – Name of coin used internally (e.g. "Bcash Testnet").
    * `shortcut` – Ticker symbol (code) of coin (e.g. "TBCH").
    * `label` – Name of coin used publicly (e.g. "Bitcoin Cash Testnet").
    * `alias` – Name of coin used in file system paths and config files. We use convention that the name uses
      lowercase characters and underscore '_' as a word delimiter. Testnet versions of coins must have *_testnet*
      suffix. For example "bcash_testnet".

* `ports` – List of ports used by both back-end and Blockbook. Ports defined here are used in configuration templates
   and also as a source for generated documentation.
    * `backend_rpc` – Port of back-end RPC that is connected by Blockbook service.
    * `backend_message_queue` – Port of back-end MQ (if used) that is connected by Blockbook service.
    * `backend_*` – Additional back-end ports can be documented here. Actually the only purpose is to get them to
       the port table (prefix is removed and rest of string is used as note).
    * `blockbook_internal` – Blockbook's internal port that is used for metric collecting, debugging etc.
    * `blockbook_public` – Blockbook's public port that is used to comunicate with Trezor wallet (via Socket.IO).

* `ipc` – Defines how Blockbook connects its back-end service.
    * `rpc_url_template` – Template that defines URL of back-end RPC service. See note on templates below.
    * `rpc_user` – User name of back-end RPC service, used by both Blockbook and back-end configuration templates.
    * `rpc_pass` – Password of back-end RPC service, used by both Blockbook and back-end configuration templates.
    * `rpc_timeout` – RPC timeout used by Blockbook.
    * `message_queue_binding_template` – Template that defines URL of back-end's message queue (ZMQ), used by both
       Blockbook and back-end configuration template. See note on templates below.

* `backend` – Definition of back-end package, configuration and service.
    * `package_name` – Name of package. See convention note above.
    * `package_revision` – Revision of package. It specifies the version of the back-end package based on the upstream
       version. Usually changes only when there is a change in configuration. (For details about versioning see
       [Debian policy](https://www.debian.org/doc/debian-policy/ch-controlfields.html#version).)
    * `system_user` – User used to run back-end service. See convention note above.
    * `version` – Upstream version. Generated package have version in format:
       *&lt;backend.version&gt;-&lt;backend.package_revision&gt;*. (For details about versioning see
       [Debian policy](https://www.debian.org/doc/debian-policy/ch-controlfields.html#version).)
    * `binary_url` – URL of back-end archive. See note on back-end package build.
    * `verification_type` – Type of back-end archive verification. Possible values are *gpg*, *gpg-sha256*, *sha256*.
       See note on back-end package build.
    * `verification_source` – Source of sign/checksum of back-end archive. See note on back-end package build.
    * `extract_command` – Command to extract back-end archive. It is required to extract content of the archive to the
       *backend* directory.
    * `exclude_files` – List of files from back-end archive to exclude. Some files are not necessary for server
       deployment, some binaries have unnecessary dependencies, so it is good idea to extract these files from output
       package. Note that paths are relative to the *backend* directory where the archive is extracted.
    * `exec_command_template` – Template of command to execute back-end node daemon. Every back-end node daemon has its
       service that is managed by systemd. The template is evaluated to *ExecStart* option in *Service* section of
       service unit. See note on templates below.
    * `logrotate_files_template` – Template that define log files rotated by logrotate daemon. See note on templates
       below.
    * `postinst_script_template` – Additional steps in postinst script. See [ZCash definition](configs/coins/zcash.json)
       for more information.
    * `service_type` – Type of service. Services that daemonize must have *forking* type and write their PID to
       *PIDFile*. Services that don't support daemonization must have *simple* type. See examples above.
    * `service_additional_params_template` – Additional parameters in service unit. See
       [ZCash definition](configs/coins/zcash.json) for more information.
    * `protect_memory` – Enables *MemoryDenyWriteExecute* option in service unit if *true*.
    * `mainnet` – Set *false* for testnet back-end.
    * `config_file` – Name of template of back-end configuration file. Templates are defined in *build/backend/config*.
       For Bitcoin-like coins it is not necessary to add extra template, most options can be added via
       *additional_params*. For coins that don't require configuration the option should be empty (e.g. Ethereum).
    * `additional_params` – Object of extra parameters that are added to back-end configuration file as key=value pairs.
       Exception is *addnode* key that contains list of nodes that is expanded as addnode=item lines.

* `blockbook` – Definition of Blockbook package, configuration and service.
    * `package_name` – Name of package. See convention note above.
    * `system_user` – User used to run Blockbook service. See convention note above.
    * `internal_binding_template` – Template for *-internal* parameter. See note on templates below.
    * `public_binding_template` – Template for *-public* parameter. See note on templates below.
    * `explorer_url` – URL of blockchain explorer.
    * `additional_params` – Additional params of exec command (see [Dogecoin definition](configs/coins/dogecoin.json)).
    * `block_chain` – Configuration of BlockChain type that ensures communication with back-end service. All options
       must be tweaked for each individual coin separely.
        * `parse` – Use binary parser for block decoding if *true* else call verbose back-end RPC method that returns
           JSON. Note that verbose method is slow and not every coin support it. However there are coin implementations
           that don't support binary parsing (e.g. ZCash).
        * `mempool_workers` – Number of workers for UTXO mempool.
        * `mempool_sub_workers` – Number of subworkers for UTXO mempool.
        * `block_addresses_to_keep` – Number of blocks that are to be kept in blockaddresses column.
        * `additional_params` – Object of coin-specific params.

* `meta` – Common package metadata.
    * `package_maintainer` – Full name of package maintainer.
    * `package_maintainer_email` – E-mail of package maintainer.


TODO:
* name conventions baliky + uzivatele - nize
* prepsat system_user: neco jako mame konvenci, prefix, coin alias, s pomlckama, ale testnety suffix _testnet nepouzivaji
* Go template evaluation note
* verification, generate ports.md, download backend, path conventions, BB versioning
* hele dopsat tam, ze vetsinou musi stejne zmenit jen porty a url archivu a verification*, balicky, uzivatele, verzi

#### Add coin implementation

#### Add tests

#### Deploy public server
