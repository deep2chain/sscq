## gist
    * private ip for validator node
    * public ip for sentry node
    * ip filtering by firewall
## [cases](https://medium.com/irisnet-blog/tech-choices-for-cosmos-validators-27c7242061ea)
#### single node validator setup - firewall
#### single node sentry node setup
#### two layer sentry node setup
#### relay network setup
## [setup](https://cosmos.network/docs/cosmos-hub/validators/security.html#validator-security)
#### validator node
    [config.toml]
    persistent_peers =[list of sentry nodes]
    pex = false
#### sentry node
    [config.toml]
    private_peer_ids = "node_ids_of_private_peers"
#### delegator node
    ?

## reference
#### [Anti-DDoS](https://medium.com/irisnet-blog/tech-choices-for-cosmos-validators-27c7242061ea)