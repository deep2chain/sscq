# junying-todo, 2019-08-06

# contract tx hash: 141D3CD4EB7B0BDB41C2313CF2170F6FBF0410898BEB05FBB7D88BA182E5028B
#					58A1FD7C7FD874FB72E9776638DF4C9107EF263CC8401738AC2FD60AB5BC08B6
# contract addr: sscq1hx484kx6nntkv5ascuw20c5pn7a879agvmurh7
# issuer:		 sscq1ha7ryup8nc2avgesfunx2pm22waqv2cx6dj0ac
# query data:	 0x07546172
# log:			 "log": "{\"code\":0,\"message\":\"ok\",\"contract_address\":\"\",\"evm_output\":\"000000000000000000000000bf7c3270279e15d623304f2665076a53ba062b06\"}"
# result:		 ssutils hex2json bf7c3270279e15d623304f2665076a53ba062b06
#				 sscq1ha7ryup8nc2avgesfunx2pm22waqv2cx6dj0ac

calc.method.id:
	@make -sC ../ get.method.id

param.address:
	@make -sC ../ param.address
	
param.int:
	@make -sC ../ param.int

run.contract:
	@read -p "fromaddress: " fromaddr;\
	 read -p "contract address: " contract_addr;\
	 read -p "data: " data;\
	 gas=900000;\
	 gasprice=1;\
	 replkey From $$fromaddr ../unsigned.tx y;\
	 replkey To $$contract_addr ../unsigned.tx y;\
	 replkey Data $$data ../unsigned.tx y;\
	 replkey GasLimit $$gas ../unsigned.tx y;\
	 replkey GasPrice $$gasprice ../unsigned.tx y;\
	 contract=$$(oneline ../unsigned.tx | excludestr "    ");\
	 echo "############################################################";\
	 quoted="'"$$contract"'";\
	 echo $$quoted;

BLK_TIME = 5

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

# function name: minter
# 0x07546172
# caution: 0 must be included. DO NOT OMITT IT.
minter:
	@replkey From $$(sscli accounts list|row 1) unsigned.tx y;\
	 replkey Data 07546172 unsigned.tx y;\
	 make -sC . execute;

# function name: balances
# parameters: 'address'
# 0x27e235e3
balances: #param.address
	@queryaddr=$$(make -sC . param.address);\
	 replkey From $$(sscli accounts list|row 1) unsigned.tx y;\
	 replkey Data 27e235e3$$queryaddr unsigned.tx y;\
	 make -sC . execute;

# function name: mint
# parameters: 'address','uint256'
# 0x40c10f19
mint: #param.address param.int
	@receiver=$$(make -sC . param.address);\
	 amount=$$(make -sC . param.int);\
	 replkey From $$(sscli accounts list|row 1) unsigned.tx y;\
	 replkey Data 40c10f19$$receiver$$amount unsigned.tx y;\
	 make -sC . execute;
	 
# function name: send
# parameters: 'address','uint256'
# 0xd0679d34
send: #param.address param.int
	@read -p "fromaddress: " sender;\
	 echo toaddress@; receiver=$$(make -sC . param.address);\
	 echo amount@; amount=$$(make -sC . param.int);\
	 replkey From $$sender unsigned.tx y;\
	 replkey Data d0679d34$$receiver$$amount unsigned.tx y;\
	 make -sC . execute;