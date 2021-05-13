# junying-todo, 2020-01-13
# sscli config text format based
# ex: 
# $ sscli query distr rewards sscq1j9fhvpfrhlmx2m3xn487qpuyhe56x7wfptssqz
# $ 271771554720.370000000000000000satoshi
# [metrics]
# commission rate, community pool rate
# [checkpoints]
# mining rewards = fixed mining rewards +  tx fees
# mining rewards = validator rewards(commission) + delegator rewards + community pool
# validator_reward = validator rewards + delegator rewards
# 
# 
# 
DELEGATOR_CNTS = $$(findstr delegator_address ~/.ssd/config/genesis.json | linecount)
VALIDATOR_CNTS = $$(findstr validator_address ~/.ssd/config/genesis.json | linecount)

clean:
	@rm -rf sscq*

initialize:
	@for i in $$(python -c "print ' '.join(str(item) for item in range(${VALIDATOR_CNTS}))"); \
	 do \
	 validator=$$(findstr validator_address ~/.ssd/config/genesis.json|row $$(($$i+1))|fromstr ': "'|endstr '"');\
	 if [ -d "$$validator" ]; then rm -rf $$validator; mkdir $$validator; else mkdir $$validator;fi;\
	 delegator=$$(findstr delegator_address ~/.ssd/config/genesis.json|row $$(($$i+1))|fromstr ': "'|endstr '"');\
	 echo $$delegator >> $$validator/delegators.list;\
	 done

# VALIDATOR_CNTS = $$(ls|grep sscq|linecount)

print.validator.rewards:
	@for i in $$(python -c "print ' '.join(str(item) for item in range(${VALIDATOR_CNTS}))"); \
	 do \
	 validator=$$(ls|grep sscq|row $$(($$i+1)));\
	 echo V$$i:$$validator;\
	 reward=$$(sscli query distr commission $$validator);\
	 echo VRewards:$$reward;\
	 done;

print.delegator.rewards:
	@for i in $$(python -c "print ' '.join(str(item) for item in range(${VALIDATOR_CNTS}))"); \
	 do \
	 rewards=$$(echo DRewards:);\
	 validator=$$(ls|grep sscq|row $$(($$i+1)));\
	 echo ------------------------------------------------------;\
	 echo V$$i:$$validator;\
	 delcnts=$$(linecount $$validator/delegators.list);\
	 for j in $$(python -c "print ' '.join(str(item) for item in range($$delcnts))"); \
	 do \
	 delegator=$$(cat $$validator/delegators.list |row $$(($$j+1)));\
	 echo D$$j:$$delegator;\
	 reward=$$(sscli query distr rewards $$delegator $$validator);\
	 rewards=$$(echo $$rewards $$reward);\
	 done;\
	 echo $$rewards;\
	 done;
	 
print.all:
	@for i in $$(python -c "print ' '.join(str(item) for item in range(${VALIDATOR_CNTS}))"); \
	 do \
	 echo ------------------------------------------------------;\
	 validator=$$(ls|grep sscq|row $$(($$i+1)));\
	 vreward=$$(sscli query distr commission $$validator|excludestr satoshi|endstr .);\
	 drewards=$$(echo );\
	 echo V$$(($$i+1)): $$validator, $$vreward;\
	 delcnts=$$(linecount $$validator/delegators.list);\
	 dsum=0;\
	 for j in $$(python -c "print ' '.join(str(item) for item in range($$delcnts))"); \
	 do \
	 delegator=$$(cat $$validator/delegators.list|row $$(($$j+1)));\
	 echo D$$(($$j+1)): $$delegator;\
	 dreward=$$(sscli query distr rewards $$delegator $$validator|excludestr satoshi|endstr .);\
	 drewards=$$(echo $$drewards $$dreward);\
	 dsum=$$(python -c "print $$dsum+int($$dreward);");\
	 done;\
	 echo $$vreward $$dsum $$drewards;\
	 dsum=0;\
	 done;

community.pool:
	@reward=$$(sscli query distr community-pool  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$reward;

check.rewards:
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
	 sleep 5;\
	 done\
