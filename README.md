# cess-indexer
## Introduction
indexer is a simple cacher search program. It only needs simple configuration to run the program to obtain services. Please run the program in the Linux operating system. Before doing so, please make sure that the Golang environment is installed on your machine.
## Start
1. First, you need to make a simple configuration. The configuration file is config.toml under the config directory,Please fill in all configuration options.
```yaml
#CESS network ws address
RpcAddr="wss://devnet-rpc.cess.cloud/ws/"
#You indexer server port
ServerPort="8081"
#You CESS account and seed
AccountSeed="plug increase image genius rabbit unable once sponsor tail castle wet fan"
AccountID="cXgP2nH3GkFweiAUcPFWJjxRdiScooTneMgdowKe6Lfe3hQQ8"
#This configuration is used to adjust the size of the transaction pool
MaxTxPoolSize=128
#order expiration(hours)
ExpDate=24
```
2. Then you can go back to the project main directory and use the command `go run main.go` to start the service
## Unit Test
You can use the test samples in the test directory for unit testing. Note that you should set the configuration file before testing
```shell
cd test
go test indexer_test.go
```
## Code Walkthrough
1. The user calls the CreateCacheBill service through the http protocol to generate a new cache order.Then the data will be stored in the cache pool. When the cache pool is full or reaches the predetermined time, it will be submitted to the blockchain together.
https://github.com/CESSProject/cess-indexer/blob/87325c6f5da41e7682d0c065d2f59684dd795ff2/server/service/index.go#L17-L39
2. The user can also call the CreateCacheBillBySliceIndex method to create a cache order. The difference from the former is that the file partition index is used instead of the hash value to complete the location.
https://github.com/CESSProject/cess-indexer/blob/87325c6f5da41e7682d0c065d2f59684dd795ff2/server/service/index.go#L41-L62
The two services are called through the same URL. The only difference is that you choose to use the hash or sine parameter.When creating an bill, all active miners on the blockchain will be queried, and the miners who have cached the file fragments will be preferentially selected to provide download services. If there are multiple miners, the miners with the best unit price and distance will be selected first.
https://github.com/CESSProject/cess-indexer/blob/87325c6f5da41e7682d0c065d2f59684dd795ff2/base/index/index.go#L269-L284
3. After the cache bill is created, users can call the GenerateFileToken service through HTTP request to obtain the download link according to the response data obtained previously.
https://github.com/CESSProject/cess-indexer/blob/87325c6f5da41e7682d0c065d2f59684dd795ff2/server/service/index.go#L64-L89
At this time, the indexer will send a request to the cache to generate a token. After receiving the response data, the indexer will splice the cache file download address and the token into a download link and return.
https://github.com/CESSProject/cess-indexer/blob/87325c6f5da41e7682d0c065d2f59684dd795ff2/base/client/client.go#L26-L45

## API

There are currently two APIs available

### Create cache bill

This API is used to generate cache bill. The indexer will automatically select the nearest and cheapest cacher to provide you with cache services,The use method is as follows:

**POST** /order/bill/create

- Request Header

| key          | value            |
| ------------ | ---------------- |
| Content-Type | application/json |

- Request Body

| field  | value                                                        |
| ------ | ------------------------------------------------------------ |
| fid    | file hash in CESS storage                                    |
| sname  | file slice hash,is optional,you can choose to use sname or sindex |
| sindex | index of file slice ,startiing from 0,                       |

- Responses

Response schema: `application/json`

| HTTP Code | Message     | Description                            |
| --------- | ----------- | -------------------------------------- |
| 200       | message map | return Bill ID,cacher address and tips |
| 400       | error       | specific error reminder                |

### Generate token

This API is used to generate cache file download token. Indexer will return a download link with token for you. You can use this link to download files from cacher on any device.This API needs to use the bill info generated in the previous API.

**POST** /order/token/gen

- Request Header

| key          | value            |
| ------------ | ---------------- |
| Content-Type | application/json |

- Request Body

| field  | value                                                        |
| ------ | ------------------------------------------------------------ |
| cacher | cacher address,please obtain it from 'Create cache bill' API response |
| BID    | bill ID,please obtain it from 'Create cache bill' API response |

- Responses

Response schema: `application/json`

| HTTP Code | Message     | Description                       |
| --------- | ----------- | --------------------------------- |
| 200       | message map | return cache download url or tips |
| 400       | error       | specific error reminder           |
