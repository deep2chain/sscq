# junying-todo, 2020-06-09

Contract_Address = sscq1aluln32wfnsd0r868t6j262la58xcw8h289462
TxHash = 6494C1263587445D8F2F3509DE6D4A54E38D7F154532100CD07FEE196EDF2210

calc.method.id:
	@make -sC ../ get.method.id

param.address:
	@make -sC ../ param.address
	
param.int:
	@make -sC ../ param.int

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

#############################################################################
#############################################################################
# Airdrop　                                                                 #
# transferFrom A----->B														#
#                 |  														#
#           (tx fee on C)										            #
# ,where C has privilege to handle the approved amount of tokens from A as  #
# well as paying fee.														#
# C is airdrop's contract address, A is airdrop's msg.sender                #
# r:https://stevenocean.github.io/2018/07/08/token-airdrop-with-web3py.html #
#############################################################################
# [HRC20 Contract]
# hrc20.het.contract.addr:sscq15zuywqfgvw955wg9hvdgkr3hxsdgja5z2tvrmc(000000000000000000000000a0b8470128638b4a3905bb1a8b0e37341a897682)
# approve A(airdrop.msg.sender) approved B(airdrop.contract.addr) to spend Z
# A: 00000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe9(sscq1sh8d3h0nn8t4e83crcql80wua7u3xtlft9sr5d)
# B: 000000000000000000000000eff9f9c54e4ce0d78cfa3af525695fed0e6c38f7(sscq1aluln32wfnsd0r868t6j262la58xcw8h289462)
# Z: 00000000000000000000000000000000000000000000021e19e0c9bab2400000(10000000000000000000000)
# approve.data:095ea7b3000000000000000000000000eff9f9c54e4ce0d78cfa3af525695fed0e6c38f700000000000000000000000000000000000000000000021e19e0c9bab2400000
# allowance: sscli query contract sscq15zuywqfgvw955wg9hvdgkr3hxsdgja5z2tvrmc dd62ed3e00000000000000000000000085ced8ddf399d75c9e381e01f3bddcefb9132fe90000000000000000000000002ff97e52f5418eaf5d7d88cb919cc74a735a7128
#**********************************************************************************************************************************************************
# [Airdrop Contract](airdrop:sscq1vcv6jz4nmunv8fndyhznuwuc96kuccsj4cuz34)
# contract FutureEdgeAirdrop {
#     function drop(address tokenAddr, address[] dests, uint256 value) public payable {
#         uint256 valuePerCount = value / dests.length；
#         for (uint i = 0; i < dests.length; i++) {
#             ERC20Token(tokenAddr).transferFrom(msg.sender, dests[i], valuePerCount);
#         }
#     }
# }
# airdrop.contract.addr: sscq1aluln32wfnsd0r868t6j262la58xcw8h289462(000000000000000000000000eff9f9c54e4ce0d78cfa3af525695fed0e6c38f7)
# msg.sender:sscq1sh8d3h0nn8t4e83crcql80wua7u3xtlft9sr5d
# tokenaddr: sscq15zuywqfgvw955wg9hvdgkr3hxsdgja5z2tvrmc(000000000000000000000000a0b8470128638b4a3905bb1a8b0e37341a897682)
# receivers: sscq1dqrh6tf6js9flf3clmzq4jt3rg7nz02ka6f305(00000000000000000000000068077d2d3a940a9fa638fec40ac9711a3d313d56)
#			 sscq1jcz5glz5zn4z9j903qmwkyydtvpx58n6vynk25(0000000000000000000000009605447c5414ea22c8af8836eb108d5b026a1e7a)
#			 sscq14p4k25d0g8qfavre5nz8dlcud927j5z2cu5lqn(000000000000000000000000a86b6551af41c09eb079a4c476ff1c6955e9504a)
#			 sscq1c59cde6ga8qpf48leenq6u68gdmjmz3ue6jgxz(000000000000000000000000c50b86e748e9c014d4ffce660d734743772d8a3c)
#			 sscq1cu68kn2uydclu6aq549gqe33mfsgay55vh4t5u(000000000000000000000000c7347b4d5c2371fe6ba0a54a806631da608e9294)
# value:     10000000000000000000(0000000000000000000000000000000000000000000000008ac7230489e80000)
# method.id + tokenaddr + pos.addr(hex(32*3)) + pos.value(hex(32*(3+1+array.len))) + addr.array.len + receiver.addr.1 + receiver.addr.2 + ... + value.array.len + value.1 + value.2 +...
# 7da5efc8000000000000000000000000f1c106c2361723e17c2d5253e473fbb544a4702e000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c50b86e748e9c014d4ffce660d734743772d8a3c000000000000000000000000c7347b4d5c2371fe6ba0a54a806631da608e9294000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000001110000000000000000000000000000000000000000000000000000000000000112
# sscli query contract sscq178qsds3kzu37zlpd2ff7gulmk4z2gupwd8qm3p 70a08231
#*******************************************************************************************************

# function name: drop
# parameters: 'address','address[]','uint256[]'
# 0x7da5efc8
# drop: #param.address param.int
# 	@read -p "fromaddress: " sender;\
# 	 echo tokenaddr@; tokenaddr=$$(make -sC . param.address);\
# 	 echo toaddr@; toaddr=$$(make -sC . param.address);\
# 	 echo amount@; amount=$$(make -sC . param.int);\
# 	 replkey From $$sender unsigned.tx y;\
# 	 replkey Data 7da5efc8$$tokenaddr$$toaddr$$amount unsigned.tx y;\
# 	 make -sC . execute;