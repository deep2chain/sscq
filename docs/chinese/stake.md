    delegator-addr: sscq1zf07fyt2an2ral8zve0u4y7lzqa6x4lqfeyl8m
    validator-addr: sscqvaloper1zf07fyt2an2ral8zve0u4y7lzqa6x4lqrquxss
    amount: 10000000stake
    samount: 100satoshi
### 设置chain-id
    sscli config chain-id [chain-id]
### 如果代表地址里面没钱的话，转账一笔
    sscli query accounts [delegator-addr]
    sscli tx send [fromaddr] [delegator-addr] [samount] --gas-price=100
### [stake-抵押，解绑](https://github.com/deep2chain/sscq/blob/master/x/staking/client/cli/tx.go)   
    [抵押]
    sscli tx staking delegate [delegator-addr] [validator-addr] [amount] --gas-price=100
    
    [解绑]
    sscli tx staking unbond [delegator-addr] [validator-addr] [amount]  --gas-price=100
### [奖励-查询，回收](https://github.com/deep2chain/sscq/blob/master/x/distribution/client/cli/tx.go)
    [查询]
    sscli query distr rewards [delegator-addr]

    [回收]
    sscli tx distr withdraw-rewards [delegator-addr] [validator-addr]  --gas-price=100
### [管理者，设置delegator解绑状态]
    [查询]
    sscli query staking delegation [delegator-addr] [validator-addr]

    [许可]
    sscli tx staking uds [delegator-addr] --delegator-manager=[validator-addr] --delegator-status=true --gas-price=100
