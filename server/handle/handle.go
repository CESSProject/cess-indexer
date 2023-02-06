package handle

import (
	resp "cess-indexer/server/response"
	"cess-indexer/server/service"
	"errors"

	"github.com/gin-gonic/gin"
)

func CreateBillHandler(c *gin.Context) {
	fhash := c.PostForm("filehash")
	shash := c.PostForm("slicehash")
	if fhash == "" || shash == "" {
		resp.RespError(c, resp.NewError(400, errors.New("bad params")))
		return
	}
	res, err := service.CreateCacheBill(fhash, shash)
	if err != nil {
		resp.RespError(c, err)
		return
	}
	resp.RespOk(c, res)
}

func GenerateFileTokenHandler(c *gin.Context) {
	cacher := c.PostForm("cacher")
	bid := c.PostForm("id")
	if cacher == "" || bid == "" {
		resp.RespError(c, resp.NewError(400, errors.New("bad params")))
		return
	}
	res, err := service.GenerateFileToken(cacher, bid)
	if err != nil {
		if err.Status() == 0 {
			resp.RespOkWithFlag(c, false, res)
			return
		}
		resp.RespError(c, err)
		return
	}
	resp.RespOk(c, res)
}

func QueryHandler(c *gin.Context) {
	resp.RespOk(c, service.QueryAllCachers())
}
