package main

import (
	"fmt"
	"log"
	"nScrapper/helpers"
	"os"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

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
	logs := helpers.InitLogger()
	defer logs.Close()

	helpers.InitDataBase()
	if redstart := helpers.InitRediGo(os.Getenv("REDIS"), os.Getenv("REDIS_PWD")); !redstart {
		helpers.LogError("Unable to connect to Prod Redis, check logs", nil)
	}

	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).Headless(false).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	defer browser.MustClose()

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
	fmt.Println((num))
	// uniqueTags := make(map[string]bool)
	uniqueTags := make(chan string, 100)

	go func() {
		for lk := range uniqueTags {
			helpers.InsertRedis(lk)
		}
	}()

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
	fmt.Println(len(uniqueTags))

	close(uniqueTags)
	helpers.Redigo.Close()
}
