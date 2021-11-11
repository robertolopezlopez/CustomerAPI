package tools

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateUUID4() (string, error) {
	bytes := make([]byte, 16)
	rand.Seed(time.Now().UnixNano())
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16]), nil
}
