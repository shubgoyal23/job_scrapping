package main

import (
	"fmt"
	"log"
	"nScrapper/helpers"

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

	helpers.InitDataBase()
	logs := helpers.InitLogger()
	defer logs.Close()

	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).Headless(false).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	// browser := rod.New().MustConnect().NoDefaultDevice()
	defer browser.MustClose()

	page := browser.MustPage("https://www.naukri.com/").MustWaitStable()
	uniqueTags := make(map[string]bool)

	for i := 1; i < 2; i++ {
		linkstr := fmt.Sprintf("https://www.naukri.com/it-jobs-%d?src=gnbjobs_homepage_srch", i)
		// pageNErr := page.Navigate("https://www.naukri.com/it-jobs?src=gnbjobs_homepage_srch")
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
			uniqueTags[lk] = true
		}
	}
	fmt.Println(len(uniqueTags))
	// for k := range uniqueTags {
	// 	helpers.LogError(k, nil)
	// }

	helpers.Insert(uniqueTags)

}
