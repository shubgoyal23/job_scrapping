package helpers

import (
	"fmt"
	"nScrapper/types"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/supabase-community/supabase-go"
)

var Conn *supabase.Client
var Redigo *redis.Pool

func InitDataBase() {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
	Conn = client
}

func Insert(val types.JobListing) {
	res, intr, err := Conn.From("job_listings").Insert(val, false, "replace", "returning", "id").Execute()
	if err != nil {
		fmt.Println("cannot insert", err)
	}
	fmt.Println(res, intr)
}

func InitRediGo(r string, pwd string) bool {
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
		Redigo = nil
		return false
	} else {
		Redigo = pool
		return true
	}
}

func InsertRedisList(val string) {
	rc := Redigo.Get()
	defer rc.Close()
	_, err := rc.Do("LPUSH", "all_job_lisks", val)
	if err != nil {
		LogError("cannot insert", err)
	}
}
func InsertRedisSet(val string) {
	rc := Redigo.Get()
	defer rc.Close()
	_, err := rc.Do("SADD", "job_links_posted", val)
	if err != nil {
		LogError("cannot insert", err)
	}
}
func CheckRedisSetMemeber(val string) bool {
	rc := Redigo.Get()
	defer rc.Close()
	f, err := rc.Do("SISMEMBER", "job_links_posted", val)
	fmt.Println(f)
	if err != nil {
		LogError("cannot insert", err)
		return false
	}
	return true
}

func InsertSupabase(val types.JobListing) bool {
	_, _, err := Conn.From("job_listings").Insert(val, false, "replace", "returning", "id").Execute()
	if err != nil {
		fmt.Println("cannot insert", err)
		return false
	}
	return true
}
