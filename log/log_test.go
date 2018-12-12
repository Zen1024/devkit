package log

import (
	"testing"
	"time"
)

func TestLogLevel(t *testing.T) {
	Debug("debug")
	Info("info")
	Warn("warn")
	Fatal("fatal")
	Error("error")
}

func TestSplitByMinute(t *testing.T) {
	SetFile("test.log")
	SetFlags(Lsplitm)
	for i := 0; i < 70; i++ {
		Debug("debug")
		time.Sleep(time.Second)
	}
}

func TestSplitBySize(t *testing.T) {
	SetFile("test.log")
	SetMaxSize(1024)
	for i := 0; i < 100; i++ {
		Debug("debug")
		time.Sleep(time.Millisecond * 10)
	}
}

func TestSplitBySizeMinute(t *testing.T) {
	SetFile("test.log")
	SetMaxSize(1024)
	SetFlags(Lsplitm)
	for i := 0; i < 180; i++ {
		Debug("debugdebugdebugdebugdebugdebugdebugdebugdebugdebug")
		time.Sleep(time.Second)
	}
}
