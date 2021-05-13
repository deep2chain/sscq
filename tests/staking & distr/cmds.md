### account
    sscli accounts new 12345678
### query
    sscli query account $(sscli accounts list| sed -n '2p')
### bank
    sscli tx send $(sscli accounts list| sed -n '1p') $(sscli accounts list| sed -n '2p') 20000000satoshi --gas-price=100
    sscli tx create $(sscli accounts list| sed -n '1p') $(sscli accounts list| sed -n '2p') 1000satoshi --gas-price=100 --encode=false
    sscli tx sign [rawcode] --encode=false
### staking
    [validating]
    sscli tx staking create-validator $(sscli accounts list| sed -n '1p') \
                                       --pubkey=$(ssd tendermint show-validator)\
                                       --amount=100000000satoshi \
                                       --moniker=client \
                                       --commission-rate=0.10 \
                                       --commission-max-rate=0.20 \
                                       --commission-max-change-rate=0.01 \
                                       --min-self-delegation=1 \
                                       --gas-price=100
    sscli tx staking edit-validator $(sscli accounts list| sed -n '2p') --gas-price=100

    [delegating]
    sscli tx staking delegate $(sscli accounts list| sed -n '2p') <validatorAddress> 10000000satoshi --gas-adjustment=1.5 --gas-price=100