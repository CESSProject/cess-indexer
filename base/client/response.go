package client

type Response struct {
	Result bool `json:"result"`
	Data   any  `json:"data"`
}

type MinerStats struct {
	GeoLocation string      `json:"geoLocation"`
	BytePrice   uint        `json:"bytePrice"`
	MinerStatus string      `json:"status"`
	NetStats    NetStats    `json:"netStats"`
	MemoryStats MemoryStats `json:"memStats"`
	CPUStats    CPUStats    `json:"cpuStats"`
	DiskStats   DiskStats   `json:"diskStats"`
}

type FileStat struct {
	Price      uint     `json:"price"`
	Size       uint64   `json:"size"`
	ShardCount int      `josn:"shardCount"`
	Shards     []string `json:"shards"`
}

type DiskStats struct {
	Total     int64   `json:"total"`
	Used      int64   `json:"used"`
	Available int64   `json:"available"`
	UseRate   float32 `json:"useRate"`
}

type MemoryStats struct {
	Total     int64 `json:"total"`
	Free      int64 `json:"free"`
	Available int64 `json:"available"`
}

type CPUStats struct {
	Num      int     `json:"cpuNum"`
	LoadAvgs float32 `json:"loadAvgs"`
}

type NetStats struct {
	Download int64 `json:"cacheSpeed"`
	Upload   int64 `json:"downloadSpeed"`
}
