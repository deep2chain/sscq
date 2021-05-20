### configuration
    sscli config chain-id [chain-id]

### accounts cmds
    sscli accounts newaccount
    sscli accounts listaccounts
    sscli accounts genprivkey [addr]
    sscli accounts getbalance [addr]

### transaction cmds
    sscli tx send [fromaddr] [toaddr] [amount]
    sscli tx create [fromaddr] [toaddr] [amount]
    sscli tx sign [rawdata]
    sscli tx broadcast [rawdata]

### query cmds
```
sscli query accounts [addr]
sscli query block
sscli query txs
sscli query tx

[additional]
sscli query rewards [block-height]
sscli query total-provisions

[contractcall]
contract-addr: sscq1l03rqalmg58wgw9ya39wwc3466lyy20xnpeaee
callcode: 27e235e300000000000000000000000027681ceb7de9bae3c5f7f10f81ff5106e2ca48a7
sscli query contract [contract-addr] [callcode]
```
### check
    sscli query staking pool
    sscli query staking params
    sscli query distr params

### [staking cmds](https://github.com/deep2chain/sscq/blob/master/x/staking/client/cli/tx.go)
    delegator-addr: sscq1zf07fyt2an2ral8zve0u4y7lzqa6x4lqfeyl8m
    validator-addr: sscqvaloper1zf07fyt2an2ral8zve0u4y7lzqa6x4lqrquxss
    amount: 100000stake
    
    [unbound]
    sscli tx staking unbond [delegator-addr] [validator-addr] [amount] --gas-adjustment 1.5 --gas-price=100

    [delegate]
    sscli tx staking delegate [delegator-addr] [validator-addr] [amount] --gas-adjustment=1.5 --gas-price=100
### [rewards](https://github.com/deep2chain/sscq/blob/master/x/distribution/client/cli/tx.go)
    [query]
    sscli query distr rewards [delegator-addr]
    sscli query distr rewards <delegator_address> <validator_address>
    sscli query distr commission <validator_address>
    sscli query distr community-pool
    sscli query rewards 1

    [withdraw]
    sscli tx distr withdraw-rewards [delegator-addr] [validator-addr] --gas-adjustment 1.5 --gas-price=100
    sscli tx distr withdraw-rewards [delegator-addr] [validator-addr] --commission --gas-adjustment 1.5 --gas-price=100

### upgrade
```
[query]
sscli query gov proposal [proposal-id]
sscli query gov proposal 2
sscli query gov proposals
sscli query gov votes [proposal-id] 
sscli query gov votes 1

[submit]
sscli tx gov submit-proposal [flags]
sscli tx gov submit-proposal sscq1sh8d3h0nn8t4e83crcql80wua7u3xtlft9sr5d --gas-price=100  --switch-height=4100 --description="third proposal"  --title="test2" --type="software_upgrade" --deposit="1000000000satoshi" --version="1"

[vote]
sscli tx gov vote [voter-addr] [proposal-id] [option] [flags]
sscli tx gov vote  sscq148asterza2u7ww0vptntmy8ut84hdeetr927hl 3  yes --gas-price=100 

[deposit]
sscli tx gov deposit [proposal-id] [deposit] [flags]
```
### unjail
```
sscli tx slashing unjail [validator-address] --gas-price=100
sscli query staking validators|egrep -e "jail|status|token|share"
```