package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Pairs map[string]string

type AuthReq struct {
	BID  string `json:"bid"`
	Sign []byte `json:"sign"`
}

func GenerateToken(addr, bid string, sign []byte) (string, error) {
	var token string
	u, err := url.JoinPath("http://"+addr, "gen")
	if err != nil {
		return token, errors.Wrap(err, "generate token error")
	}
	headers := Pairs{"Content-Type": "application/json"}
	jbytes, err := json.Marshal(AuthReq{bid, sign})
	if err != nil {
		return token, errors.Wrap(err, "generate token error")
	}
	resp, err := HttpRequest("POST", u, headers, strings.NewReader(string(jbytes)))
	if err != nil {
		return token, errors.Wrap(err, "generate token error")
	}
	if PraseData(&token, resp.Data) != nil {
		return token, errors.Wrap(err, "generate token error")
	}
	return token, nil
}

func QueryMinerStats(addr string) (MinerStats, error) {
	var stats MinerStats
	u, err := url.JoinPath("http://"+addr, "query", "stats")
	if err != nil {
		return stats, errors.Wrap(err, "query miner stats error")
	}
	resp, err := HttpRequest("GET", u, make(Pairs), nil)
	if err != nil {
		return stats, errors.Wrap(err, "query miner stats error")
	}
	if PraseData(&stats, resp.Data) != nil {
		return stats, errors.Wrap(err, "query miner stats error")
	}
	return stats, nil
}

func QueryMinerCacheFiles(addr string) ([]string, error) {
	var files []string
	u, err := url.JoinPath("http://"+addr, "query", "cached")
	if err != nil {
		return files, errors.Wrap(err, "query cached files error")
	}
	resp, err := HttpRequest("GET", u, make(Pairs), nil)
	if err != nil {
		return files, errors.Wrap(err, "query cached files error")
	}
	if PraseData(&files, resp.Data) != nil {
		return files, errors.Wrap(err, "query cached files error")
	}
	return files, nil
}

func QueryMinerFile(addr, hash string) (FileStat, error) {
	var stats FileStat
	u, err := url.JoinPath("http://"+addr, "query", "file", hash)
	if err != nil {
		return stats, errors.Wrap(err, "query file stats error")
	}
	resp, err := HttpRequest("GET", u, make(Pairs), nil)
	if err != nil {
		return stats, errors.Wrap(err, "query file stats error")
	}
	if PraseData(&stats, resp.Data) != nil {
		return stats, errors.Wrap(err, "query file stats error")
	}
	return stats, nil
}

func HttpRequest(method string, url string, headers Pairs, body io.Reader) (Response, error) {
	var result Response
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return result, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, errors.Wrap(err, "send http request error")
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, errors.Wrap(err, "send http post request error")
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result, errors.Wrap(err, "send http post request error")
	}
	return result, nil
}

func PraseData(target any, data any) error {
	jbytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jbytes, target)
}
