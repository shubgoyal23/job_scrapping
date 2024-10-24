package main

import (
	"fmt"
	"log"
	"nScrapper/helpers"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

var uniqueTags = make(chan string, 100)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			helpers.LogError("naukri.com "+fmt.Sprintf("begin handler crashed because %s", r), nil)
		}
	}()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(os.Getenv("REDIS"), os.Getenv("REDIS_PWD"))
	logs := helpers.InitLogger()
	defer logs.Close()

	helpers.InitDataBase()
	if redstart := helpers.InitRediGo(os.Getenv("REDIS"), os.Getenv("REDIS_PWD")); !redstart {
		helpers.LogError("Unable to connect to Prod Redis, check logs", nil)
	}
	go func() {
		for lk := range uniqueTags {
			helpers.InsertRedisList(lk)
		}
	}()

	go naukriDataFetch()
	go PudDataToSupabase()

	select {}
}

func naukriDataFetch() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			helpers.LogError("naukri.com "+fmt.Sprintf("begin handler crashed because %s", r), nil)
		}
	}()
	for t := range time.Tick(time.Hour * 12) {
		helpers.LogError(fmt.Sprintf("running link scrapper at time: %s", t), nil)

		path, _ := launcher.LookPath()
		u := launcher.New().Bin(path).Headless(true).MustLaunch()
		browser := rod.New().ControlURL(u).MustConnect()

		page := browser.MustPage("https://www.naukri.com/").MustWaitStable()
		page.Navigate("https://www.naukri.com/it-jobs?src=gnbjobs_homepage_srch")

		count, countErr := page.Element(".styles_count-string__DlPaZ")
		if countErr != nil {
			helpers.LogError("countErr", countErr)
		}
		countTxt, countTxtErr := count.Text()
		if countTxtErr != nil {
			helpers.LogError("countTxtErr", countTxtErr)
		}
		res := strings.Split(countTxt, " of ")
		num, conErr := strconv.Atoi(res[1])
		if conErr != nil {
			helpers.LogError("conErr", conErr)
		}
		totalPage := (num / 20) + 1

		for i := 1; i <= totalPage; i++ {
			linkstr := fmt.Sprintf("https://www.naukri.com/it-jobs-%d?src=gnbjobs_homepage_srch", i)
			pageNErr := page.Navigate(linkstr)
			if pageNErr != nil {
				helpers.LogError("naukri.com", pageNErr)
			}
			page.MustWaitDOMStable()

			aTags, aTagErr := page.Elements("a")
			fmt.Println("aTags", len(aTags))
			if aTagErr != nil {
				helpers.LogError("aTagErr", aTagErr)
			}

			for _, a := range aTags {
				aTag, aTagErr := a.Attribute("href")
				if aTag == nil {
					continue
				}
				if aTagErr != nil {
					helpers.LogError("aTagErr", aTagErr)
				}
				lk := helpers.CleanUrl(*aTag, "https://www.naukri.com")
				if lk == "" {
					continue
				}
				uniqueTags <- lk
			}
		}
		browser.MustClose()
	}
}

func PudDataToSupabase() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			helpers.LogError("naukri.com "+fmt.Sprintf("begin handler crashed because %s", r), nil)
		}
	}()
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).Headless(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	page := browser.MustPage("https://www.naukri.com/").MustWaitStable()

	for {
		conn := helpers.Redigo.Get()
		if conn.Err() != nil {
			helpers.LogError("cannot get from redis", conn.Err())
			time.Sleep(time.Hour * 1)
			continue
		}
		val, err := conn.Do("SPOP", "all_job_lisks")
		if err != nil {
			helpers.LogError("cannot get from redis", err)
			time.Sleep(time.Hour * 1)
			continue
		}
		if val == nil {
			time.Sleep(time.Hour * 1)
			continue
		}

		page.Navigate(string(val.([]string)[0]))
		page.MustWaitDOMStable()
		jobD := helpers.NaukriElements(page)
		if c := helpers.InsertSupabase(jobD); c {
			helpers.InsertRedisSet(string(val.([]string)[0]))
		} else {
			helpers.InsertRedisList(string(val.([]string)[0]))
		}
	}
}
