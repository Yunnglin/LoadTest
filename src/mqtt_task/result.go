package mqtt_task

import "time"

type Result struct {
	Err           error
	StatusCode    int
	Duration      time.Duration
	ConnDuration  time.Duration // connection setup(DNS lookup + Dial up) duration
	DnsDuration   time.Duration // dns lookup duration
	ReqDuration   time.Duration // request "write" duration
	ResDuration   time.Duration // response "read" duration
	DelayDuration time.Duration // delay between response and request
	ContentLength int64
	Body          []byte
}
