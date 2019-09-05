package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/imroc/req"
	"net/http"
	"time"
)

func BlockHash(bts []byte) string {
	h := sha1.New()
	h.Write(bts)
	return hex.EncodeToString(h.Sum(nil))
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func SetConnPool() {
	client := &http.Client{}
	client.Transport = &http.Transport{
		MaxIdleConnsPerHost: 500,
		// 无需设置MaxIdleConns
		// MaxIdleConns controls the maximum number of idle (keep-alive)
		// connections across all hosts. Zero means no limit.
		// MaxIdleConns 默认是0，0表示不限制
	}

	req.SetClient(client)
	req.SetTimeout(5 * time.Second)
}
