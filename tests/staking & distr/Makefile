run.withdraw:
	# @read -p "delegator_address: " delegator_address;\
	#  read -p "validator_address: " validator_address;\
	#  replkey delegator_address $$delegator_address withdraw.tx y;\
	#  replkey validator_address $$validator_address withdraw.tx y;
	 contract=$$(oneline withdraw.tx | excludestr "    ");\
	 echo "unsigned############################################################";\
	 unsigned=$$(ssutils json2hex "$$contract");\
	 echo "$$unsigned";\
	 signed=$$(sscli tx sign "$$unsigned");\
	 echo "signed############################################################";\
	 echo "$$signed";
