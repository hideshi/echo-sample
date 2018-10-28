package utils

import (
	"strconv"
	"time"

	"github.com/hideshi/echo-sample/structs"
)

func GetCurrentUnixTime() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func GetExpirationOfActivationKey() string {
	return strconv.Itoa(int(time.Now().Unix() + structs.Conf.Auth.ExpirationOfActivationKey))
}
