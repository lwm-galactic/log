package log

import (
	"go.uber.org/zap"
	"sync"
)

type logger struct {
	opt       *options
	logger    *zap.Logger
	mu        sync.Mutex
	entryPool *sync.Pool
}
