package client

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/valyala/fasthttp"
)

type FastHTTPClient struct {
	*fasthttp.HostClient
	*pool
}

func NewClient(addr string, maxCap, poolCap uint64) *FastHTTPClient {
	c := &fasthttp.HostClient{
		Addr: addr,
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialTimeout(addr, 2*time.Second)
		},
		MaxIdleConnDuration: 1 * time.Second,
		MaxConns:            int(poolCap),
	}
	return &FastHTTPClient{c, newPool(maxCap, poolCap)}
}

func (c *FastHTTPClient) GetBlock(url string, i uint64) (*Response, error) {
	body := c.getBytes()

	statusCode, body, err := c.Get(body, url)
	if err != nil {
		c.putBytes(body)
		return nil, fmt.Errorf("Error when loading page %s through local proxy: %v", url, err)
	}
	if statusCode != fasthttp.StatusOK {
		c.putBytes(body)
		return nil, fmt.Errorf("Unexpected status code: %d. Expecting %d, url: %s", statusCode, fasthttp.StatusOK, url)
	}

	r := &Response{}
	err = json.Unmarshal(body, r)
	if err != nil {
		c.putBytes(body)
		return nil, fmt.Errorf("Failed to unmarshal response body %s, url: %s, error: %v", string(body), url, err)
	}

	c.putBytes(body)
	return r, nil
}
