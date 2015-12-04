package env
import (
	"os"
	"errors"
	"strconv"
)

func GetOr(key string, fallback string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}
	return fallback
}

func GetOrInt(key string, fallback int) (int, error) {
	if strVal, found := os.LookupEnv(key); found != false {
		if val,err := strconv.Atoi(strVal); err == nil {
			return val, nil
		} else {
			return 0, errors.New("Failed to make int from ENV " + key + ": " + err.Error())
		}
	} else {
		return fallback, nil
	}
}
