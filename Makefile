# detect operating system
ifeq ($(OS),Windows_NT)
    CURRENT_OS := Windows
else
    CURRENT_OS := $(shell uname -s)
endif

export GO11MODULE=on

# default: log
export LOG_LEVEL=debug

#GOBIN
GOBIN = $(shell pwd)/build/bin
GO ?= latest

# variables
DEBUGAPI=ON  # disable DEBUGAPI by default

PACKAGES = $(shell go list ./... | grep -Ev 'vendor|importer')
COMMIT_HASH := $(shell git rev-parse --short HEAD)
GIT_BRANCH :=$(shell git branch 2>/dev/null | grep "^\*" | sed -e "s/^\*\ //")
# tool checking
DEP_CHK := $(shell command -v dep 2> /dev/null)
GOLINT_CHK := $(shell command -v golint 2> /dev/null)
GOMETALINTER_CHK := $(shell command -v gometalinter.v2 2> /dev/null)
UNCONVERT_CHK := $(shell command -v unconvert 2> /dev/null)
INEFFASSIGN_CHK := $(shell command -v ineffassign 2> /dev/null)
MISSPELL_CHK := $(shell command -v misspell 2> /dev/null)
ERRCHECK_CHK := $(shell command -v errcheck 2> /dev/null)
UNPARAM_CHK := $(shell command -v unparam 2> /dev/null)
#
LEDGER_ENABLED ?= true

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

BUILD_FLAGS = -tags "$(build_tags)" -ldflags '-X version.GitCommit=${COMMIT_HASH} -X main.GitCommit=${COMMIT_HASH} -X main.DEBUGAPI=${DEBUGAPI} -X main.GitBranch=${GIT_BRANCH}'
BUILD_FLAGS_STATIC_LINK = -tags "$(build_tags)" -ldflags '-X github.com/deep2chain/sscq/version.GitCommit=${COMMIT_HASH} -X main.GitCommit=${COMMIT_HASH} -X main.DEBUGAPI=${DEBUGAPI} -X main.GitBranch=${GIT_BRANCH} -linkmode external -w -extldflags "-static"'

all: build

# all: tools deps build

# tools:
# ifndef DEP_CHK
# 	@echo "Installing dep"
# 	go get -u -v github.com/golang/dep/cmd/dep
# else
# 	@echo "Dep is already installed..."
# endif

# deps:
# 	@echo "--> Generating vendor directory via dep ensure"
# 	@rm -rf .vendor-new
# 	@dep ensure -v -vendor-only

# update:
# 	@echo "--> Running dep ensure"
# 	@rm -rf .vendor-new
# 	@dep ensure -v -update

buildquick: go.sum
ifeq ($(CURRENT_OS),Windows)
	@echo BUILD_FLAGS=$(BUILD_FLAGS)
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bin/hsd.exe ./cmd/hsd
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bin/hscli.exe ./cmd/hscli
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bin/hsutils.exe ./cmd/hsutil
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bin/hscli.exe ./cmd/hsinfo
else
	@echo BUILD_FLAGS=$(BUILD_FLAGS)
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bin/hsd ./cmd/hsd
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bin/hscli ./cmd/hscli
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bin/hsutils ./cmd/hsutil
	@go build -mod=readonly $(BUILD_FLAGS) -o build/bin/hsinfo ./cmd/hsinfo
endif

# https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
# https://stackoverflow.com/questions/36279253/go-compiled-binary-wont-run-in-an-alpine-docker-container-on-ubuntu-host,
# failed because dependency path modified
build.CGO_DISABLED: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(MAKE) buildquick

build.static: go.sum
	@echo BUILD_FLAGS=$(BUILD_FLAGS_STATIC_LINK)
	@go build -mod=readonly $(BUILD_FLAGS_STATIC_LINK) -o build/testnet/hsd ./cmd/hsd
	@go build -mod=readonly $(BUILD_FLAGS_STATIC_LINK) -o build/testnet/hscli ./cmd/hscli

