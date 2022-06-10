package async_func

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gw123/glog"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
)

type NextFunc struct {
	DelayLevel int    `json:"delay_level"`
	FuncId     string `json:"function_id"`
	Body       []byte `json:"body"`
}

type NextEvent struct {
	DelayLevel int    `json:"delay_level"`
	Event      string `json:"event"`
	Body       []byte `json:"body"`
}

type NextCall struct {
	NextFuncs  []NextFunc  `json:"next_functions"`
	NextEvents []NextEvent `json:"next_events"`
}

func AsyncCall(ctx context.Context, functionName string, resp *http.Response) {
	headers := resp.Header
	isCallNext := headers.Get("IsCallNext")
	if isCallNext == "true" && resp.StatusCode == http.StatusOK {
		bodyRes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			glog.WithOTEL(ctx).WithError(err).Errorf("asyncCall err")
			return
		}
		resbody := ioutil.NopCloser(bytes.NewReader(bodyRes))
		resp.Body = resbody

		var next NextCall
		if err := json.Unmarshal(bodyRes, &next); err != nil {
			glog.WithOTEL(ctx).WithError(err).Errorf("parse next call err, next functions %s", isCallNext)
			return
		}

		aGW := os.Getenv("AsyncGatewayAddr")
		if aGW == "" {
			aGW = "http://async-gateway.fission.svc/async"
		}
		url := aGW + "/callNext"

		if err := DoPost(next, nil, url, functionName); err != nil {
			glog.WithOTEL(ctx).WithError(err).Errorf("send msg to asyncGateway parse err, next functions %s", isCallNext)
			return
		}
	}
}

type Error struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type RestResponseWithTraceId struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []Error     `json:"errors,omitempty"`
	Total   int64       `json:"total,omitempty" `
	TraceId string      `json:"trace_id,omitempty" `
}

func DoPost(reqData, respData interface{}, url, source string) error {
	data, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("X-Source", source)

	// req.Header.Set("Authorization", viper.GetString("token"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("status is %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cresp := &RestResponseWithTraceId{}
	if err := json.Unmarshal(body, cresp); err != nil {
		return err
	}

	if cresp.Code != 0 || len(cresp.Errors) > 0 {
		return errors.Errorf("%+v", cresp.Errors[0])
	}

	if respData == nil {
		return nil
	}

	if err := json.Unmarshal(body, respData); err != nil {
		return err
	}
	return nil
}
