package test

import (
	"cess-indexer/base/chain"
	"cess-indexer/base/index"
	"cess-indexer/config"
	"cess-indexer/logger"
	"cess-indexer/server/service"
	"cess-indexer/utils"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

// test chain
var testCli chain.IChain

func GetTestChainCli() chain.IChain {
	return testCli
}

func InitTestChainClient() error {
	var err error
	testCli, err = chain.NewChainClient(
		"wss://testnet-rpc0.cess.cloud/ws/",
		"lunar talent spend shield blade when dumb toilet drastic unique taxi water",
		"cXgZo3RuYkAGhhvCHjAcc9FU13CG44oy8xW6jN39UYvbBaJx5",
		time.Duration(time.Second*15),
	)
	if err != nil {
		return errors.Wrap(err, "init chain client error")
	}
	return nil
}

//

func TestGetCachers(t *testing.T) {
	err := InitTestChainClient()
	if err != nil {
		t.Fatal("init chain client error", err)
	}
	cachers, err := GetTestChainCli().GetCachers()
	if err != nil {
		t.Fatal("get cachers error", err)
	}
	t.Log("cachers info on chain:", cachers)
}

func TestSubmitBills(t *testing.T) {
	err := InitTestChainClient()
	if err != nil {
		t.Fatal("init chain client error", err)
	}
	pubkey := GetTestChainCli().GetPublicKey()
	acc := types.AccountID{}
	copy(acc[:], pubkey)
	txhash, err := GetTestChainCli().CreateAndSendCacheBills(
		[]chain.Bill{
			{
				Id:              utils.CreateUUID(),
				To:              acc,
				File_hash:       types.Hash{},
				Slice_hash:      types.Hash{},
				Amount:          types.NewU128(*big.NewInt(1000)),
				Expiration_time: uint64(time.Now().Add(24 * time.Hour).Unix()),
			},
		},
	)
	if err != nil {
		t.Fatal("submit cache bills on chain error", err)
	}
	t.Log("submit cache bills success,tx hash is", txhash)
}

func TestQueryCachers(t *testing.T) {
	if err := config.InitConfig(""); err != nil {
		log.Fatal("init config error", err)
	}

	logger.InitLogger()

	if err := chain.InitChainClient(config.GetConfig()); err != nil {
		log.Fatal("init chain client error", err)
	}
	index.InitIndexer()
	t.Log(service.QueryAllCachers())
}
