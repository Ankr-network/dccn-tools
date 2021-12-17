package nsqadmin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NsqAdmin interface {
	GetTopic() ([]string, error)
	GetTopicDepth(topic string) (*TopicInfo, error)
	EmptyQueue(topic string) error
}

type nsqAdmin struct {
	addr   string
	client *http.Client
}

func NewNsqAdmin(addr string) *nsqAdmin {
	return &nsqAdmin{
		addr: addr,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        50,
				MaxIdleConnsPerHost: 50,
				MaxConnsPerHost:     100,
			},
		},
	}
}

func (n *nsqAdmin) GetTopic() ([]string, error) {
	resp, err := n.client.Get(fmt.Sprintf(getTopicsUrl, n.addr))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var tp Topics
	err = json.Unmarshal(body, &tp)
	if err != nil {
		return nil, err
	}
	return tp.Topics, nil
}

func (n *nsqAdmin) GetTopicDepth(topic string) (*TopicInfo, error) {
	resp, err := n.client.Get(fmt.Sprintf(topicsUrl, n.addr, topic))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tp TopicInfo
	err = json.Unmarshal(body, &tp)
	if err != nil {
		return nil, err
	}
	return &tp, nil
}

func (n *nsqAdmin) EmptyQueue(topic string, channel string) error {
	var url string
	if channel == "" {
		url = fmt.Sprintf(topicsUrl, n.addr, topic)
	} else {
		url = fmt.Sprintf(emptyQueueUrl, n.addr, topic, channel)
	}
	resp, err := n.client.Post(url, contentType, bytes.NewBufferString(emptyTopicBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	return errors.New(fmt.Sprintf("failed to empty quest status :%d", resp.StatusCode))
}
