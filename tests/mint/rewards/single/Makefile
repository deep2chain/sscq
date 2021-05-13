# junying-todo, 2020-01-13
# sscli config json format based
# ex: 
# $ sscli query distr rewards sscq1gjtlth85y2ucym7rncrz5y6jcxp77vaz6drjuz sscqvaloper1gjtlth85y2ucym7rncrz5y6j
# [
#  {
#    "denom": "satoshi",
#    "amount": "3175575109510898.808000000000000000"
#  }
# ]
acc.rewards:
	@acc=$(shell sscli accounts list| sed -n '1p');\
	 reward=$$(sscli query distr rewards $${acc}  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$reward;

second.rewards:
	@acc=$(shell sscli accounts list| sed -n '2p');\
	 reward=$$(sscli query distr rewards $${acc}  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$reward;

community.pool:
	@reward=$$(sscli query distr community-pool  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$reward;

commision:
	@validator_address=$$(findstr validator_address ~/.ssd/config/genesis.json  | fromstr ': "'|endstr '"');\
	 reward=$$(sscli query distr commission $$validator_address  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$reward;

rewards:
	@read -p "accaddr: " accaddr; \
	 rewards=$$(sscli query distr rewards $$accaddr  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$rewards

height:
	@height=$$(sscli query block | sed -n "16p"| fromstr ': "'|endstr '"');\
	 echo $$height

total_txs:
	@total_txs=$$(sscli query block | sed -n "19p"| fromstr ': "'|endstr '"');\
	 echo $$total_txs

num_txs:
	@num_txs=$$(sscli query block | sed -n "18p"| fromstr ': "'|endstr '"');\
	 echo $$num_txs

print.rewards:
	@r_old=0;cp_old=0;cm_old=0;num_txs_old=0;total_txs_old=0;\
	 for i in $$(python -c "print ' '.join(str(item) for item in range(100))"); \
	 do \
	 r=$$($(MAKE) -s acc.rewards);\
	 cp=$$($(MAKE) -s community.pool);\
	 cm=$$($(MAKE) -s commision);\
	 d1=$$(python -c "print $$r-$$r_old");\
	 d2=$$(python -c "print $$cp-$$cp_old");\
	 d3=$$(python -c "print $$cm-$$cm_old");\
	 r_old=$$r;cp_old=$$cp;cm_old=$$cm;\
	 echo $$r $$cp $$cm;\
	 mr=$$(python -c "print $$d1+$$d2+$$d3");\
	 height=$$($(MAKE) -s height);\
	 total_txs=$$($(MAKE) -s total_txs);\
	 num_txs=$$($(MAKE) -s num_txs);\
	 echo delta:$$height:$$num_txs:$$total_txs--- $$d1 $$d2 $$d3 $$mr $$(python -c "print $$mr-$$num_txs_old*3000000");\
	 num_txs_old=$$num_txs;total_txs_old=$$total_txs;\
	 echo ratio:$$height:$$num_txs:$$total_txs--- $$(python -c "print '%.2f'%($$d1/$$mr)") $$(python -c "print '%.2f'%($$d2/$$mr)") $$(python -c "print '%.2f'%($$d3/$$mr)");\
	 sleep 4;\
	 done\
