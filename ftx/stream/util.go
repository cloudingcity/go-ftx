package stream

import (
	"time"
)

var unixTime = func() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
