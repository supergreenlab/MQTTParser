package mqtt

import (
	"strconv"
	"strings"

	mqttparser "github.com/SuperGreenLab/MQTTParser/pkg"
)

func newRawLog(topic, payload string) mqttparser.RawLog {
	ts := strings.Split(topic, ".")
	id := ts[0]
	channel := ts[1]
	rl := mqttparser.RawLog{
		ID:      id,
		Channel: channel,
		Topic:   topic,
		Payload: payload,
	}
	return rl
}

func newLog(rl mqttparser.RawLog) mqttparser.Log {
	payload := string(colorTrimExpr.ReplaceAllString(rl.Payload, ""))
	sm := msgExpr.FindStringSubmatch(payload)

	level := sm[1]
	ts, _ := strconv.Atoi(sm[2])

	tag := sm[3]
	module := sm[4]
	msg := sm[5]

	l := mqttparser.Log{
		RawLog:    rl,
		Level:     level,
		Timestamp: ts,
		Tag:       tag,
		Module:    module,
		Msg:       msg,
	}

	return l
}

func newKeyValueLog(l mqttparser.Log) mqttparser.KeyValueLog {
	kvl := mqttparser.KeyValueLog{
		Log: l,
		Kvs: map[string]string{},
		Kvi: map[string]float64{},
	}
	vars := kvExpr.FindAllStringSubmatch(l.Msg, -1)
	for _, varMatch := range vars {
		varName := varMatch[2]
		varValue := varMatch[3]
		numValue, err := strconv.ParseFloat(varValue, 64)
		if err == nil {
			kvl.Kvi[varName] = numValue
		} else {
			kvl.Kvs[varName] = varValue
		}
	}
	return kvl
}
