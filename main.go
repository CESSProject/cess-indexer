package main

import (
	"cess-indexer/base/chain"
	"cess-indexer/base/index"
	"cess-indexer/config"
	"cess-indexer/logger"
	"cess-indexer/server"
	"log"
)

func main() {
	if err := config.InitConfig(""); err != nil {
		log.Fatal("init config error", err)
	}

	logger.InitLogger()

	if err := chain.InitChainClient(config.GetConfig()); err != nil {
		log.Fatal("init chain client error", err)
	}

	//test chain
	if err := chain.InitTestChainClient(); err != nil {
		log.Fatal("init chain client error", err)
	}

	index.InitIndexer()
	server.SetupGinServer()
}
