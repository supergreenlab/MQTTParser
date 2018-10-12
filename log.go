package main

import (
	"strconv"
	"strings"
)

type RawLog struct {
	Id      string `json:"id"`
	Channel string `json:"channel"`
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

func newRawLog(topic, payload string) RawLog {
	ts := strings.Split(topic, ".")
	id := ts[0]
	channel := ts[1]
	rl := RawLog{
		id, channel, topic, payload,
	}
	return rl
}

type Log struct {
	RawLog
	Level     string `json:"level"`
	Timestamp int    `json:"timestamp"`
	Tag       string `json:"tag"`
	Module    string `json:"module"`
	Msg       string `json:"msg"`
}

func newLog(rl RawLog) Log {
	payload := string(colorTrimExpr.ReplaceAllString(rl.Payload, ""))
	sm := msgExpr.FindStringSubmatch(payload)

	level := sm[1]
	ts, _ := strconv.Atoi(sm[2])

	tag := sm[3]
	module := sm[4]
	msg := sm[5]

	l := Log{
		rl, level, ts, tag, module, msg,
	}

	return l
}

type KeyValueLog struct {
	Log
	Kvs map[string]string `json:"kvs"`
	Kvi map[string]int    `json:"kvi"`
}

func newKeyValueLog(l Log) KeyValueLog {
	kvl := KeyValueLog{
		l, map[string]string{}, map[string]int{},
	}
	vars := kvExpr.FindAllStringSubmatch(l.Msg, -1)
	for _, varMatch := range vars {
		varName := varMatch[2]
		varValue := varMatch[3]
		numValue, err := strconv.Atoi(varValue)
		if err == nil {
			kvl.Kvi[varName] = numValue
		} else {
			kvl.Kvs[varName] = varValue
		}
	}
	return kvl
}
