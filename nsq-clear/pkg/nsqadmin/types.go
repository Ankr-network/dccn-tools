package nsqadmin

var (
	contentType = "application/json"

	getTopicsUrl   = `%s/api/topics`
	topicsUrl      = `%s/api/topics/%s`
	emptyTopicBody = `{"action": "empty"}`
)

type Topics struct {
	Message string   `json:"message"`
	Topics  []string `json:"topics"`
}

type TopicInfo struct {
	BackendDepth int64  `json:"backend_depth"`
	Depth        uint64 `json:"depth"`
	Hostname     string `json:"hostname"`
	MemoryDepth  int64  `json:"memory_depth"`
	Message      string `json:"message"`
	MessageCount int64  `json:"message_count"`
	Node         string `json:"node"`
	TopicName    string `json:"topic_name"`
}
