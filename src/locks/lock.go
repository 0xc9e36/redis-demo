/**
 * @Author: wei.tan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2019-10-13 11:48
 */

package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

/*

	方案1: 进程在 do sth 期间挂掉， redis 锁永远不会释放, 导致死锁
	setnx lock_a random_value
	// do sth
	delete lock_a


	方案二：同样的问题, 有可能没来得及设置超时时间进程就挂掉, 导致死锁, SETNX/SETEX 不是原子操作
	setnx lock_a random_value
	setex lock_a 10 random_value // 10s超时
	// do sth
	delete lock_a

	方案三: 进程 1 获取锁, 过期时间为 10s, 但是业务还在运行
           这时进程 2 获取到了锁，运行业务...
		   此时如果进行 1 执行完毕，会成功释放锁.
	SET lock_a random_value NX PX 10000 // 10s超时
	// do sth
	delete lock_a


	方案四: 类似 CAS, 比较删除, 只释放自己加的锁, 不过 redis 原生命令不支持, 需要借助 lua 脚本实现
	SET lock_a random_value NX PX 10000
	// do sth
	eval "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end" 1 lock_a random_value

	注意:
	1.超时时间: 如果业务大于锁过期时间，可以新开线程续租(更新过期时间)
	2.重试: 拿不到锁，可以轮询


	以上针对于单点实现， 集群方案可选择 redlock
	redis 实现的分布式锁缺陷很多， 需要强一致性可使用 zk, etcd 等
*/

func init() {
	Redispool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			tcp := fmt.Sprintf("%s:%d", "127.0.0.1", 6379)
			c, err := redis.Dial("tcp", tcp)
			if err != nil {
				return nil, err
			}
			fmt.Println("connect redis success!")
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

}

type RedisLock struct {
	lockKey string
	value   string
}

//保证原子性（redis是单线程），避免del删除了，其他client获得的lock
var delScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`)

func (this *RedisLock) Lock(rd *redis.Conn, timeout int) error {

	{ //随机数
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			return err
		}
		this.value = base64.StdEncoding.EncodeToString(b)
	}
	lockReply, err := (*rd).Do("SET", this.lockKey, this.value, "ex", timeout, "nx")
	if err != nil {
		return errors.New("redis fail")
	}
	if lockReply == "OK" {
		return nil
	} else {
		return errors.New("lock fail")
	}
}

func (this *RedisLock) Unlock(rd *redis.Conn) {
	delScript.Do(*rd, this.lockKey, this.value)
}



func main() {
	rd := Redispool.Get()
	defer rd.Close()

	go func() {
		Alock := RedisLock{lockKey: "xxxxx"}
		err := Alock.Lock(&rd, 5) //5 秒后自动删除Alock

		time.Sleep(7 * time.Second) //等待7秒
		fmt.Println("111", err)
		Alock.Unlock(&rd) //想删除的是Alock锁，但是Alock 已经被自动删除 ,Block由于value 不一样，所以也不会删除
	}()

	time.Sleep(6 * time.Second) //此时Alock 已经被删除
	Block := RedisLock{lockKey: "xxxxx"}
	err := Block.Lock(&rd, 5) //此时 会获取新的lock Block
	fmt.Println("222", err)

	time.Sleep(2 * time.Second)
	Clock := RedisLock{lockKey: "xxxxx"}
	err = Clock.Lock(&rd, 5) //想获取新的lock Clock，但由于 Block还存在，返回错误
	fmt.Println("333", err)

	time.Sleep(10 * time.Second)

}

var Redispool *redis.Pool

