package conf

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                  string
	BaseUrl               string
	GracefulShutdownSec   time.Duration
	ServerReadTimeoutSec  time.Duration
	ServerWriteTimeoutSec time.Duration
	ServerIdleTimeoutSec  time.Duration
	FilesCleanIntervalMil time.Duration
}

func InitConf() (Config, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:                  os.Getenv("PORT"),
		BaseUrl:               os.Getenv("BASE_URL"),
		GracefulShutdownSec:   time.Second * time.Duration(intEnv("GRACEFUL_SHUTDOWN_SEC")),
		ServerReadTimeoutSec:  time.Second * time.Duration(intEnv("SERVER_READ_TIMEOUT_SEC")),
		ServerWriteTimeoutSec: time.Second * time.Duration(intEnv("SERVER_WRITE_TIMEOUT_SEC")),
		ServerIdleTimeoutSec:  time.Second * time.Duration(intEnv("SERVER_IDLE_TIMEOUT_SEC")),
		FilesCleanIntervalMil: time.Millisecond * time.Duration(intEnv("FILE_CLEAN_INTERVAL_MIL")),
	}, nil
}

func intEnv(env string) int {
	tm := os.Getenv(env)
	if tm == "" {
		return 0
	}

	i, err := strconv.Atoi(tm)
	if err != nil {
		return 0
	}

	return i
}
