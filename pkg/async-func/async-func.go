package async_func

import (
	"bytes"
	"context"
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gw123/glog"
	"io/ioutil"
	"net/http"
)

type NextFunc struct {
	FuncId         string `json:"function_id"`
	Body           []byte `json:"body"`
	ResponseAsBody string `json:"response_as_body"`
}

type NextEvent struct {
	Event          string `json:"event"`
	Body           []byte `json:"body"`
	ResponseAsBody string `json:"response_as_body"`
}

type FuncResp struct {
	NextFuncs  []NextFunc  `json:"next_functions"`
	NextEvents []NextEvent `json:"next_events"`
}

func AsyncCall(ctx context.Context, functionName string, resp *http.Response) {
	headers := resp.Header
	isCallNext := headers.Get("IsCallNext")
	if isCallNext == "true" && resp.StatusCode == http.StatusOK {
		bodyRes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			glog.WithErr(err).Errorf("asyncCall err")
			return
		}
		resbody := ioutil.NopCloser(bytes.NewReader(bodyRes))
		resp.Body = resbody

		var next FuncResp
		if err := json.Unmarshal(bodyRes, &next); err != nil {
			glog.WithErr(err).Errorf("send msg to queue parse err, next functions %s", isCallNext)
			return
		}

		for _, nextFunc := range next.NextFuncs {
			newEvent := cloudevents.NewEvent()
			newEvent.SetSource(functionName)
			newEvent.SetType("function")
			newEvent.SetSubject(nextFunc.FuncId)
			newEvent.SetData(cloudevents.ApplicationCloudEventsJSON, nextFunc.Body)
			SendMessage(ctx, &newEvent)
		}

		for _, nextEvent := range next.NextEvents {
			newEvent := cloudevents.NewEvent()
			newEvent.SetSource(functionName)
			newEvent.SetType(nextEvent.Event)
			newEvent.SetData(cloudevents.ApplicationCloudEventsJSON, nextEvent.Body)
			SendMessage(ctx, &newEvent)
		}
	}
}
