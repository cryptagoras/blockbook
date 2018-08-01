# Registry of ports

| coin                     | blockbook internal port | blockbook public port | backend rpc port   | zmq port |
|--------------------------|-------------------------|-----------------------|--------------------|----------|
| Bitcoin                  | 9030                    | 9130                  | 8030               | 38330    |
| Bcash                    | 9031                    | 9131                  | 8031               | 38331    |
| Zcash                    | 9032                    | 9132                  | 8032               | 38332    |
| Dash                     | 9033                    | 9133                  | 8033               | 38333    |
| Litecoin                 | 9034                    | 9134                  | 8034               | 38334    |
| Bgold                    | 9035                    | 9135                  | 8035               | 38335    |
| Ethereum                 | 9036                    | 9136                  | 8036 ws, 8136 http | 38336*   |
| Ethereum Classic         | 9037                    | 9137                  | 8037               | 38337*   |
| Dogecoin                 | 9038                    | 9138                  | 8038               | 38338    |
| Namecoin                 | 9039                    | 9139                  | 8039               | 38339    |
| Vertcoin                 | 9040                    | 9140                  | 8040               | 38340    |
| Bitcoin Testnet          | 19030                   | 1913                  | 18030              | 48330    |
| Bcash Testnet            | 19031                   | 1913                  | 18031              | 48331    |
| Zcash Testnet            | 19032                   | 1913                  | 18032              | 48332    |
| Dash Testnet             | 19033                   | 1913                  | 18033              | 48333    |
| Litecoin Testnet         | 19034                   | 1913                  | 18034              | 48334    |
| Ethereum Testnet Ropsten | 19036                   | 19136                 | 18036              | 48336*   |
| Vertcoin Testnet         | 19040                   | 19140                 | 18040              | 48340    |

\* geth listens on this port, however not as zmq service

> NOTE: This document is generated from coin definitions in `configs/coins`.
