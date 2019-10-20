/**
 * @Author: wei.tan
 * @Description:  redis 客户端使用
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2019-10-14 22:52
 */

package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

func main()  {
	coon, err := redis.Dial("tcp", ":6379")
	if err != nil{
		log.Fatalln(err)
		return
	}
	defer coon.Close()


	//执行 set
	rs, err := coon.Do("SET", "xxx", "hello")
	if err != nil{
		log.Fatalln(err)
		return
	}
	log.Println("设置 key", rs)

	//执行 get
	rs, err = redis.String(coon.Do("GET", "xxx"))
	if err != nil{
		log.Fatalln(err)
		return
	}
	log.Println("获取 key", rs)


	//获取所有 key
	rs, err = redis.Strings(coon.Do("Keys", "*"))
	if err != nil{
		log.Fatalln(err)
		return
	}
	log.Println("获取所有 key", rs)
}