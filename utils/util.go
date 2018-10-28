package utils

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

func GetCurrentUnixTime() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func GetExpirationOfActivationKey(expirationOfActivationKey int64) string {
	return strconv.Itoa(int(time.Now().Unix() + expirationOfActivationKey))
}

func GenerateActivationKey(email string, password string, activationSalt string) string {
	h := sha256.New()
	unixtime := GetCurrentUnixTime()
	h.Write([]byte(email + password + activationSalt + unixtime))
	return fmt.Sprintf("%x", h.Sum(nil))
}
