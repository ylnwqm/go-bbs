package utils
import (
	"github.com/garyburd/redigo/redis"
	"fmt"
)
var RedisPool *redis.Pool  //创建redis连接池

func LPush(key,value string){
	redisClient := RedisPool.Get()
	defer redisClient.Close()
	// fmt.Printf("%s,%s\n",key,value)
	_,err := redisClient.Do("LPush",key,string(value))
	if err != nil {
		fmt.Printf("队列：%s,错误：%s,信息：%s\n",key,err.Error(),value)
		// LPush(key,value)
	}
}