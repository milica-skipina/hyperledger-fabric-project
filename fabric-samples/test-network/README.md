# PDASP

To run application open terminal from root folder and run:
```
$ cd fabric-samples/test-network 
$ ./network.sh down
$ ./network.sh up createChannel -ca
$ ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
$ cd ../asset-transfer-basic/application-javascript
$ rm -rf ./wallet
$ node app.js
```
