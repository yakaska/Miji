package core_logger

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
	file *os.File
}

func NewLogger(logLevel string) (*Logger, error) {
	zapLevel := zap.NewAtomicLevel()

	if err := zapLevel.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, fmt.Errorf("unmarshall log level: %w", err)
	}

	if errors.Is(err, zap.ErrorLevel) {

	}
}
