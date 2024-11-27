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
				LogError("InitRedigo", "connection redis", err)
				return nil, err
			}
			if _, err := conn.Do("AUTH", pwd); err != nil {
				conn.Close()
				LogError("InitRedigo", "connection redis", err)
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
func InsertRedisListLPush(key string, val []string) error {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("InsertRedisListLPush", "Redis not connected", er)
		return er
	}
	ar := redis.Args{}.Add(key).AddFlat(val)
	_, err := rc.Do("LPUSH", ar...)
	if err != nil {
		LogError("InsertRedisListLPush", fmt.Sprintf("cannot insert in redis list key: %s with value: %s", key, val), err)
		return err
	}
	return nil
}

// insert data in redis list
func GetRedisListRPOP(key string, n int) ([][]byte, error) {
	r := [][]byte{}
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("GetRedisListRPOP", "Redis not connected", er)
		return r, er
	}
	res, err := redis.ByteSlices(rc.Do("RPOP", key, n))
	if err != nil {
		LogError("GetRedisListRPOP", fmt.Sprintf("Cannot get items from redis list key: %s", key), err)
		return r, err
	}
	r = append(r, res...)
	return r, nil
}

// insert data in redis set
func InsertRedisSet(key string, val string) (bool, error) {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("InsertRedisSet", "Redis not connected", er)
		return false, er
	}
	_, err := rc.Do("SADD", key, val)
	if err != nil {
		LogError("InsertRedisSet", fmt.Sprintf("cannot insert in redis set key: %s with value: %s", key, val), err)
		return false, err
	}
	return true, nil
}

// insert data in redis set
func InsertRedisSetBulk(key string, val []string) (bool, error) {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("InsertRedisSetBulk", "Redis not connected", er)
		return false, er
	}
	ar := redis.Args{}.Add(key).AddFlat(val)
	_, err := rc.Do("SADD", ar...)
	if err != nil {
		LogError("InsertRedisSetBulk", fmt.Sprintf("cannot insert in redis set key: %s with value: %s", key, val), err)
		return false, err
	}
	return true, nil
}

// check if data exists in redis set
func CheckRedisSetMemeber(key string, val string) (bool, error) {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("CheckRedisSetMemeber", "Redis not connected", er)
		return false, er
	}
	f, err := rc.Do("SISMEMBER", key, val)
	if err != nil {
		LogError("CheckRedisSetMemeber", fmt.Sprintf("cannot check in redis set key: %s with value: %s", key, val), err)
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
		LogError("DeleteRedisSetMemeber", "Redis not connected", er)
		return false, er
	}
	_, err := rc.Do("SREM", key, val)
	if err != nil {
		LogError("DeleteRedisSetMemeber", fmt.Sprintf("cannot delete in redis set key: %s with value: %s", key, val), err)
		return false, err
	}
	return true, nil
}

func GetRedisKeyVal(key string) (string, error) {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("GetRedisKeyVal", "Redis not connected", er)
		return "", er
	}
	res, err := redis.String(rc.Do("GET", key))
	if err != nil {
		LogError("GetRedisKeyVal", fmt.Sprintf("cannot get in redis key: %s", key), err)
		return "", err
	}
	return res, nil
}

func SetRedisKeyVal(key string, val string) error {
	rc := RedigoConn.Get()
	defer rc.Close()
	if _, er := rc.Do("PING"); er != nil {
		LogError("SetRedisKeyVal", "Redis not connected", er)
		return er
	}
	_, err := rc.Do("SET", key, val)
	if err != nil {
		LogError("SetRedisKeyVal", fmt.Sprintf("cannot set in redis key: %s with value: %s", key, val), err)
		return err
	}
	return nil
}
