package flog

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const DefalutLogFile = "./flog/flog.txt"

var flogger zerolog.Logger

func Flogger() *zerolog.Logger {
	return &flogger
}

func Init(path string) error {
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Error().Caller().Msgf("opening log file failed during initializing the logger: %s", err)
		return fmt.Errorf("opening log file failed during initializing the logger: %s", err)
	}

	writers := zerolog.MultiLevelWriter(fd, os.Stderr)

	flogger = zerolog.New(writers).With().Timestamp().Caller().Logger().Level(zerolog.InfoLevel)

	return nil
}

var loggerwriters zerolog.LevelWriter

func InitLoggerWiters(path string) (zerolog.LevelWriter, error) {
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Error().Caller().Msgf("opening log file failed: %s", err)
		return nil, fmt.Errorf("opening log file failed: %s", err)
	}

	return zerolog.MultiLevelWriter(fd, os.Stderr), nil
}
