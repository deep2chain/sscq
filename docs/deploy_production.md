


## 系统参数
Moniker （网络的绰号）： sscqnet
ChainName（链的名称） ： mainchain


## 账户
 两个重要的账户
- HTDF_ISSUE_ACCT  HTDF基金会/央行账户（HTDF token 发行账户）
- STAKE_BIG_ACCT   sscq抵押账户（stake 大账户，只保存stake）


 四个验证节点的账户
- VALIDATORACCT1
- VALIDATORACCT2
- VALIDATORACCT3
- VALIDATORACCT4
 

## 发行计划
账户               | 发行内容  
------------------|:-------------------
 HTDF_ISSUE_ACCT  | 伍亿 HTDF  token
 STAKE_BIG_ACCT   | 10^14 stake
 VALIDATORACCT1   | 10^13 stake
 VALIDATORACCT2   | 10^13 stake
 VALIDATORACCT3   | 10^13 stake
 VALIDATORACCT4   | 10^13 stake
 


```
    
    $HTDF_ISSUE_ACCT 50000000000000000satoshi
    STAKE_BIG_ACCT   100000000000000stake
    $VALIDATORACCT1  10000000000000stake
    $VALIDATORACCT2  10000000000000stake
    $VALIDATORACCT3  10000000000000stake
    $VALIDATORACCT4  10000000000000stake
    
    sscli accounts list
    
```

## 开始部署

- Makefile文件所在
```
    $HTDF_SRC_HOME/networks/remote/2-layers/Makefile
```


- 部署相关的配置文件
```
    $(HOME)/config 目录拷贝 *.pem + servers.info
```


- 部署
```
    cd ${SRC_HOME}
    make install

    cd ${SRC_HOME}/networks/remote/2-layers
    
    #根据服务器信息，编辑和生成  ~/config/servers.info    
    
    #根据 servers.info 修改  ${SRC_HOME}/networks/remote/2-layers/Makefile 中的参数
    #   如 DEFAULT_VAL_COUNT 等服务器配置参数
    
    
    make reset
    make dist-makefile
    make clean
    make stop
    make reset
    make regen
    chmod 400 ~/config/*.pem
    make dist
    make dist-ssd
    make dist-sscli
    make start-daemon
    make start-rest

```


## 检查运行情况
- 检查链的chain-id
> "mainchain"
- 是否出块
- 网络连接检查
> P2P、RCP端口和连接是否正常

- 账户检查
 > 按以上 "发行计划"，检查六个账户余额是否正常

- 端口检查
> 26657 端口只允许监听在 localhost地址

## keystores 备份

```
    #备份keystores
    把 ~/config 目录中的keystores 子目录，拷贝到U盘
    
    #删除keystores
    先确认以上备份做好
    ~/config 目录 只保留 .pem 文件，其他文件和子目录要删除
```





