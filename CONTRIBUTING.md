# Blockbook Contributor Guide

Blockbook is a back-end service for Trezor hardware wallet. Although it is open source, design and development of
core packages is done by Trezor developers in order to keep Blockbook compatible with Trezor. If you feel you can
use Blockbook for another purposes, it is recommended to make your own fork.

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

* `coin` –
    * `name` –
    * `shortcut` –
    * `label` –
    * `alias` –

* `ports` –
    * `backend_rpc` –
    * `backend_message_queue` –
    * `backend_*` –
    * `blockbook_internal` –
    * `blockbook_public` –

* `ipc` –
    * `rpc_url_template` –
    * `rpc_user` –
    * `rpc_pass` –
    * `rpc_timeout` –
    * `message_queue_binding_template` –

* `backend` –
    * `package_name` –
    * `package_revision` –
    * `system_user` –
    * `version` –
    * `binary_url` –
    * `verification_type` –
    * `verification_source` –
    * `extract_command` –
    * `exclude_files` –
    * `exec_command_template` –
    * `logrotate_files_template` –
    * `postinst_script_template` –
    * `service_type` –
    * `service_additional_params_template` –
    * `protect_memory` –
    * `mainnet` –
    * `config_file` –
    * `additional_params` –

* `blockbook` –
    * `package_name` –
    * `system_user` –
    * `internal_binding_template` –
    * `public_binding_template` –
    * `explorer_url` –
    * `additional_params` –
    * `block_chain` –
        * `parse` –
        * `mempool_workers` –
        * `mempool_sub_workers` –
        * `block_addresses_to_keep` –
        * `additional_params` –

* `meta` –
    * `package_maintainer` –
    * `package_maintainer_email` –

#### Add coin implementation

#### Add tests

#### Deploy public server
