package service

import (
	"cess-indexer/base/index"
)

func QueryAllCachers() map[string]index.CacherFactor {
	res := make(map[string]index.CacherFactor)
	index.GetCachers().Range(func(key, value any) bool {
		res[key.(string)] = value.(index.CacherFactor)
		return true
	})
	return res
}

// func QueryCacherCachedFiles(addr string) ([]string, resp.Error) {
// 	files, err := client.QueryMinerCacheFiles(addr)
// 	if err != nil {
// 		return nil, resp.NewError(500, err)
// 	}
// 	return files, nil
// }

// func QueryFileStat(addr, hash string) (client.FileStat, resp.Error) {
// 	stat, err := client.QueryMinerFile(addr, hash)
// 	if err != nil {
// 		return stat, resp.NewError(500, err)
// 	}
// 	return stat, nil
// }