build: unittest buildquick

build-batchsend:
	@build/env.sh go run build/ci.go install ./cmd/hsbatchsend

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/hsd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/hscli

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify


# test part
test:
	@go test --vet=off $(PACKAGES)
	@echo $(PACKAGES)

unittest:
	@go test -v ./evm/...
	@go test -v ./types/...
	@go test -v ./store/...
	@go test -v ./utils/...
	@go test -v ./x/mint/...
	@go test -v ./x/bank/...
	@go test -v ./x/core/...

	@go test -v ./accounts/...
	@go test -v ./app/...
	@go test -v ./client/...
	@go test -v ./init/...
	@go test -v ./crypto/...
	@go test -v ./server/...
	@go test -v ./tools/...
	@go test -v ./x/auth/...
	@go test -v ./x/crisis/...
	@go test -v ./x/distribution/...
	@go test -v ./x/gov/...
	@go test -v ./x/guardian/...
	@go test -v ./x/ibc/...
	@go test -v ./x/params/...
	@go test -v ./x/slashing/...
	@go test -v ./x/staking/...

CHAIN_ID = testchain
GENESIS_ACCOUNT_PASSWORD = 12345678
GENESIS_ACCOUNT_BALANCE = 3000000000000000satoshi
MINIMUM_GAS_PRICES = 100satoshi

new: install clear hsinit accs conf vals

new.pure: clear hsinit accs conf vals

hsinit:
	@hsd init mynode --chain-id $(CHAIN_ID)

accs:
	@echo create new accounts....;\
    $(eval ACC1=$(shell hscli accounts new $(GENESIS_ACCOUNT_PASSWORD)))\
	$(eval ACC2=$(shell hscli accounts new $(GENESIS_ACCOUNT_PASSWORD)))
	@hsd add-genesis-account $(ACC1) $(GENESIS_ACCOUNT_BALANCE)
	@hsd add-guardian-account $(ACC1) 
	@hsd add-genesis-account $(ACC2) $(GENESIS_ACCOUNT_BALANCE)

conf:
	@echo setting config....
	@hscli config chain-id $(CHAIN_ID)
	@hscli config output json
	@hscli config indent true
	@hscli config trust-node true

vals:
	@echo setting validators....
	@hsd gentx $(ACC1)
	@hsd collect-gentxs

start: start.daemon start.rest

start.daemon:
	@echo starting daemon....
	@nohup hsd start >> ${HOME}/.hsd/app.log  2>&1  &

start.rest:
	@echo starting rest server...
	@nohup hscli rest-server --chain-id=${CHAIN_ID} --trust-node=true --laddr=tcp://0.0.0.0:1317 >> ${HOME}/.hsd/restServer.log  2>&1  &

stop:
	@pkill hsd
	@pkill hscli

# clean part
clean:
	@find build -name bin | xargs rm -rf

clear: clean
	@rm -rf ~/.hs*

DOCKER_VALIDATOR_IMAGE = falcon0125/hsdnode
DOCKER_CLIENT_IMAGE = falcon0125/hsclinode
VALIDATOR_COUNT = 4
TESTNODE_NAME = client
TESTNETDIR = build/testnet
LIVENETDIR = build/livenet

##############################################################################################################################
# Run a 4-validator testnet locally
##############################################################################################################################

# docker-compose part[multi-node part, also test mode]
# Local validator nodes using docker and docker-compose
hsnode: clean build.static# hstop
	$(MAKE) -C tools/deploy/docker/local

echotest:
	@echo  $(CURDIR)/${TESTNETDIR}

hsinit-v4: 
	@if ! [ -f ${TESTNETDIR}/node0/.hsd/config/genesis.json ]; then\
	 docker run --rm -v $(CURDIR)/build/testnet:/root:Z ${DOCKER_VALIDATOR_IMAGE} testnet \
																				  --chain-id ${CHAIN_ID} \
																				  --v ${VALIDATOR_COUNT} \
																				  -o . \
																				  --starting-ip-address 192.168.10.2 \
																				  --minimum-gas-prices ${MINIMUM_GAS_PRICES}; fi
