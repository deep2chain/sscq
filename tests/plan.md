## plot
  * [ ] **topology**: networks
    * [x] **local**: standalone
    * [x] **local**: docker-compose
    * [ ] **remote**: ethernet
    * [ ] **remote**: internet
  * [ ] **checkpoint**: curl test
  * [ ] **checkpoint**: cmd test
  * [ ] **checkpoint**: potential & historical hacks
  * [ ] **checkpoint**: possible security holes(leaks)
    * [ ] *backdoors/uap*(rest/cmd apis)
      * [ ] sscq backdoors
      * [ ] tendermint backdoors
      * [ ] cosmos backdoors
    * [ ] *system-issuer* manipulation attack
    * [ ] [DDoS attack on *persistant peers* written in config.toml](https://cosmos.network/docs/cosmos-hub/validators/validator-setup.html#what-is-a-validator)
    * [ ] *signature bypassing* attack
    * [ ] fee & reward system
      * [ ]  minimum gas prices logic
      * [ ]  gas used, gas wanted logic
      * [ ]  deduction, distribution logic
    * [ ] validators add/remove
  * [ ] **checkpoint**: private project security issues
  * [ ] **checkpoint**: validator integrity check

## self-checking points
#### 1. what if priv_validator_key.json is removed in validator nodes? Any bad effects on signing or validating? In security aspect, privs keys should be disconnected to public network.
#### 2. what if node_key.json is removed in validator nodes?
#### 3. what if addrbook.json is removed in validator nodes?
#### 4. what if .sscli is removed in validator nodes?
#### 5. what will happen when validators are upgraded?
#### 6. 

## simulations
#### 1. DDoS attacks on Sentry Nodes
#### 2. Load Test
#### 3. BruteForce Attacks on Key Management System
#### 4. Routing Attacks / Sybil Attacks / ...