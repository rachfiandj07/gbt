package ping

import (
	"encoding/json"
	"github.com/bitly/go-nsq"
	"github.com/pkg/errors"
	"log"
)

type (
	//PingHandler contains topic and channel
	PingHandler struct {
		Topic   string `cache:"topic" json:"topic"`
		Channel string `cache:"channel" json:"channel"`
	}

	//Data holds unmarshalled msg
	Data struct {
		Ping string `json:"ping"`
	}
)

//HandleMessage to handle ping message
func (p PingHandler) HandleMessage(message *nsq.Message) error {
	data := Data{}

	//unmarshall message
	err := json.Unmarshal(message.Body, &data)
	if err != nil {
		err = errors.Wrapf(err, "[Consumer] error when unmarshalling: %v, err: ", message.Body)
		log.Println(err.Error())
		message.Finish()
		return err
	}

	log.Println("[Consumer] ", data)

	message.Finish()
	return nil
}
