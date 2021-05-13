# junying-todo, 2019-08-06
# constructor is called when contract created. you have to add constructor argument to the end of contract bytecode,
# then combined bytecode is used to create a contract.
# it will call constructor with arguments when creating a contract.

# contract tx hash: 2861E1B659A72195DFDC25FBDF97EE7A17BBA7557051AC6BCDEB7ACB756A6B26
# contract addr: sscq1guam0nv5ky049g6zg05qtx9y3ayv0v3pvaavu9
# issuer:		 sscq1ha7ryup8nc2avgesfunx2pm22waqv2cx6dj0ac
# bid
# 0x1998aeef
#---------------------
# beneficiary
# 0x38af3eed
#---------------------
# withdraw
# 0x3ccfd60b
#---------------------
# auctionEnd
# 0x2a24f46c
#---------------------
# highestBidder
# 0x91f90157
#---------------------
# highestBid
# 0xd57bde79
#---------------------
# remaining
# 0x55234ec0
#---------------------
calc.method.id:
	@make -sC ../ get.method.id

param.address:
	@make -sC ../ param.address
	
param.int:
	@make -sC ../ param.int

run.contract:
	@read -p "data: " data;\
	 replkey Data $$data unsigned.tx y;\
	 contract=$$(oneline unsigned.tx | excludestr "    ");\
	 echo "############################################################";\
	 echo $$contract; \
	 contract_hex=$$(ssutils json2hex "$$contract");\
	 echo $$contract_hex;

BLK_TIME = 7

TMP_PATH = /tmp/result.json

execute:
	@contract=$$(oneline unsigned.tx | excludestr "    ");\
	 echo "############################################################";\
	 echo $$contract; \
	 unsigned=$$(ssutils json2hex "$$contract");\
	 echo $$unsigned;\
	 signed=$$(sscli tx sign $$unsigned);\
	 echo $$signed;\
	 sscli tx broadcast $$signed > ${TMP_PATH};\
	 txid=$$(findkey txhash ${TMP_PATH});\
	 echo $$txid;\
	 sleep ${BLK_TIME};\
	 sscli query tx $$txid > ${TMP_PATH};\
	 findkey logs ${TMP_PATH};

# [PREASSUMPTION]
# You should set static items such as To(Contract Address), Gas, GasPrice manually first
# so that you can use unsigned template with the following few changes.
# In:  []
# Out: []
# test: passed
bid:
	@read -p "bidder address: " bidder;\
	 read -p "bid amount: " amount;\
	 replkey From $$bidder bid.tx y;\
	 replvalue amount 0 $$amount bid.tx y;\
	 replkey Data 1998aeef bid.tx y;\
	 contract=$$(oneline bid.tx | excludestr "    ");\
	 echo "############################################################";\
	 echo $$contract; \
	 unsigned=$$(ssutils json2hex "$$contract");\
	 echo $$unsigned;\
	 signed=$$(sscli tx sign $$unsigned);\
	 echo $$signed;\
	 sscli tx broadcast $$signed > result.json;\
	 txid=$$(findkey txhash result.json);\
	 echo $$txid;\
	 sleep ${BLK_TIME};\
	 sscli query tx $$txid > result.json;\
	 findkey logs result.json;\
	 replvalue amount $$amount 0 bid.tx y;

# In:  []
# Out: address
# error case: open contract error when data involved '0x'
# test: passed
# result: 000000000000000000000000eb51f19caa7c653edeca6e78644c0f772f2c131d
beneficiary:
	@replkey From $$(sscli accounts list|row 1) unsigned.tx y;\
	 replkey Data 38af3eed unsigned.tx y;\
	 make -sC . execute;


# In:  []
# Out: address
# test: passed
# result: 0000000000000000000000002132306e77dd63d06e648e98abc61cfa8844b951
highestBidder:
	@replkey From $$(sscli accounts list|row 1) unsigned.tx y;\
	 replkey Data 91f90157 unsigned.tx y;\
	 make -sC . execute;

# In:  []
# Out: uint256
# test: passed
# result: 0000000000000000000000000000000000000000000000000000000000989680
highestBid:
	@replkey From $$(sscli accounts list|row 1) unsigned.tx y;\
	 replkey Data d57bde79 unsigned.tx y;\
	 make -sC . execute;

# In: []
# Out: []
# YOU NEED TO CALL THIS TO FINALIZE AUCTION. THEN, THE BENEFICIARY WILL RECEIVE THE HIGHTEST BID.
auctionEnd:
	@replkey From $$(sscli accounts list|row 1) unsigned.tx y;\
	 replkey Data 2a24f46c unsigned.tx y;\
	 make -sC . execute;

# In: []
# Out: [bool]
# 
withdraw:
	@replkey From $$(sscli accounts list|row 1) unsigned.tx y;\
	 replkey Data 3ccfd60b unsigned.tx y;\
	 make -sC . execute;

# function name: pending
# parameters: 'address'
# 0x5eebea20
pending: #param.address
	@queryaddr=$$(make -sC . param.address);\
	 replkey From $$(sscli accounts list|row 1) unsigned.tx y;\
	 replkey Data 5eebea20$$queryaddr unsigned.tx y;\
	 make -sC . execute;

# function name: remaining
# 0x55234ec0
remaining:
	@replkey From $$(sscli accounts list|row 1) unsigned.tx y;\
	 replkey Data 55234ec0 unsigned.tx y;\
	 make -sC . execute;

