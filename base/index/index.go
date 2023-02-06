package index

import (
	"cess-indexer/base/chain"
	"cess-indexer/base/client"
	"cess-indexer/config"
	"cess-indexer/logger"
	"cess-indexer/utils"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

const (
	CACHERS_FRESHTIME = time.Second * 10
	TX_POOL_FRESHTIME = time.Second * 6
	DEFAULT_POOL_SIZE = 128
	MAX_POOL_SIZE     = 1024
)

type BID [16]byte

type CacherFactor struct {
	Account   types.AccountID `json:"account"`
	BytePrice chain.Balance   `json:"bytePrice"`
	Distance  float64         `json:"distance"`
}

type CachersHandle interface {
	Load(key any) (value any, ok bool)
	Store(key any, value any)
	Delete(key any)
	Range(f func(key any, value any) bool)
}

type TxPoolHandle interface {
	Put(bill chain.Bill)
	GetBills() []chain.Bill
	GetTxHash(bid BID) (string, bool)
	Store(hash string, keys ...BID)
	BillsNum() int
}

type TxPool struct {
	records map[BID]string
	mu      sync.RWMutex
	pool    chan chain.Bill
}

var (
	cachers sync.Map
	orders  TxPool
)

func GetCachers() CachersHandle {
	return &cachers
}

func GetOrders() TxPoolHandle {
	return &orders
}

func InitIndexer() {
	maxSize := config.GetConfig().MaxTxPoolSize
	if maxSize <= 0 {
		maxSize = DEFAULT_POOL_SIZE
	} else if maxSize > MAX_POOL_SIZE {
		maxSize = MAX_POOL_SIZE
	}
	orders = TxPool{
		records: make(map[BID]string),
		pool:    make(chan chain.Bill, maxSize),
	}
}

func (p *TxPool) Put(bill chain.Bill) {
	p.pool <- bill
}

func (p *TxPool) GetBills() []chain.Bill {
	p.mu.Lock()
	defer p.mu.Unlock()
	num := len(p.pool)
	bills := make([]chain.Bill, num)
	for i := 0; i < num; i++ {
		bills[i] = <-p.pool
	}
	return bills
}
func (p *TxPool) GetTxHash(bid BID) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	v, ok := p.records[bid]
	return v, ok
}
func (p *TxPool) Store(hash string, keys ...BID) {
	p.mu.Lock()
	for _, key := range keys {
		p.records[key] = hash
	}
	defer p.mu.Unlock()
}

func (p *TxPool) BillsNum() int {
	return len(p.pool)
}

func BillServer() {
	fullSignal := make(chan struct{}, 1)
	go func() {
		for {
			if GetOrders().BillsNum() >= MAX_POOL_SIZE {
				fullSignal <- struct{}{}
			}
		}
	}()
	ticker := time.NewTicker(TX_POOL_FRESHTIME)
	for {
		select {
		case <-ticker.C:
		case <-fullSignal:
		}
		go func() {
			bills := GetOrders().GetBills()
			txhash, err := chain.GetChainCli().CreateAndSendCacheBills(bills)
			if err != nil {
				logger.Uld.Sugar().Error(err)
				return
			}
			bids := make([]BID, len(bills))
			for i, bill := range bills {
				bids[i] = bill.Id
			}
			GetOrders().Store(txhash, bids...)
		}()
		if len(fullSignal) == 1 {
			<-fullSignal
		}
	}
}

func GetBestCacher(addrs ...string) string {
	var minPrice, maxPrice uint64 = math.MaxUint64, 0
	var minDist, maxDist float64 = math.MaxFloat64, 0
	cfs := make(map[string]CacherFactor)
	for _, addr := range addrs {
		v, ok := GetCachers().Load(addr)
		if !ok {
			continue
		}
		cf := v.(CacherFactor)
		cfs[addr] = cf
		if cf.BytePrice.Uint64() < minPrice {
			minPrice = cf.BytePrice.Uint64()
		}
		if cf.BytePrice.Uint64() > maxPrice {
			maxPrice = cf.BytePrice.Uint64()
		}
		if cf.Distance < minDist {
			minDist = cf.Distance
		}
		if cf.Distance > maxDist {
			maxDist = cf.Distance
		}
	}
	var bastCacher string
	var bastScore float64
	for k, v := range cfs {
		ds := 10000 * (1 - math.Abs(v.Distance-minDist)/maxDist)
		ps := 10000 * (1 - math.Abs(float64(v.BytePrice.Uint64()-minPrice)/float64(maxPrice)))
		if ds+ps > bastScore {
			bastCacher = k
			bastScore = ps + ds
		}
	}
	return bastCacher
}

