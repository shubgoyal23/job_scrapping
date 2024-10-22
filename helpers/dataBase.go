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

func InsertRedis(val string) {
	_, err := Redigo.Get().Do("LPUSH", "job_listings", val)
	if err != nil {
		LogError("cannot insert", err)
	}
}
