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
		depth, err := testNsqAdmin.GetTopicDepth(topic)
		assert.NoError(t, err)
		t.Log(depth)
		if depth > 100 {
			err = testNsqAdmin.EmptyQueue(topic)
			assert.NoError(t, err)
		}
	}
}

func TestEmptyQueue(t *testing.T) {
	err := testNsqAdmin.EmptyQueue("test")
	assert.NoError(t, err)
}
