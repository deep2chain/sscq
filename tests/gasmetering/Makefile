# written by junying
# 2019.10.28
# index
# 1.sscqservice tx
# - non-contract tx
# - contract tx
# 2.non-contract tx

TMP_PATH = /tmp/result.json
BLK_TIME = 6
# sscqservice tx
## non-contract tx
### test cases
#### gas < DefaultTxGas(21000)
#### gasprice < 1satoshi
TX_AMOUNT = 1000000
UNIT = satoshi
TxGas = 30000
GAS_PRICE = 100
test.normal:
	@acc1=$$(sscli accounts list|row 1);\
	 acc2=$$(sscli accounts list|row 2);\
	 sscli query account $$acc1 > ${TMP_PATH};\
	 orgn1=$$(findkey amount ${TMP_PATH});\
	 sscli query account $$acc2 > ${TMP_PATH};\
	 orgn2=$$(findkey amount ${TMP_PATH});\
	 read -p "GasWanted:" GasWanted;\
	 sscli tx send $$acc1 $$acc2 ${TX_AMOUNT}${UNIT} --gas-wanted=$$GasWanted --gas-price=${GAS_PRICE} > ${TMP_PATH};\
	 txhash=$$(findstr txhash ${TMP_PATH}|fromstr ": "|extractstr '"' '"');\
	 sleep ${BLK_TIME};\
	 sscli query tx $$txhash > ${TMP_PATH};\
	 success=$$(findkey success ${TMP_PATH});\
	 gas_wanted=$$(findkey gas_wanted ${TMP_PATH});\
	 gas_used=$$(findkey gas_used ${TMP_PATH});\
	 GasPrice=$$(findkey GasPrice ${TMP_PATH});\
	 Gas=$$(findkey Gas ${TMP_PATH});\
	 sscli query account $$acc1 > ${TMP_PATH};\
	 curr1=$$(findkey amount ${TMP_PATH});\
	 sscli query account $$acc2 > ${TMP_PATH};\
	 curr2=$$(findkey amount ${TMP_PATH});\
	 margin=$$(python -c "print $$orgn1-$$curr1");\
	 python -c "result='passed!' if $$margin == ($$GasPrice*$$gas_used+${TX_AMOUNT}) else 'failed!'; print result";\
	 echo "$$orgn1 - $$curr1 == $$margin == $$gas_used * $$GasPrice + ${TX_AMOUNT}";


## contract tx
### test cases
#### 1. gas > gasUsed
#### 2. Intrinsic Gas < gas < gasUsed: DeliverTx Success, No Error, No gas_used report
#### 3. gas < Intrinsic : DeliverTx Failed, create contract error
#### 4. gas < GasMeter.ConsumedGas(): CheckTx Failed
test.contract.creation:
	@acc1=$$(sscli accounts list|row 1);\
	 sscli query account $$acc1 > ${TMP_PATH};\
	 orgn1=$$(findkey amount ${TMP_PATH});\
	 echo $$(cat ../evm/coin/coin_sol_Coin.bin);\
	 $(MAKE) -sC ../evm create.contract.default;\
	 txhash=$$(findkey txhash ${TMP_PATH});\
	 success=$$(findkey success ${TMP_PATH});\
	 gas_wanted=$$(findkey gas_wanted ${TMP_PATH});\
	 gas_used=$$(findkey gas_used ${TMP_PATH});\
	 gas_used=$$(python -c "gas_used='0' if '$$gas_used'== 'no found!' else '$$gas_used'; print gas_used");\
	 GasPrice=$$(findkey GasPrice ${TMP_PATH});\
	 Gas=$$(findkey Gas ${TMP_PATH});\
	 sscli query account $$acc1 > ${TMP_PATH};\
	 curr1=$$(findkey amount ${TMP_PATH});\
	 margin=$$(python -c "print $$orgn1-$$curr1");\
	 python -c "result='passed!' if $$margin == $$GasPrice*$$gas_used else 'failed!'; print result";\
	 echo "$$orgn1 - $$curr1 == $$margin == $$gas_used * $$GasPrice?";

# non-sscqservice tx
# staking unbond testing failed
# after-unbounding balance = before-balance
UNBOND_AMOUNT = 100000

