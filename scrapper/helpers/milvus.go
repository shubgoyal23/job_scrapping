package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"nScrapper/types"
	"net/http"
	"os"
	"strconv"
)

func PushToMilvus() {
	defer func() {
		if r := recover(); r != nil {
			LogError("PushToMilvus", fmt.Sprintf("PushToMilvus handler crashed because %s", r), nil)
		}
	}()
	url := os.Getenv("PY_MILVUS_URL")
	if url == "" {
		url = "localhost:8000"
	}
	startPage, err := GetRedisKeyVal("milvus_startpage")
	if err != nil {
		LogError("PushToMilvus", "Unable to get data from Redis key milvus_starttime, check logs", err)
		return
	}
	startP, err := strconv.Atoi(startPage)
	if err != nil {
		LogError("PushToMilvus", "Unable to convert milvus_starttime to int, check logs", err)
		return
	}
	q := "SELECT id, job_description FROM job_listings ORDER BY created_at LIMIT $1 OFFSET $2;"
	page := startP
	for {
		limit := 10
		offset := (page - 1) * limit
		r, err := GetManyDocPostgres(q, []interface{}{limit, offset})
		if err != nil {
			LogError("PushToMilvus", "Unable to get data from Postgres, check logs", err)
			return
		}
		var ids []int
		var vals []string
		for r.Next() {
			var res types.JobListing
			if err := r.Scan(&res.ID, &res.JobDescription); err != nil {
				LogError("PushToMilvus", "cannot decode doc in postgres", err)
				continue
			}
			vals = append(vals, res.JobDescription)
			ids = append(ids, res.ID)
		}
		if len(ids) == 0 {
			break
		}
		data := map[string]interface{}{
			"id":          ids,
			"description": vals,
		}
		j, _ := json.Marshal(data)
		resp, err := http.Post(url+"/vector", "application/json", bytes.NewBuffer(j))
		if err != nil {
			LogError("PushToMilvus", "Unable to push data to Milvus, check logs", err)
			break
		}
		resp.Body.Close()
		r.Close()

		page++
		if err := SetRedisKeyVal("milvus_startpage", strconv.Itoa(page)); err != nil {
			LogError("PushToMilvus", "Unable to set data in Redis key milvus_starttime, check logs", err)
			return
		}
	}
}
