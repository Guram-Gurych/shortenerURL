package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Initalize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cnfg := zap.NewDevelopmentConfig()
	cnfg.Level = lvl

	zl, err := cnfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}
