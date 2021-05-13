

## rest  send


```
curl --location --request POST "http://127.0.0.1:1317/ss/send" \
  --header "Content-Type: application/x-www-form-urlencoded" \
  --data "    { \"base_req\": 
                          { \"from\": \"sscq103x7taejyqwxrvyadu2yxd7u04wdqs5stq5a40\", 
                            \"memo\": \"\",
                            \"password\": \"12345678\", 
                            \"chain_id\": \"testchain\", 
                            \"account_number\": \"3\", 
                            \"sequence\": \"1\", 
                            \"gas_wanted\": \"200000\", 
                            \"gas_price\": \"20\", 
                            \"simulate\": false
                          },          
                \"amount\": [ 
                        { \"denom\": \"sscq\", 
                          \"amount\": \"0.00001\" } ],
                \"to\": \"sscq1ucwpvw99u9tj3sxuwtnesge5tl90c6y0zcnc73\",
                \"data\": \"\"
    }"


```