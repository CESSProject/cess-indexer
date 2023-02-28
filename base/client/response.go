package client

type Response struct {
	Result bool `json:"result"`
	Data   any  `json:"data"`
}

type MinerStats struct {
	GeoLocation string      `json:"geoLocation"`
	BytePrice   uint64      `json:"bytePrice"`
	MinerStatus string      `json:"status"`
	NetStats    NetStats    `json:"netStats"`
	MemoryStats MemoryStats `json:"memStats"`
	CPUStats    CPUStats    `json:"cpuStats"`
	DiskStats   DiskStats   `json:"diskStats"`
}

type FileStat struct {
	Cached bool   `json:"cached"`
	Price  uint64 `json:"price"`
	Size   uint64 `json:"size"`
}

type DiskStats struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	UseRate   float32 `json:"useRate"`
}

type MemoryStats struct {
	Total     uint64 `json:"total"`
	Free      uint64 `json:"free"`
	Available int64  `json:"available"`
}

type CPUStats struct {
	Num      int     `json:"cpuNum"`
	LoadAvgs float32 `json:"loadAvgs"`
}

type NetStats struct {
	Download uint64 `json:"cacheSpeed"`
	Upload   uint64 `json:"downloadSpeed"`
}

type CacheStat struct {
	HitRate  float32 `json:"hitRate"`
	MissRate float32 `json:"missRate"`
	ErrRate  float32 `json:"errRate"`
}