GENESIS_JSON_PATH = ~/.ssd/config/genesis.json
test.staking.unbond:
	@acc=$$(sscli accounts list|row 1);\
	 sscli query account $$acc > ${TMP_PATH};\
	 orgn=$$(findkey amount ${TMP_PATH});\
	 sscli query staking validators > ${TMP_PATH};\
	 tokens=$$(findkey tokens ${TMP_PATH});\
	 echo balance:$$orgn;\
	 echo tokens:$$tokens;\
	 val_addr=$$(findkey validator_address ${GENESIS_JSON_PATH});\
	 sscli query distr rewards $$acc > ${TMP_PATH};\
	 rewards=$$(findstr amount ${TMP_PATH}|fromstr ": "|extractstr '"' '"');\
	 echo rewards:$$rewards;\
	 sscli tx staking unbond $$acc $$val_addr ${UNBOND_AMOUNT}${UNIT} --gas-adjustment 1.5 --gas-price=${GAS_PRICE} > ${TMP_PATH};\
	 txhash=$$(findstr txhash ${TMP_PATH}|fromstr ": "|extractstr '"' '"');\
	 echo txhash:$$txhash;\
	 sleep ${BLK_TIME};\
	 sscli query tx $$txhash > ${TMP_PATH};\
	 success=$$(findkey success ${TMP_PATH});\
	 echo success:$$success;\
	 if [ $$success = "False" ]; then echo success:$$success please press ctrl+c to exit; sleep 1000;fi;\
	 gas_wanted=$$(findkey gas_wanted ${TMP_PATH});\
	 echo gas_wanted:$$gas_wanted;\
	 gas_used=$$(findkey gas_used ${TMP_PATH});\
	 echo gas_used:$$gas_used;\
	 sscli query account $$acc > ${TMP_PATH};\
	 curr=$$(findkey amount ${TMP_PATH});\
	 margin=$$(python -c "print $$curr-$$orgn");\
	 fee=$$(python -c "print int($$gas_used*${GAS_PRICE})+1");\
	 python -c "result='passed!' if $$margin == $$rewards+${UNBOND_AMOUNT}-$$fee else 'failed!'; print result";\
	 echo $$curr - $$orgn == $$margin == $$rewards+${UNBOND_AMOUNT}-$$fee?;

# 1. create a new account
# 2. delegate with the account
DELEGATE_AMOUNT = 100000
PASSWORD = 12345678
test.staking.delegate:
	@acc=$$(sscli accounts new ${PASSWORD});\
	 sscli tx send $$(sscli accounts list|row 1) $$acc 10000000${UNIT} --gas-price=${GAS_PRICE};\
	 sleep ${BLK_TIME};\
	 sscli query account $$acc > ${TMP_PATH};\
	 orgn=$$(findkey amount ${TMP_PATH});\
	 echo balance:$$orgn;\
	 val_addr=$$(findkey validator_address ${GENESIS_JSON_PATH});\
	 sscli tx staking delegate $$acc $$val_addr ${DELEGATE_AMOUNT}${UNIT} --gas-adjustment 1.0 --gas-price=${GAS_PRICE} > ${TMP_PATH};\
	 txhash=$$(findstr txhash ${TMP_PATH}|fromstr ": "|extractstr '"' '"');\
	 echo txhash:$$txhash;\
	 sleep ${BLK_TIME};\
	 sscli query tx $$txhash > ${TMP_PATH};\
	 success=$$(findkey success ${TMP_PATH});\
	 echo success:$$success;\
	 if [ $$success = "False" ]; then echo success:$$success please press ctrl+c to exit; sleep 1000;fi;\
	 gas_wanted=$$(findkey gas_wanted ${TMP_PATH});\
	 echo gas_wanted:$$gas_wanted;\
	 gas_used=$$(findkey gas_used ${TMP_PATH});\
	 echo gas_used:$$gas_used;\
	 sscli query account $$acc > ${TMP_PATH};\
	 curr=$$(findkey amount ${TMP_PATH});\
	 margin=$$(python -c "print $$orgn-$$curr");\
	 fee=$$(python -c "print int($$gas_used*${GAS_PRICE})");\
	 python -c "result='passed!' if $$margin == ${DELEGATE_AMOUNT}+$$fee else 'failed!'; print result";\
	 echo  $$orgn - $$curr == $$margin == ${DELEGATE_AMOUNT}+$$gas_used*${GAS_PRICE}?;
