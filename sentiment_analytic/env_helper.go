package sentiment_analytic

import (
	"os"
"strconv"
)

func getEnvInt(key string) int{
	return func() (int){
		res, err := strconv.Atoi(os.Getenv(key))
		if err == nil{
			return res
		}else{
			return -1
		}
	}()
}

func getEnvInt64(key string) int64{
	return func() (int64){
		res := getEnvInt(key)
		return int64(res)
	}()
}

func getEnvBool(key string) bool{
	return func() bool{
		res, err := strconv.ParseBool(os.Getenv(key))
		if err == nil{
			return res
		}else{
			return false
		}
	}()
}
