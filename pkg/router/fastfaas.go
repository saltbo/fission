package router

import (
	"context"
	"github.com/gw123/glog"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var fastfaasSvcAddr = ""

func init() {
	fastfaasSvcAddr = os.Getenv("FASTFAAS_SVC_ADDR")
	if fastfaasSvcAddr == "" {
		glog.Warn("get FASTFAAS_SVC_ADDR from env faild")
		fastfaasSvcAddr = "http://fastfaas"
	}
}

var Client *http.Client = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: 180 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxConnsPerHost:       2000,
		MaxIdleConnsPerHost:   50,
		IdleConnTimeout:       180 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func GetServiceForFunction(ctx context.Context, funcNname string) (*url.URL, error) {
	ctx, _ = context.WithTimeout(ctx, time.Second*10)
	executorURL := fastfaasSvcAddr + "/v3/getServiceForFunction/" + funcNname
	resp, err := ctxhttp.Get(ctx, Client, executorURL)
	if err != nil {
		return nil, errors.Wrap(err, "error posting to getting service for function")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("resp statuscode %d", resp.StatusCode)
	}

	svcName, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body from getting service for function")
	}

	svcURL, err := url.Parse(string(svcName))
	if err != nil {
		glog.Errorf("error parsing service url",
			zap.Error(err),
			zap.String("service_url", svcURL.String()))
		return nil, err
	}

	return svcURL, nil
}

func UnTapService(funcName, addr string) error {
	executorURL := fastfaasSvcAddr + "/v3/unTapService?funcName=" + funcName + "&addr=" + addr
	resp, err := ctxhttp.Post(context.Background(), Client, executorURL, "application/json", nil)
	if err != nil {
		glog.WithErr(err).Errorf("UnTapService err")
		return errors.Wrap(err, "error posting to getting service for function")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		glog.Errorf("UnTapService err, code %d", resp.StatusCode)
		return errors.Errorf("resp statuscode %d", resp.StatusCode)
	}

	return nil
}

func TapService(funcName, addr string) error {
	executorURL := fastfaasSvcAddr + "/v3/tapService?funcName=" + funcName + "&addr=" + addr
	resp, err := ctxhttp.Post(context.Background(), Client, executorURL, "application/json", nil)
	if err != nil {
		return errors.Wrap(err, "error posting to getting service for function")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("resp statuscode %d", resp.StatusCode)
	}
	return nil
}
