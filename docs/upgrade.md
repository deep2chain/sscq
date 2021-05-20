## steps
```
[1. stop ssd]
[2. backup .ssd]
[3. update ssd]
[4. start ssd]

[5. submit proposal]
sscli tx gov submit-proposal sscq1sh8d3h0nn8t4e83crcql80wua7u3xtlft9sr5d --gas-price=100  --switch-height=320 --description="second proposal"  --title="test1" --type="software_upgrade" --deposit="1000000000satoshi" --version="1"

[6. vote]
sscli tx send sscq1sh8d3h0nn8t4e83crcql80wua7u3xtlft9sr5d sscq1xvktg68uwrtkml7m4yqul2sjm6fhtglm6y20eg 1000000000satoshi --gas-price=100
sscli tx send sscq1sh8d3h0nn8t4e83crcql80wua7u3xtlft9sr5d sscq1p97p8vckpvkzx34se7eansm3rp223r3trmn63h 1000000000satoshi --gas-price=100
sscli tx send sscq1sh8d3h0nn8t4e83crcql80wua7u3xtlft9sr5d sscq18efa8m8yudlqc765s9vkdrev5475jyqalqm5a0 1000000000satoshi --gas-price=100
sscli tx send sscq1sh8d3h0nn8t4e83crcql80wua7u3xtlft9sr5d sscq1nuxf4amphaajuwg0ph3se6kmsda9cs6v0sja7r 1000000000satoshi --gas-price=100
sscli tx gov vote  sscq1xvktg68uwrtkml7m4yqul2sjm6fhtglm6y20eg 2  yes --gas-price=100 
sscli tx gov vote  sscq1p97p8vckpvkzx34se7eansm3rp223r3trmn63h 2  yes --gas-price=100 
sscli tx gov vote  sscq18efa8m8yudlqc765s9vkdrev5475jyqalqm5a0 2  yes --gas-price=100 
sscli tx gov vote  sscq1nuxf4amphaajuwg0ph3se6kmsda9cs6v0sja7r 2  yes --gas-price=100 

[7. check]
sscli query staking params
unbonding, unslashing test

```