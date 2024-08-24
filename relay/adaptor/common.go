package adaptor

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	//"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common/client"
	"github.com/songquanpeng/one-api/relay/meta"
)

func SetupCommonRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) {
	req.Header.Set("Content-Type", c.Request.Header.Get("Content-Type"))
	req.Header.Set("Accept", c.Request.Header.Get("Accept"))
	if meta.IsStream && c.Request.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/event-stream")
	}
}

// 辅助函数，用于读取和返回请求体内容
func readRequestBody(req *http.Request) string {
	var bodyBytes []byte
	bodyBytes, _ = io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重置请求体，以便可以再次使用
	return string(bodyBytes)
}

func DoRequestHelper(a Adaptor, c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	fullRequestURL, err := a.GetRequestURL(meta)
	if err != nil {
		return nil, fmt.Errorf("get request url failed: %w", err)
	}
	req, err := http.NewRequest(c.Request.Method, fullRequestURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("new request failed: %w", err)
	}
	err = a.SetupRequestHeader(c, req, meta)
	if err != nil {
		return nil, fmt.Errorf("setup request header failed: %w", err)
	}
	// 打印请求对象的详细信息
	fmt.Printf("Request Method: %s\n", req.Method)
	fmt.Printf("Request URL: %s\n", req.URL.String())
	fmt.Printf("Request Headers: %+v\n", req.Header)
	fmt.Printf("Request Body: %s\n", readRequestBody(req))

	resp, err := DoRequest(c, req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	return resp, nil
}

// 添加这个新函数
func readResponseBody(resp *http.Response) string {
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}

func DoRequest(c *gin.Context, req *http.Request) (*http.Response, error) {
	resp, err := client.HTTPClient.Do(req)
	// 打印响应对象的详细信息
	fmt.Printf("Response Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Headers: %+v\n", resp.Header)
	fmt.Printf("Response Body: %s\n", readResponseBody(resp))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("resp is nil")
	}
	_ = req.Body.Close()
	_ = c.Request.Body.Close()
	return resp, nil
}
