package async_func

import (
	"context"
	"encoding/json"
	rocketmq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gw123/glog"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type RocketManager struct {
	Producer    rocketmq.Producer
	Nameservers []string
	topicName   string
}

var rocketManger *RocketManager

func init() {
	rockerAddrStr := os.Getenv("ROCKET_ADDRS")
	topicName := os.Getenv("ROCKET_TOPIC_NAME")
	rocketAddrs := strings.Split(rockerAddrStr, ",")
	var err error
	rocketManger, err = NewConsumerManager(rocketAddrs, topicName)
	if err != nil {
		glog.WithErr(err).Error("NewConsumerManager error")
	}
}

func SendMessage(ctx context.Context, event *event.Event) {
	glog.Infof("sendMessage subject:%s type:%s", event.Subject(), event.Type())
	if rocketManger != nil {
		rocketManger.SendMessage(ctx, event)
	} else {
		glog.Error("sendMessage error rocketManager is nil")
	}
}

func NewConsumerManager(Nameservers []string, topicName string) (*RocketManager, error) {
	cmg := &RocketManager{}
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(Nameservers)),
		producer.WithRetry(3),
	)
	if err != nil {
		return nil, errors.Wrap(err, "NewProducer")
	}
	err = p.Start()
	if err != nil {
		return nil, errors.Wrap(err, "p.Start()")
	}
	cmg.Producer = p
	cmg.Nameservers = Nameservers
	cmg.topicName = topicName
	return cmg, nil
}

const QueueTopic = "queuetopic"

func GetQueueTopic(event *event.Event) string {
	ext := event.Extensions()
	var callbackUrl string
	if val, ok := ext[QueueTopic]; ok {
		callbackUrl, _ = val.(string)
	}
	return callbackUrl
}

func (cmg *RocketManager) SendMessage(ctx context.Context, event *event.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	topic := cmg.topicName

	selfTopic := GetQueueTopic(event)
	if selfTopic != "" {
		topic = selfTopic
	}

	message := primitive.NewMessage(topic, bytes)
	message.WithTag(event.Type())
	_, err = cmg.Producer.SendSync(ctx, message)
	if err != nil {
		errors.Wrap(err, "RocketManager.SendMessage "+event.Subject())
	}
	return nil
}
