### minimum gas price
    ~/.ssd/config/ssd.toml
    params/fee.go - DefaultMinGasPrice
    init/testnet.go - FlagMinGasPrices
### persistant peers
    ~/.ssd/config/config.toml
### gentxs
    ~/.ssd/config/genesis.json
### chain-id
    sscli config chain-id testchain
    ssd testnet --chain-id testchain
    ssd init [moniker] --chain-id testchain
    init/testnet.go - FlagChainID
### trust node
    sscli config trust-node true