hsinit-test: 
	@hsd testnet --chain-id ${CHAIN_ID} \
				 --v ${VALIDATOR_COUNT} \
				 -o ${TESTNETDIR} \
				 --starting-ip-address 192.168.10.2 \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES}
hsinit-o1:
	@mkdir -p ${TESTNETDIR}/node4/.hsd ${TESTNETDIR}/node4/.hscli
	@hsd init node4 --home ${TESTNETDIR}/node4/.hsd
	@cp ${TESTNETDIR}/node0/.hsd/config/genesis.json ${TESTNETDIR}/node4/.hsd/config
	# @cp ${TESTNETDIR}/node0/.hsd/config/hsd.toml ${TESTNETDIR}/node4/.hsd/config
	@cp ${TESTNETDIR}/node0/.hsd/config/config.toml ${TESTNETDIR}/node4/.hsd/config
	@sed -i s/node0/node4/g ${TESTNETDIR}/node4/.hsd/config/config.toml
	@cp -rf ${TESTNETDIR}/node0/.hscli/* ${TESTNETDIR}/node4/.hscli

hsinit-o2:
	@mkdir -p ${TESTNETDIR}/node5/.hsd ${TESTNETDIR}/node5/.hscli
	@hsd init node5 --home ${TESTNETDIR}/node5/.hsd
	@cp ${TESTNETDIR}/node0/.hsd/config/genesis.json ${TESTNETDIR}/node5/.hsd/config
	# @cp ${TESTNETDIR}/node0/.hsd/config/hsd.toml ${TESTNETDIR}/node5/.hsd/config
	@cp ${TESTNETDIR}/node0/.hsd/config/config.toml ${TESTNETDIR}/node5/.hsd/config
	@sed -i s/node0/node5/g ${TESTNETDIR}/node5/.hsd/config/config.toml
	@cp -rf ${TESTNETDIR}/node1/.hscli/* ${TESTNETDIR}/node5/.hscli

hstart: build.static hsinit-test hsinit-o1 hsinit-o2
	@docker-compose up -d

hstart.debug: build hsinit-test hsinit-o1 hsinit-o2
	@docker-compose up

hsattach:
	@docker attach hsclinode1

# Stop testnet
hstop:
	docker-compose down

hscheck:
	@docker logs -f hsdnode0

hsclean:
	@docker rmi ${DOCKER_VALIDATOR_IMAGE} ${DOCKER_CLIENT_IMAGE}

##############################################################################################################################
# ethernet part
##############################################################################################################################
clean-t:
	@find build -name testnet |xargs rm -rf
	
# addrs:
# 	@if [ -f ipaddrs.conf ]; then rm ipaddrs.conf ;fi
# 	# modify conf files
# 	@for index in $$(seq -s ' ' 4); do \
# 	 read -p "Enter node$$index IP addr: " ipaddr;\
# 	 echo $$ipaddr >> ipaddrs.conf; done

# killall:
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '1p') pkill -9 hsd
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '2p') pkill -9 hsd
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '3p') pkill -9 hsd
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '4p') pkill -9 hsd

# startall:
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '1p') nohup hsd start & > /dev/null
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '2p') nohup hsd start & > /dev/null
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '3p') nohup hsd start & > /dev/null
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '4p') nohup hsd start & > /dev/null

# cleanall:
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '1p') rm -rf /root/.hsd /root/.hscli
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '2p') rm -rf /root/.hsd /root/.hscli
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '3p') rm -rf /root/.hsd /root/.hscli
# 	@sshpass -p miss16980 ssh root@$$(cat networks/remote/ipaddrs.conf | sed -n '4p') rm -rf /root/.hsd /root/.hscli

# copyall:
# 	# upload files
# 	### 1st server
# 	@sshpass -p miss16980 scp -r ${TESTNETDIR}/node0/.hsd root@$$(cat networks/remote/ipaddrs.conf | sed -n '1p'):/root
# 	@sshpass -p miss16980 scp -r ${TESTNETDIR}/node0/.hscli root@$$(cat networks/remote/ipaddrs.conf | sed -n '1p'):/root
# 	@sshpass -p miss16980 scp -r build/bin/hsd root@$$(cat networks/remote/ipaddrs.conf | sed -n '1p'):/usr/local/bin
# 	### 2nd server
# 	@sshpass -p miss16980 scp -r ${TESTNETDIR}/node1/.hsd root@$$(cat networks/remote/ipaddrs.conf | sed -n '2p'):/root
# 	@sshpass -p miss16980 scp -r ${TESTNETDIR}/node1/.hscli root@$$(cat networks/remote/ipaddrs.conf | sed -n '2p'):/root
# 	@sshpass -p miss16980 scp -r build/bin/hsd root@$$(cat networks/remote/ipaddrs.conf | sed -n '2p'):/usr/local/bin
# 	### 3rd server
# 	@sshpass -p miss16980 scp -r ${TESTNETDIR}/node2/.hsd root@$$(cat networks/remote/ipaddrs.conf | sed -n '3p'):/root
# 	@sshpass -p miss16980 scp -r build/testnet/node2/.hscli root@$$(cat networks/remote/ipaddrs.conf | sed -n '3p'):/root
# 	@sshpass -p miss16980 scp -r build/bin/hsd root@$$(cat networks/remote/ipaddrs.conf | sed -n '3p'):/usr/local/bin
# 	### 4th server
# 	@sshpass -p miss16980 scp -r ${TESTNETDIR}/node3/.hsd root@$$(cat networks/remote/ipaddrs.conf | sed -n '4p'):/root
# 	@sshpass -p miss16980 scp -r ${TESTNETDIR}/node3/.hscli root@$$(cat networks/remote/ipaddrs.conf | sed -n '4p'):/root
# 	@sshpass -p miss16980 scp -r build/bin/hsd root@$$(cat networks/remote/ipaddrs.conf | sed -n '4p'):/usr/local/bin

# resetall: #clean-4 install-
# 	@if ! [ -d ${TESTNETDIR} ]; then mkdir -p ${TESTNETDIR}; fi
# 	@hsd testnet --chain-id mainchain \
# 				 --v 4 \
# 				 -o ${TESTNETDIR} \
# 				 --validator-ip-addresses $(CURDIR)/networks/remote/ipaddrs.conf \
# 				 --minimum-gas-prices ${MINIMUM_GAS_PRICES}

# clean-testnet:
# 	@rm -rf $(CURDIR)/build/testnet

# testnet: clean-testnet install resetall #copyall startall # killall cleanall 

# chketh:
# 	@sshpass -p miss16980 ssh root@192.168.10.69

##############################################################################################################################
# ethernet distribution part
##############################################################################################################################
clean-livenet:
	@rm -rf $(CURDIR)/build/livenet

distall: #clean-4 install-
	@if ! [ -d ${LIVENETDIR} ]; then mkdir -p ${LIVENETDIR}; fi
	@hsd livenet --chain-id livenet \
				 --v $$(wc $(CURDIR)/networks/remote/ipaddrs.conf | awk '{print$$1F}') \
				 -o ${LIVENETDIR} \
				 --validator-ip-addresses $(CURDIR)/networks/remote/ipaddrs.conf \
				 --minimum-gas-prices ${MINIMUM_GAS_PRICES}

livenet: clean-livenet install distall

##############################################################################################################################
# load test part
##############################################################################################################################
loadtest:
	@locust -f $(CURDIR)/tests/locustfile.py --host=http://127.0.0.1:1317

.PHONY: build install build- install- \
		test clean clean-t \
		testnet livenet \
		stop
