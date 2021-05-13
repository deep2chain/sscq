### monitoring system
   monitoring system aims to monitor blockchain integrity.
   
### check point
   * [ ] abnomality detection
      * [ ] double sign
      * [ ] breaches
   * [ ] account balance
      * [ ] big fishes
      * [ ] suspicious startup
   * [ ] tps
   * [ ] block time
   * [ ] top accounts

### solution
#### main functions
   * accept signal
   * record
   * analyze
   * alert/report
#### signal handling
   new block signal --> 
   - record block time
   - analyze block
   - record block info(block time)
   - record transaction info
#### database
   * (address,transactions)
   * (blocknumber, blocktime)
   * (blocknumber, txcount)