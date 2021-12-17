package nsqadmin

var (
	contentType = "application/json"

	getTopicsUrl   = `%s/api/topics`
	topicsUrl      = `%s/api/topics/%s`
	emptyQueueUrl  = `%s/api/topics/%s/%s`
	emptyTopicBody = `{"action": "empty"}`
)

type Topics struct {
	Message string   `json:"message"`
	Topics  []string `json:"topics"`
}

type Channel struct {
	ChannelName   string `json:"channel_name"`
	DeferredCount int64  `json:"deferred_count"`
	Depth         uint64 `json:"depth"`
	InFlightCount int64  `json:"in_flight_count"`
	MemoryDepth   int64  `json:"memory_depth"`
	MessageCount  int64  `json:"message_count"`
	Paused        bool   `json:"paused"`
	RequeueCount  int64  `json:"requeue_count"`
	TimeoutCount  int64  `json:"timeout_count"`
}

type TopicInfo struct {
	Channels     []Channel `json:"channels"`
	Depth        uint64    `json:"depth"`
	Message      string    `json:"message"`
	MessageCount int64     `json:"message_count"`
}
