package helpers

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var RedigoConn *redis.Pool

// init redis
func InitRediGo(r string, pwd string) error {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", r)
			if err != nil {
				//log to local as could not connect to Redis
				LogError("connection redis", err)
				return nil, err
			}
			if _, err := conn.Do("AUTH", pwd); err != nil {
				conn.Close()
				LogError("connection redis", err)
				return nil, err
			}
			return conn, nil
		},
	}
	if pool.Get().Err() != nil {
		RedigoConn = nil
		return pool.Get().Err()
	} else {
		RedigoConn = pool
		return nil
	}
}

// insert data in redis list
func InsertRedisList(key string, val string) (bool, error) {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("Redis not connected", er)
		return false, er
	}
	_, err := rc.Do("LPUSH", key, val)
	if err != nil {
		LogError(fmt.Sprintf("cannot insert in redis list key: %s with value: %s", key, val), err)
		return false, err
	}
	return true, nil
}

// insert data in redis set
func InsertRedisSet(key string, val string) (bool, error) {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("Redis not connected", er)
		return false, er
	}
	_, err := rc.Do("SADD", key, val)
	if err != nil {
		LogError(fmt.Sprintf("cannot insert in redis set key: %s with value: %s", key, val), err)
		return false, err
	}
	return true, nil
}

// insert data in redis set
func InsertRedisSetBulk(key string, val []string) (bool, error) {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("Redis not connected", er)
		return false, er
	}
	ar := redis.Args{}.Add(key).AddFlat(val)
	_, err := rc.Do("SADD", ar...)
	if err != nil {
		LogError(fmt.Sprintf("cannot insert in redis set key: %s with value: %s", key, val), err)
		return false, err
	}
	return true, nil
}

// check if data exists in redis set
func CheckRedisSetMemeber(key string, val string) (bool, error) {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("Redis not connected", er)
		return false, er
	}
	f, err := rc.Do("SISMEMBER", key, val)
	if err != nil {
		LogError(fmt.Sprintf("cannot check in redis set key: %s with value: %s", key, val), err)
		return false, err
	}
	if f.(int64) < 1 {
		return false, nil
	}
	return true, nil
}

// delete redis set member
func DeleteRedisSetMemeber(key string, val string) (bool, error) {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("Redis not connected", er)
		return false, er
	}
	_, err := rc.Do("SREM", key, val)
	if err != nil {
		LogError(fmt.Sprintf("cannot delete in redis set key: %s with value: %s", key, val), err)
		return false, err
	}
	return true, nil
}
