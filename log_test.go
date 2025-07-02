package log

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	l := NewLogger()
	l.SetOptions(
		WithLevel(DebugLevel),
		WithOutputPaths(FileOutputPath, StdoutOutputPath),
		WithName("test"),
		WithCompress(true),
		WithMaxSize(10),
		WithMaxBackups(5),
		WithMaxAge(7),
		WithEnableColor(true),
		WithFilePath("./log/zap.log"),
	)
	err := l.Build()
	if err != nil {
		fmt.Println(err)
		return
	}
	l.Debugf("测试debug %v", "test")
	l.Infof("测试info %v", "test")
	l.Errorf("测试error %v", "test")
}
