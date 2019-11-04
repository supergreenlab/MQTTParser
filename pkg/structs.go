package mqttparser

// RawLog logs as they come out of MQTT
type RawLog struct {
	ID      string `json:"id"`
	Channel string `json:"channel"`
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

// Log text only log
type Log struct {
	RawLog
	Level     string `json:"level,omitempty"`
	Timestamp int    `json:"log_timestamp,omitempty"`
	Tag       string `json:"tag,omitempty"`
	Module    string `json:"module,omitempty"`
	Msg       string `json:"msg,omitempty"`
}

// KeyValueLog log with key value pairs
type KeyValueLog struct {
	Log
	Kvs map[string]string  `json:"kvs,omitempty"`
	Kvi map[string]float64 `json:"kvi,omitempty"`
}
