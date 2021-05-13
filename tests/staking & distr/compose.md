# MainNet
### action
    make sstart
    make sstest
    make sstop
    sudo make clean
### sstest
    [alias]
    cp ss* /usr/local/bin
    alias sscli="sscli --home node4/.sscli"
    alias ssd="ssd --home node4/.ssd"

    [config/start]
    sscli config chain-id testchain
    sscli config trust-node true
    sscli config node http://192.168.10.2:26657
    # nohup ssd start & > /dev/null
    # clear

    [transactions]
    sscli accounts new 12345678
    sscli accounts list
    sscli query account $(sscli accounts list| sed -n '1p')
    sscli tx send $(sscli accounts list| sed -n '1p') $(sscli accounts list| sed -n '2p') 20000000satoshi --gas-price=100
    sscli tx send $(sscli accounts list| sed -n '1p') $(sscli accounts list| sed -n '2p') 20000000satoshi --gas-price=100

    [validators]
    - show yours
    ssd tendermint show-validator
    - check status
    sscli query staking validators
    sscli  query staking validator [sscqvaloper]
    - confirm running
    sscli query tendermint-validator-set
    sscli query tendermint-validator-set | grep [sscqvalcons/sscqvalconspub]
    sscli query tendermint-validator-set | grep "$(ssd tendermint show-validator)"
    - start yours(tip: 100,000,000 for voting power 1, 1,000,000,000 for 10, 10,000,000,000 for 100)
    sscli tx staking create-validator $(sscli accounts list| sed -n '2p') \
                                       --pubkey=$(ssd tendermint show-validator)\
                                       --amount=100000000satoshi \
                                       --moniker=client \
                                       --commission-rate=0.10 \
                                       --commission-max-rate=0.20 \
                                       --commission-max-change-rate=0.01 \
                                       --min-self-delegation=1 \
                                       --gas-price=100
    sscli tx staking edit-validator $(sscli accounts list| sed -n '2p') --gas-price=100
    or
    sscli tx staking edit-validator $(sscli accounts list| sed -n '2p')\
                --moniker=client \
                --chain-id=testchain \
                --website="https://sscq.network" \
                --identity=23870f5bb12ba2c4967c46db \
                --details="To infinity and beyond!" \
                --gas-price=100 \
                --commission-rate=0.10
    - unjail
    sscli tx slashing unjail $(sscli accounts list| sed -n '2p') --gas-price=100
    - log
    sscli query slashing signing-info [sscqvalconspub]
    - check

    [delegators]
    sscli tx staking delegate $(sscli accounts list| sed -n '2p') \
                               $(grep -nr validator_address  ~/.ssd/config/genesis.json |sed -n '1p'|awk '{print $3F}' | cut -d'"' -f 2)\
                               100000000satoshi --gas-adjustment=1.5 --gas-price=100
    sscli tx staking redelegate $(sscli accounts list| sed -n '2p') \
                                 $(grep -nr validator_address  ~/.ssd/config/genesis.json |sed -n '2p'|awk '{print $3F}' | cut -d'"' -f 2) \
                                 $(grep -nr validator_address  ~/.ssd/config/genesis.json |sed -n '3p'|awk '{print $3F}' | cut -d'"' -f 2) \
                                 100000000satoshi --gas-adjustment=1.5 --gas-price=100
    sscli query distr rewards $(sscli accounts list| sed -n '2p')
    sscli query account $(sscli accounts list| sed -n '2p')
    sscli tx staking unbond $(sscli accounts list| sed -n '1p') \
                             $(grep -nr validator_address  ~/.ssd/config/genesis.json |sed -n '1p'|awk '{print $3F}' | cut -d'"' -f 2) \
                             100000000satoshi \
                             --gas-adjustment 1.5 --gas-price=100
    sscli tx distr withdraw-rewards $(sscli accounts list| sed -n '1p') \
                                     $(grep -nr validator_address  ~/.ssd/config/genesis.json |sed -n '1p'|awk '{print $3F}' | cut -d'"' -f 2) \ --gas-adjustment 1.5 --gas-price=100