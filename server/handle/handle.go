package handle

import (
	resp "cess-indexer/server/response"
	"cess-indexer/server/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderReq struct {
	FileHash  string `json:"fid"`
	SliceHash string `json:"sname,omitempty"`
	Index     string `json:"sindex,omitempty"`
}

type TokenReq struct {
	Cacher string `json:"cacher"`
	BID    string `json:"BID"`
}

func CreateBillHandler(c *gin.Context) {
	var orderReq OrderReq
	e := c.BindJSON(&orderReq)
	if e != nil {
		resp.RespError(c, resp.NewError(400, errors.New("bad params")))
		return
	}
	if orderReq.FileHash == "" || (orderReq.SliceHash == "" && orderReq.Index == "") {
		resp.RespError(c, resp.NewError(400, errors.New("bad params")))
		return
	}
	if orderReq.SliceHash != "" {
		res, err := service.CreateCacheBill(orderReq.FileHash, orderReq.SliceHash)
		if err != nil {
			resp.RespError(c, err)
			return
		}
		resp.RespOk(c, res)
		return
	}
	i, e := strconv.Atoi(orderReq.Index)
	if e != nil {
		resp.RespError(c, resp.NewError(400, e))
		return
	}
	res, err := service.CreateCacheBillBySliceIndex(orderReq.FileHash, i)
	if err != nil {
		resp.RespError(c, err)
		return
	}
	resp.RespOk(c, res)
}

func GenerateFileTokenHandler(c *gin.Context) {
	var tokenReq TokenReq
	c.BindJSON(&tokenReq)
	if tokenReq.Cacher == "" || tokenReq.BID == "" {
		resp.RespError(c, resp.NewError(400, errors.New("bad params")))
		return
	}
	res, err := service.GenerateFileToken(tokenReq.Cacher, tokenReq.BID)
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
