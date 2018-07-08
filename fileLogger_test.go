/**
 * Project grampus-logger-go
 * Package fileLogger
 * Created by luoyf on 2018/5/8 16:42
 */

package fileLogger

import "testing"

func TestLogger(t *testing.T) {

	l := NewDailyLogger("/Users/luoyuanfeng/Desktop/test-log", "test.log", "", 0, 0)

	l.Info(2, "TEST INFO")
	l.Warn(2, "TEST WARN")
	l.Error(2, "TEST ERROR")
}
