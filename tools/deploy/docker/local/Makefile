# Makefile for the "ssdnode" docker image.
DOCKER_VALIDATOR_IMAGE = falcon0125/ssdnode
DOCKER_CLIENT_IMAGE = falcon0125/ssclinode

all: daemon client

daemon:
	docker build --tag ${DOCKER_VALIDATOR_IMAGE} ssdnode

client:
	docker build --tag ${DOCKER_CLIENT_IMAGE} ssclinode

CHAIN_ID = testchain
GENESIS_ACCOUNT_PASSWORD = 12345678
GENESIS_ACCOUNT_BALANCE = 6000000000000000satoshi
MINIMUM_GAS_PRICES = 100satoshi

testing:
	@docker run --rm -it -v ~/.ssd:/root/.ssd -v ~/.sscli:/root/.sscli ${DOCKER_CLIENT_IMAGE} ssd init --chain-id $(CHAIN_ID) 
	@docker run --rm -it -v ~/.ssd:/root/.ssd -v ~/.sscli:/root/.sscli ${DOCKER_CLIENT_IMAGE} ssd add-genesis-account $$(sscli accounts new ${GENESIS_ACCOUNT_PASSWORD}) ${GENESIS_ACCOUNT_BALANCE}
	@docker run --rm -it -v ~/.ssd:/root/.ssd -v ~/.sscli:/root/.sscli ${DOCKER_CLIENT_IMAGE} ssd add-genesis-account $$(sscli accounts new ${GENESIS_ACCOUNT_PASSWORD}) ${GENESIS_ACCOUNT_BALANCE}
	@docker run --rm -it -v ~/.ssd:/root/.ssd -v ~/.sscli:/root/.sscli ${DOCKER_CLIENT_IMAGE} sscli config chain-id ${CHAIN_ID}
	@docker run --rm -it -v ~/.ssd:/root/.ssd -v ~/.sscli:/root/.sscli ${DOCKER_CLIENT_IMAGE} sscli config output json
	@docker run --rm -it -v ~/.ssd:/root/.ssd -v ~/.sscli:/root/.sscli ${DOCKER_CLIENT_IMAGE} sscli config indent true
	@docker run --rm -it -v ~/.ssd:/root/.ssd -v ~/.sscli:/root/.sscli ${DOCKER_CLIENT_IMAGE} sscli config trust-node true
	@docker run --rm -it -v ~/.ssd:/root/.ssd -v ~/.sscli:/root/.sscli ${DOCKER_CLIENT_IMAGE} ssd start --minimum-gas-prices=$(MINIMUM_GAS_PRICES)

.PHONY: all, daemon, client, testing
