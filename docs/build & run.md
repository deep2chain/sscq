### Build    
    # source code
    mkdir -p sourceBuffer/github.com/orientwalt
    cd sourceBuffer/github.com/orientwalt
    git clone https://github.com/deep2chain/sscq.git    
    cd sscq
    
    # warning
    # we use go mod to manage dependency package
    # so sourceBuffer must not be in $GOPATH
    # more about go mod , see "./use_go_mod_to_manage_dependency_package.md"
    
    # GO111MODULE on;  no more depend on deps or vendor 
    export GO111MODULE=on
    
    # use proxy
    export GOPROXY=https://goproxy.io
    
    # use proxy (another proxy)
    # when depend on your private github.com repository , use 'https://goproxy.cn,direct' , and go get your  repository 
    #   $export GOPROXY=https://goproxy.cn,direct
    #   $go get github.com/orientwalt/tendermint    
    
    # set ApiSecuritylevel
    #   see below #ApiSecuritylevel for more detail 
    export DEBUGAPI=ON  ##  "ON", develop and test version (ApiSecuritylevel=low); "OFF", default value, production version (ApiSecuritylevel=high); 
    
    # compile and install
    make install    
        
    # turn off proxy temporary when orientwalt/tendermint can not download
    export GOPROXY=
    make install
    
    # ... after download orientwalt/tendermint success, turn on proxy again
    export GOPROXY=https://goproxy.io
    make install
    
    # if turn off proxy not work for orientwalt/tendermint
    #  copy directory $GOPATH/pkg/mod/github.com/orientwalt/tendermint@vXX.X.X+incompatible from your workmate
    #  and setup GOPRIVATE
    export GOPRIVATE=github.com/orientwalt/tendermint 

    # print the version and ApiSecuritylevel
    # make sure the version, git commit hash, ApiSecuritylevel is what you need
    ssd   version
    sscli version     
    
    
### Api Security Level
for security, api need to classification, called  Api Security Level;
- "high" : disable operate type API, like new account, send tx ,and so on; only query type API is enable  
- "low": enable all API  
-  high(level) is default


compile command recommand
- production environment:    

```
make install
```
  
- develop and test environment:  

```
export DEBUGAPI=ON   ## "ON", develop and test version (ApiSecuritylevel=low); "OFF", default value, production version (ApiSecuritylevel=high);
make install
```
  
- print the api-security-level
```
sscli version
```


### Config
    # Initialize configuration files and genesis file
    ssd init [moniker] --chain-id testchain

    # Copy the `Address` output here and save it for later use
    sscli accounts new  [password] (password can not null)
    or
    sscli accounts new OFF (input passthrase must at least 8 characters)  

    # Show all local accounts keyfile
    sscli accounts list

    # Add both accounts, with coins to the genesis file
    ssd add-genesis-account [addr] [amount] (amount:xxxstake,xxxsatoshi  , e.g: 10000000000000stake,1000000000000000satoshi)
    ssd add-guardian-account  [addr] 
    ssd gentx [genesis-account]
    ssd collect-gentxs

    # Configure your CLI to eliminate need for chain-id flag
    sscli config chain-id testchain
    sscli config output json
    sscli config indent true
    sscli config trust-node true
  
### RUN & TEST
#### Run Daemon
    ssd start
#### Run REST Server
    sscli rest-server
    sscli rest-server --chain-id=testchain --trust-node=true
    sscli rest-server --chain-id=testchain --trust-node=true --laddr=tcp://0.0.0.0:1317
                      
#### CLI TEST
    [newaccount]
    sscli accounts new 123... 
    sscli accounts new OFF       
    
    [list]
    sscli accounts list
    
    [getbalance]
    sscli query account [addr]
    
    # transaction
    sscli tx send [fromaddr] [toaddr] [amount]
    sscli tx create [fromaddr] [toaddr] [amount]
    sscli tx sign [rawdata]
    sscli tx broadcast [rawdata]

#### REST TEST
    Tip: http, not https
    
    [newaccount]
    curl -X POST "http://localhost:1317/accounts/newaccount" -H "accept: application/json" -d "{\"password\": \"12345678\"}"

    [get account information]
    curl -X GET "http://localhost:1317/auth/accounts/cosmos1ytczrhg8anm6a4z2rjhhs4rz0cvrxc5yna0f68" -H "accept: application/json"

    [getbalance]
    curl -X GET "http://localhost:1317/bank/balances/cosmos1ytczrhg8anm6a4z2rjhhs4rz0cvrxc5yna0f68" -H "accept: application/json"

    [sendTx]
    curl -X POST "http://localhost:1317/ss/send" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"base_req\": { \"from\": \"cosmos1jj4aqger28lwgpd4mfr35x59g249jnflhqdyxq\", \"memo\": \"Sent via Cosmos Voyager \",\"password\": \"xxxx\", \"chain_id\": \"testchain\", \"account_number\": \"0\", \"sequence\": \"0\", \"gas\": \"200000\", \"gas_adjustment\": \"1.2\", \"fees\": [ { \"denom\": \"sscqtoken\", \"amount\": \"10\" } ], \"simulate\": false }, \"amount\": [ { \"denom\": \"sscqtoken\", \"amount\": \"10\" } ],\"to\": \"cosmos1gncjp5n8jurnuz5hnj0t2eyvqdms7gzzg8ycjx\"}"
    
    [createTx]
    curl -X POST "http://localhost:1317/ss/create" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"base_req\": { \"from\": \"cosmos1extcaaktdfcle4areslzvxx82q5rncvyrjf8m4\", \"memo\": \"Sent via Cosmos Voyager \",\"password\": \"12345678\", \"chain_id\": \"testchain\", \"account_number\": \"0\", \"sequence\": \"0\", \"gas\": \"200000\", \"gas_adjustment\": \"1.2\", \"fees\": [ { \"denom\": \"sscqtoken\", \"amount\": \"10\" } ], \"simulate\": false }, \"amount\": [ { \"denom\": \"sscqtoken\", \"amount\": \"10\" } ],\"to\": \"cosmos1ehdzkfgvqana4gc6keuymweuhm60x73uayk0kt\",\"encode\":true}"
    
    [signTx]
    curl -X POST "http://localhost:1317/ss/sign" -H "accept: application/json" -H "Content-Type: application/json" -d "{\"tx\":\"xxx\", \"passphrase\":\"xxx\",\"offline\":false,\"encode\":true}"
    
    [broadcastTx]
    curl -X POST "http://localhost:1317/ss/broadcast" -H "accept: aplication/json" -H "Content-Type: application/json" -d "{\"tx\":\"xxxx\"}"
