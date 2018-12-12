package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Debug("debug")
	Info("info")
	Warn("warn")
	Fatal("fatal")
	Error("error")

}
