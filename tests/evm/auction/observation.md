## experiments
#### single node test
    1. first-after bid error
    bid failure for all accounts except first two genesis accounts.
    first-after bidding invalid. reverting for second genesis account.
    transaction created, but highestBid, highestBidder doesn't update.

    [log analysis]
    op:REVERT
    contract err:evm: execution reverted
    Reverting Snapshot
    balanceChange reverting...
    step into account.SetBalance.

    [Caution]
    auctionEnd Time Out!!!
    SOLVED!!!!

    2. non-genesis account bid failure]
    returns no found transaction!
    it reports "recovered: invalid account type for state object"

    [panic analysis]
    state.newObject
    state.(*CommitStateDB).getStateObject
    state.(*CommitStateDB).GetBalance
    app.(*SscqServiceApp).openContract
    app.(*SscqServiceApp).Transition
    baseapp.(*BaseApp).runMsgs
    baseapp.(*BaseApp).runTx
    baseapp.(*BaseApp).DeliverTx
    client.(*localClient).DeliverTxAsync
    proxy.(*appConnConsensus).DeliverTxAsync
    state.execBlockOnProxyApp
    state.(*BlockExecutor).ApplyBlock
    consensus.(*ConsensusState).finalizeCommit
    consensus.(*ConsensusState).tryFinalizeCommit
    consensus.(*ConsensusState).enterCommit
    consensus.(*ConsensusState).addVote
    consensus.(*ConsensusState).tryAddVote
    consensus.(*ConsensusState).handleMsg
    consensus.(*ConsensusState).receiveRoutine
    consensus.(*ConsensusState).OnStart

#### Questions & Analysis
    1. now(how to know current time in solidity on blockchain)
    2. difference in statedb between genesis accounts & non-genesis accounts
    - log analysis & comparison(openContract to getStateObject)
    - statedb difference
