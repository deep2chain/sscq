### minimum gas price
    ~/.hsd/config/hsd.toml
    params/fee.go - DefaultMinGasPrice
    init/testnet.go - FlagMinGasPrices
### persistant peers
    ~/.hsd/config/config.toml
### gentxs
    ~/.hsd/config/genesis.json
### chain-id
    hscli config chain-id testchain
    hsd testnet --chain-id testchain
    hsd init [moniker] --chain-id testchain
    init/testnet.go - FlagChainID
### trust node
    hscli config trust-node true