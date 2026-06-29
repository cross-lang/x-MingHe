package httpclient

import (
	"portal/internal/pkg/log"

	"go.uber.org/zap"

	"time"

	"github.com/valyala/fasthttp"
)

func FastHttpPost(url string, header map[string]string, body []byte, timeout int) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源
	// 默认是PostDeactivateUserlication/x-www-form-urlencoded
	req.Header.SetMethod("POST")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.SetRequestURI(url)
	req.SetBody(body)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源
	var err error
	if timeout > 0 {
		if err = fasthttp.Do(req, resp); err != nil {
			log.Logger.Debug("http请求失败", zap.Any("err", err))
			return nil, err
		}
	} else {
		if err := fasthttp.DoTimeout(req, resp, time.Duration(timeout)*time.Second); err != nil {
			log.Logger.Debug("http请求失败", zap.Any("err", err))
			return nil, err
		}
	}
	b := resp.Body()
	return b, nil
}

func FastHttpPostForm(url string, header map[string]string, body map[string]string, timeout int) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源
	// 默认是PostDeactivateUserlication/x-www-form-urlencoded
	req.Header.SetMethod("POST")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.SetRequestURI(url)
	args := &fasthttp.Args{}
	for k, v := range body {
		args.Add(k, v)
	}
	req.SetBody(args.QueryString())
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源
	var err error
	if timeout == 0 {
		if err = fasthttp.Do(req, resp); err != nil {
			log.Logger.Debug("http请求失败", zap.Any("err", err))
			return nil, err
		}
	} else {
		if err := fasthttp.DoTimeout(req, resp, time.Duration(timeout)*time.Second); err != nil {
			log.Logger.Debug("http请求失败", zap.Any("err", err))
			return nil, err
		}
	}
	b := resp.Body()
	return b, nil
}
func FastHttpGet(url string, header map[string]string, body map[string]string, timeout int) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源
	// 默认是PostDeactivateUserlication/x-www-form-urlencoded
	req.Header.SetMethod("GET")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	if len(body) > 0 {
		url += "?"
		for k, v := range body {
			url += k + "=" + v + "&"
		}
		url = url[0 : len(url)-1]
	}
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源
	var err error
	if timeout == 0 {
		if err = fasthttp.Do(req, resp); err != nil {
			log.Logger.Debug("http请求失败", zap.Any("err", err))
			return nil, err
		}
	} else {
		if err := fasthttp.DoTimeout(req, resp, time.Duration(timeout)*time.Second); err != nil {
			log.Logger.Debug("http请求失败", zap.Any("err", err))
			return nil, err
		}
	}
	b := resp.Body()
	return b, nil
}
