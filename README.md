# algotx
Tool to perform Algorand transactions

## disclaimer
> The use of this tool does not guarantee security or usability for any
> particular purpose. Please review the code and use at your own risk.

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system.

```bash
go install github.com/kubetrail/algotx@latest
```
Add autocompletion for `bash` to your `.bashrc`
```bash
source <(algotx completion bash)
```

## runtime prerequisite
`algotx` acts as a client to `algod` daemon, which can run on localhost. In order
to run `algod` on localhost please follow the instructions
[hre](https://github.com/algorand/sandbox)
```bash
git clone https://github.com/algorand/sandbox.git
cd sandbox
./sandbox up testnet
```
You should see output similar to:
```text
Starting sandbox for: testnet
see sandbox.log for detailed progress, or use -v.
* docker containers started!
* waiting for services to initialize.
* services ready!

algod version
12885295106
3.6.2.stable [rel/stable] (commit #5e6bc6fc)
go-algorand is licensed with AGPLv3.0
source code available at https://github.com/algorand/go-algorand

Indexer version
Indexer disabled for this configuration.

Postgres version
postgres (PostgreSQL) 13.6

algod - goal node status
Last committed block: 15
Time since last block: 0.0s
Sync Time: 0.7s
Last consensus protocol: https://github.com/algorand/spec/tree/a26ed78ed8f834e2b9ccb6eb7d3ee9f629a6e622
Next consensus protocol: https://github.com/algorand/spec/tree/a26ed78ed8f834e2b9ccb6eb7d3ee9f629a6e622
Round for next consensus protocol: 16
Next consensus protocol supported: true
Last Catchpoint: 
Genesis ID: testnet-v1.0
Genesis hash: SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=

indexer - health
Indexer disabled for this configuration.

Starting fast-catchup with catchpoint: 21580000#PAWLVBLZ7ABAVWDNYE3R6AO7WYIOXYCVMJM2NDY7APXCXUYZWZVA
* Account processing complete.
* Blocks downloaded.

Fast-catchup complete! Printing status...

algod - goal node status
Last committed block: 21580042
Time since last block: 0.0s
Sync Time: 0.7s
Last consensus protocol: https://github.com/algorandfoundation/specs/tree/d5ac876d7ede07367dbaa26e149aa42589aac1f7
Next consensus protocol: https://github.com/algorandfoundation/specs/tree/d5ac876d7ede07367dbaa26e149aa42589aac1f7
Round for next consensus protocol: 21580043
Next consensus protocol supported: true
Last Catchpoint: 
Genesis ID: testnet-v1.0
Genesis hash: SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=

indexer - health
Indexer disabled for this configuration.
```

At this point RPC calls to `testnet` should go through.

## generate keys
[bip39](https://github.com/kubetrail/bip39) and
[algokey](https://github.com/kubetrail/algokey) can be used to generate new Algorand keys. 
The `mnemonic`, generated from `bip39` can optionally be securely saved on Google
secrets engine using
[mksecret](https://github.com/kubetrail/mksecret)

Following commands generate separate keys and addresses for the
sender and the receiver.

For sender:
```bash
bip39 gen \
  | mksecret set --name=algorand-sender \
  | algokey gen --output-format=yaml
```
```yaml
seed: 301915a64d9745206b2cfd11f5c7d77b5f912d1ac2ad14b573881ea780ff0bfd
prvHex: 301915a64d9745206b2cfd11f5c7d77b5f912d1ac2ad14b573881ea780ff0bfdd02eb69ff2188fed1b708538e1d280294e372288ca3b666d5127c395619b35d7
pubHex: d02eb69ff2188fed1b708538e1d280294e372288ca3b666d5127c395619b35d7
addr: 2AXLNH7SDCH62G3QQU4ODUUAFFHDOIUIZI5WM3KRE7BZKYM3GXLWALTFDA
keyType: ed25519
```

Similarly, for the receiver:
```bash
bip39 gen \
  | mksecret set --name=algorand-receiver \
  | algokey gen --output-format=yaml
```
```yaml
seed: 2db65f01c25bb9dbc25933700f9eac4e195611dd5df70337ca37f2667b9083ea
prvHex: 2db65f01c25bb9dbc25933700f9eac4e195611dd5df70337ca37f2667b9083eab810cc886cacc7066ab60dd6ce348da027b1ace47e4f1efb1750660836430b75
pubHex: b810cc886cacc7066ab60dd6ce348da027b1ace47e4f1efb1750660836430b75
addr: XAIMZCDMVTDQM2VWBXLM4NENUAT3DLHEPZHR56YXKBTAQNSDBN2WKXK25Q
keyType: ed25519
```

Now that we have two addresses and the private key for the sender, we can invoke
a transaction.

However, senders account first needs to be funded, which can be done
using a [testnet faucet](https://bank.testnet.algorand.network/)

## view balance
View balance of the sender:
```bash
algotx balance 2AXLNH7SDCH62G3QQU4ODUUAFFHDOIUIZI5WM3KRE7BZKYM3GXLWALTFDA
```
```text
15000000
```

Confirm the receiver balance as well:
```bash
algotx balance XAIMZCDMVTDQM2VWBXLM4NENUAT3DLHEPZHR56YXKBTAQNSDBN2WKXK25Q
```
```text
0
```

## send transaction
```bash
algotx send \
  --addr=XAIMZCDMVTDQM2VWBXLM4NENUAT3DLHEPZHR56YXKBTAQNSDBN2WKXK25Q \
  --key=301915a64d9745206b2cfd11f5c7d77b5f912d1ac2ad14b573881ea780ff0bfdd02eb69ff2188fed1b708538e1d280294e372288ca3b666d5127c395619b35d7 \
  --amount=5000000
```

This will produce a transaction ID.

At this point the sender will have a balance of `9999000` and the receiver will have a 
balance of `5000000`

## stopping node
Go back to the folder where 
[sandbox](https://github.com/algorand/sandbox) 
code was cloned, then run
```bash
./sandbox down
```
```text
Stopping sandbox containers...
[+] Running 3/3
 ⠿ Container algorand-sandbox-indexer   Stopped                                                                                              0.2s
 ⠿ Container algorand-sandbox-algod     Stopped                                                                                              0.2s
 ⠿ Container algorand-sandbox-postgres  Stopped
```

At this point running an RPC command will result in:
```text
Error: failed to get account info: Get "http://localhost:4001/v2/accounts/XAIMZCDMVTDQM2VWBXLM4NENUAT3DLHEPZHR56YXKBTAQNSDBN2WKXK25Q": dial tcp [::1]:4001: connect: connection refused
```
