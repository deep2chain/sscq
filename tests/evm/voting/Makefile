# junying-todo, 2019-08-06
# https://solidity.readthedocs.io/en/v0.4.21/solidity-by-example.html
# fails in creating

# contract tx hash: 
# contract addr: 
# issuer:		 sscq1ha7ryup8nc2avgesfunx2pm22waqv2cx6dj0ac
# vote
## parameters: 'uint256'
## query:	0x0121b93f
#---------------------------------
# proposals
## parameters: 'bytes32','uint256'
## query:	0x930dc17c
#---------------------------------
# chairperson
## parameters: 
## query:	0x2e4176cf
#---------------------------------
# chairperson
## parameters: 
## query:	0x2e4176cf

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