package service

import (
	"cess-indexer/base/chain"
	"cess-indexer/base/client"
	"cess-indexer/base/index"
	"cess-indexer/config"
	resp "cess-indexer/server/response"
	"encoding/hex"
	"time"

	"github.com/CESSProject/go-keyring"

	"github.com/pkg/errors"
)

func CreateCacheBill(fhash, shash string) (map[string]string, resp.Error) {
	resMap := make(map[string]string)
	fmeta, err := chain.GetTestChainCli().GetFileMetaInfo(fhash)
	if err != nil {
		return resMap, resp.NewError(500, errors.Wrap(err, "create cache bill error"))
	}
	var size uint64
	if isInvalidFileSlice(shash, fmeta, &size) {
		err = errors.New("bad file slice hash")
		return resMap, resp.NewError(400, errors.Wrap(err, "create cache bill error"))
	}
	bid, cacher, err := index.CreateAndSendCacheBill(
		fhash, shash, uint64(size),
		time.Duration(config.GetConfig().ExpDate)*time.Hour,
	)
	if err != nil {
		return resMap, resp.NewError(400, errors.Wrap(err, "create cache bill error"))
	}
	resMap["BID"] = hex.EncodeToString(bid[:])
	resMap["cacherAddr"] = cacher
	resMap["tips"] = "please get token after 15 s"
	return resMap, nil
}

func CreateCacheBillBySliceIndex(fhash string, i int) (map[string]string, resp.Error) {
	resMap := make(map[string]string)
	fmeta, err := chain.GetTestChainCli().GetFileMetaInfo(fhash)
	if err != nil {
		return resMap, resp.NewError(500, errors.Wrap(err, "create cache bill error"))
	}
	if i >= len(fmeta.BlockInfo) || i < 0 {
		err = errors.New("invalid  file slice index")
		return resMap, resp.NewError(400, errors.Wrap(err, "create cache bill error"))
	}
	bid, cacher, err := index.CreateAndSendCacheBill(
		fhash, string(fmeta.BlockInfo[i].BlockId[:]), uint64(fmeta.BlockInfo[i].BlockSize),
		time.Duration(config.GetConfig().ExpDate)*time.Hour,
	)
	if err != nil {
		return resMap, resp.NewError(400, errors.Wrap(err, "create cache bill error"))
	}
	resMap["BID"] = hex.EncodeToString(bid[:])
	resMap["cacherAddr"] = cacher
	resMap["tips"] = "please get token after 15 s"
	return resMap, nil
}

func GenerateFileToken(cacher, hexBid string) (map[string]string, resp.Error) {
	resMap := make(map[string]string)
	b, err := hex.DecodeString(hexBid)
	if err != nil {
		return resMap, resp.NewError(400, errors.Wrap(err, "generate file token error"))
	}
	bid := index.BID{}
	copy(bid[:], b)
	hash, ok := index.GetOrders().GetTxHash(bid)
	if !ok {
		resMap["tips"] = "the bill is not ready"
		return resMap, resp.NewError(0, nil)
	}
	kr, _ := keyring.FromURI(config.GetConfig().AccountSeed, keyring.NetSubstrate{})
	// sign message
	sign, err := kr.Sign(kr.SigningContext([]byte(hash + hexBid)))
	if err != nil {
		return resMap, resp.NewError(400, errors.Wrap(err, "generate file token error"))
	}
	token, err := client.GenerateToken(cacher, hash, hexBid, sign[:])
	if err != nil {
		return resMap, resp.NewError(400, errors.Wrap(err, "generate file token error"))
	}
	resMap["url"] = "http://" + cacher + "/download/file/" + token
	return resMap, nil
}

func isInvalidFileSlice(shash string, fmeta chain.FileMetaInfo, size *uint64) bool {
	for _, slice := range fmeta.BlockInfo {
		if string(slice.BlockId[:]) == shash {
			*size = uint64(slice.BlockSize)
			return false
		}
	}
	return true
}