func WatchCachers() {
	ticker := time.NewTicker(CACHERS_FRESHTIME)
	for range ticker.C {
		cachers, err := chain.GetChainCli().GetCachers()
		if err != nil {
			continue
		}
		for _, c := range cachers {
			ip := fmt.Sprintf("%d.%d.%d.%d",
				c.Ip.IPv4.Value[0],
				c.Ip.IPv4.Value[1],
				c.Ip.IPv4.Value[2],
				c.Ip.IPv4.Value[3],
			)
			port := fmt.Sprint(c.Ip.IPv4.Port)
			if v, ok := GetCachers().Load(ip + ":" + port); ok {
				cf := v.(CacherFactor)
				if cf.BytePrice.Cmp(c.Byte_price.Int) != 0 {
					cf.BytePrice = c.Byte_price
					GetCachers().Store(ip+":"+port, cf)
				}
				continue
			}
			dist := calcDist(ip)
			if dist < 0 {
				logger.Uld.Sugar().Errorf("can not calc peer's distance,addr:", ip+":"+port)
				continue
			}
			cf := CacherFactor{
				Account:   c.Acc,
				BytePrice: c.Byte_price,
				Distance:  dist,
			}
			GetCachers().Store(ip+":"+port, cf)
		}
	}
}

func calcDist(addr string) float64 {
	ip1 := strings.Split(addr, ":")[0]
	lat1, lng1, err := utils.ParseLocationByIp(ip1)
	if err != nil {
		return -1
	}
	ip2, err := utils.GetExternalIp()
	if err != nil {
		return -1
	}
	lat2, lng2, err := utils.ParseLocationByIp(ip2)
	if err != nil {
		return -1
	}
	return utils.EarthDistance(lat1, lng1, lat2, lng2)
}

func CreateCacheBill(cacher, fhash, shash string, size uint64, exp time.Duration) (chain.Bill, error) {
	v, ok := GetCachers().Load(cacher)
	if !ok {
		err := errors.New("cacher not fount")
		return chain.Bill{}, errors.Wrap(err, "create cache bill error")
	}
	cf := v.(CacherFactor)
	var b big.Int
	hfh, err := hex.DecodeString(fhash)
	if err != nil {
		return chain.Bill{}, errors.Wrap(err, "create cache bill error")
	}
	hsh, err := hex.DecodeString(shash)
	if err != nil {
		return chain.Bill{}, errors.Wrap(err, "create cache bill error")
	}
	return chain.Bill{
		Id:              utils.CreateUUID(),
		To:              cf.Account,
		Amount:          types.NewU128(*b.Mul(cf.BytePrice.Int, b.SetUint64(size))),
		File_hash:       types.NewHash(hfh),
		Slice_hash:      types.NewHash(hsh),
		Expiration_time: uint64(time.Now().Add(exp).Unix()),
	}, nil
}

func CreateAndSendCacheBill(fhash, shash string, size uint64, exp time.Duration) (BID, string, error) {
	cachers, err := client.QueryMinersCachedFile(fhash, GetCachers())
	if err != nil {
		return BID{}, "", errors.Wrap(err, "create and send cache bill error")
	}
	if len(cachers) <= 0 {
		cachers = getAllCaches()
	}
	cacher := GetBestCacher(cachers...)
	bill, err := CreateCacheBill(cacher, fhash, shash, size, exp)
	if err != nil {
		return BID{}, "", errors.Wrap(err, "create and send cache bill error")
	}
	GetOrders().Put(bill)
	return bill.Id, cacher, nil
}

func getAllCaches() []string {
	var cachers []string
	GetCachers().Range(func(key, value any) bool {
		cachers = append(cachers, key.(string))
		return true
	})
	return cachers
}
