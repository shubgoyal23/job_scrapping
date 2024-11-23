package main

import (
	"fmt"
	"log"
	"nScrapper/helpers"
	"nScrapper/types"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			helpers.LogError(fmt.Sprintf("begin handler crashed because %s", r), nil)
		}
	}()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	logs := helpers.InitLogger()
	defer logs.Close()

	if err := helpers.InitPostgresDataBase(); err != nil {
		helpers.LogError("Unable to connect to Postgres, check logs", nil)
		return
	}
	if rederr := helpers.InitRediGo(os.Getenv("REDIS"), os.Getenv("REDIS_PWD")); rederr != nil {
		helpers.LogError("Unable to connect to Prod Redis, check logs", rederr)
		return
	}
	if err := helpers.InitMongoDB(); err != nil {
		helpers.LogError("Unable to connect to MongoDB, check logs", err)
		return
	}

	if err := helpers.InitBrowser(); !err {
		helpers.LogError("Unable to connect to Browser, check logs", nil)
		return
	}
	m, err := helpers.GetManyDocMongoDB("jobsScrapeMap", bson.M{})
	if err != nil {
		helpers.LogError("Unable to get data from MongoDB, check logs", err)
		return
	}
	for _, v := range m {
		jsond, err := bson.Marshal(v)
		if err != nil {
			helpers.LogError("Unable to get data from MongoDB, check logs", err)
			return
		}
		var data types.JobDataScrapeMap
		if err := bson.Unmarshal(jsond, &data); err != nil {
			helpers.LogError("Unable to get data from MongoDB, check logs", err)
			return
		}
		helpers.ScrapeMap[data.Homepage] = data
	}
	go helpers.GetDataFromLink()
	for _, v := range helpers.ScrapeMap {
		helpers.LinkDupper(v)
	}
	go func() {
		for range time.Tick(time.Hour * 12) {
			helpers.LogError(fmt.Sprintf("running Round trip of 12 hours at time: %s", time.Now().String()), nil)
			for _, v := range helpers.ScrapeMap {
				helpers.LinkDupper(v)
				time.Sleep(time.Hour * 1)
			}
		}
	}()
	go func() {
		for range time.Tick(time.Hour * 1) {
			helpers.LogError(fmt.Sprintf("running Round trip of 1 hours at time: %s", time.Now().String()), nil)
			helpers.GetDataFromLink()
		}
	}()

	go func() {
		for range time.Tick(time.Hour * 24) {
			helpers.UpdateDataFromLink()
		}
	}()

	select {}
}
