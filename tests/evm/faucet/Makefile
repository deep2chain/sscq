all: sscq_faucet call_faucet sscq_faucet_with_bug hack_faucet

sscq_faucet:
	solcjs  --bin --abi  sscq_faucet.sol

call_faucet:
	solcjs  --bin --abi  call_faucet.sol sscq_faucet.sol

sscq_faucet_with_bug:
	cd hack && solcjs  --bin --abi  sscq_faucet_with_bug.sol && cd ..

hack_faucet:
	cd hack && solcjs  --bin --abi  hack_faucet.sol  sscq_faucet_with_bug.sol && cd ..

clean:
	- rm ./*.bin 
	- rm ./hack/*.bin  
	- rm ./*.abi
	- rm ./hack/*.abi

.PHONY: clean sscq_faucet call_faucet sscq_faucet_with_bug hack_faucet
