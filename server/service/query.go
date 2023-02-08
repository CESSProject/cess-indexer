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
