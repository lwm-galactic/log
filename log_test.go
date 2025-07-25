package log_test

import (
	"github.com/lwm-galactic/log"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func Test_WithName(t *testing.T) {
	defer log.Flush() // used for record logger printer

	logger := log.WithName("test")
	logger.Infow("Hello world!", "foo", "bar") // structed logger
}

func Test_WithValues(t *testing.T) {
	defer log.Flush() // used for record logger printer

	logger := log.WithValues("key", "value") // used for record context
	logger.Info("Hello world!")
	logger.Info("Hello world!")
}

func Test_V(t *testing.T) {
	defer log.Flush() // used for record logger printer

	log.V(0).Infow("Hello world!", "key", "value")
	log.V(1).Infow("Hello world!", "key", "value")
}

func Test_Option(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ExitOnError)
	opt := log.NewOptions()
	opt.AddFlags(fs)

	args := []string{"--log.level=debug"}
	err := fs.Parse(args)
	assert.Nil(t, err)

	assert.Equal(t, "debug", opt.Level)
}
