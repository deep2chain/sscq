### cmds
    [min-gas-prices]
    ssd start --minimum-gas-price=20

    [fee & gasprice]
    sscli tx send [fromaddr] [toaddr] [amount] --gas-price=100
    sscli tx send [fromaddr] [toaddr] [amount] --gas-wanted=30000 --gas-price=100
### references
#### [client cmd: fee & gas](https://cosmos.network/docs/gaia/sscli.html#fees-gas)
#### [gaiad.toml: minimum-gas-prices](https://cosmos.network/docs/gaia/join-mainnet.html#set-minimum-gas-prices)