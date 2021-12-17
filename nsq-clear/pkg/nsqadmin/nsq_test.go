package nsqadmin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testNsqAdmin = NewNsqAdmin("http://127.0.0.1:4171")

func TestGetTopic(t *testing.T) {
	topics, err := testNsqAdmin.GetTopic()
	assert.NoError(t, err)
	t.Log(topics)
	for _, topic := range topics {
		topicInfo, err := testNsqAdmin.GetTopicDepth(topic)
		assert.NoError(t, err)
		t.Log(topic, topicInfo.Depth)
		for _, v := range topicInfo.Channels {
			t.Log(v.ChannelName, v.Depth)
		}
	}
}

func TestEmptyQueue(t *testing.T) {
	err := testNsqAdmin.EmptyQueue("test6", "ch1")
	assert.NoError(t, err)
}
