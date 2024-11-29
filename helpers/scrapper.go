package helpers

import (
	"context"
	"fmt"
	"math/rand"
	"nScrapper/types"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

var ScrapeMap = make(map[string]types.JobDataScrapeMap)
var Browser *rod.Browser

// get browser
func InitBrowser() bool {
	var PORT = ""
	h := os.Getenv("BROWSER_POST")
	if h == "" {
		PORT = "7317"
	}
	u := launcher.MustResolveURL("http://rodbrower:" + PORT)
	// u := launcher.New().Bin(path).Headless(Headless).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	browser.MustPage()
	Browser = browser
	return true
}

// this function collects all the data from the page
func ScrapperElements(page *rod.Page, jobDMap types.JobDataScrapeMap) types.JobListing {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			LogError("ScrapperElements", fmt.Sprintf("ScrapperElements: handler crashed for %s because %s", jobDMap.Homepage, r), nil)
		}
	}()

	var jobDetails types.JobListing

	valuess := reflect.ValueOf(jobDMap.JobData)
	typeofv := reflect.TypeOf(jobDMap.JobData)

	for i := 0; i < valuess.NumField(); i++ {
		fname := typeofv.Field(i).Name
		fVal := valuess.Field(i).Interface().(types.TagField)

		if fVal.TagType == "string" {
			if element, err := page.Timeout(5 * time.Second).Element(fVal.Element); err == nil {
				if text, err := element.Text(); err == nil {
					text = CleanText(text, fVal.Cleaner)
					SetField(&jobDetails, fname, text)
				} else {
					LogError("ScrapperElements", "cannot get "+fname, err)
				}
			}
		} else if fVal.TagType == "numeric" {
			if element, err := page.Timeout(5 * time.Second).Element(fVal.Element); err == nil {
				if text, err := element.Text(); err == nil {
					text = CleanText(text, fVal.Cleaner)
					SetField(&jobDetails, fname, text)
				} else {
					LogError("ScrapperElements", "cannot get "+fname, err)
				}
			}
		} else if fVal.TagType == "url" {
			if element, err := page.Timeout(5 * time.Second).Element(fVal.Element); err == nil {
				if text, err := element.Attribute(fVal.AttributeTarget); err == nil {
					if text == nil {
						continue
					}
					text1 := CleanUrl(*text, jobDMap.Homepage)
					SetField(&jobDetails, fname, text1)
				} else {
					LogError("ScrapperElements", "cannot get "+fname, err)
				}
			}
		} else if fVal.TagType == "date" {
			if element, err := page.Timeout(5 * time.Second).Element(fVal.Element); err == nil {
				if text, err := element.Text(); err == nil {
					var date time.Time = time.Now()
					text = CleanText(text, fVal.Cleaner)
					d := strings.Split(text, " ")
					if agoTago, e := strconv.Atoi(d[0]); e == nil {
						if len(d) > 1 {
							if d[1] == "days" || d[1] == "day" {
								date = time.Now().AddDate(0, 0, agoTago*-1)
							} else if d[1] == "months" || d[1] == "month" {
								date = time.Now().AddDate(0, agoTago*-1, 0)
							} else if d[1] == "years" || d[1] == "year" {
								date = time.Now().AddDate(agoTago*-1, 0, 0)
							}
						}
					}
					SetField(&jobDetails, fname, date)
				} else {
					LogError("ScrapperElements", "cannot get "+fname, err)
				}
			}
		} else if fVal.TagType == "[]string" {
			if element, err := page.Timeout(5 * time.Second).Elements(fVal.Element); err == nil {
				var elms []string
				for _, elem := range element {
					if text, err := elem.Text(); err == nil {
						text = CleanText(text, fVal.Cleaner)
						elms = append(elms, text)
					} else {
						LogError("ScrapperElements", "cannot get "+fname, err)
					}
				}
				SetField(&jobDetails, fname, elms)
			}
		} else if fVal.TagType == "range" {
			if element, err := page.Timeout(5 * time.Second).Element(fVal.Element); err == nil {
				if text, err := element.Text(); err == nil {
					text = CleanText(text, fVal.Cleaner)
					text = strings.Trim(text, ".")
					d := strings.Split(text, "-")
					min, max, setdata := 0.0, 0.0, 0.0
					if len(d) < 2 {
						SetField(&jobDetails, fname, setdata)
						continue
					}
					if m, err := strconv.ParseFloat(d[0], 64); err == nil {
						min = m
					}
					if m, err := strconv.ParseFloat(d[1], 64); err == nil {
						max = m
					}
					if strings.Contains(fname, "Min") {
						setdata = min
					} else if strings.Contains(fname, "Max") {
						setdata = max
					}
					SetField(&jobDetails, fname, setdata)
				} else {
					LogError("ScrapperElements", "cannot get "+fname, err)
				}
			}
		} else if fVal.TagType == "bool" {
			if element, err := page.Timeout(5 * time.Second).Element(fVal.Element); err == nil {
				var r bool = false
				if text, err := element.Text(); err == nil {
					text = CleanText(text, fVal.Cleaner)
					if text != "" {
						r = true
					}
					SetField(&jobDetails, fname, r)
				} else {
					LogError("ScrapperElements", "cannot get "+fname, err)
				}
			}
		}
	}

	return jobDetails
}

// this function finds all the links of jobs on the page and stores them in chan UniqueLinks
func LinkDupper(jobMap types.JobDataScrapeMap) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			LogError("LinkDupper", fmt.Sprintf("LinkDupper handler crashed because %s", r), nil)
		}
	}()
	LogError("LinkDupper", fmt.Sprintf("running linkDuper for %s scrapper at time: %s", jobMap.Homepage, time.Now().String()), nil)

	// page := Browser.MustPage(jobMap.Homepage).MustWaitStable()
	page := stealth.MustPage(Browser)
	defer page.Close()

	AllTags := make(map[string]bool)
	links := []string{}
	for _, pl := range jobMap.PageLinks {
		pageNErr := page.Navigate(pl.Link)
		if pageNErr != nil {
			LogError("LinkDupper", fmt.Sprintf("error while navigating to page: %s", pl.Link), pageNErr)
			continue
		}
		errcount := 0
		for {
			if err := page.Timeout(30*time.Second).WaitDOMStable(10*time.Second, 5); err != nil {
				LogError("LinkDupper", fmt.Sprintf("error while waiting for page: %s, element: %s, errorcount : %d, err:", pl.Link, pl.Element, errcount), err)
				errcount++
			}
			aTags, aTagErr := page.Timeout(30 * time.Second).Elements(pl.Element)
			if aTagErr != nil {
				LogError("LinkDupper", fmt.Sprintf("error while getting elements from page: %s, element: %s, err:", pl.Link, pl.Element), aTagErr)
			}
			for _, a := range aTags {
				aTag, aTagErr := a.Attribute("href")
				if aTag == nil {
					continue
				}
				if aTagErr != nil {
					LogError("LinkDupper", fmt.Sprintf("error while getting href from page: %s, element: %s, err:", pl.Link, pl.Element), aTagErr)
				}
				lk := CleanUrl(*aTag, jobMap.Homepage)
				if lk == "" {
					continue
				}
				if _, ok := AllTags[lk]; ok {
					continue
				}
				AllTags[lk] = true
				links = append(links, lk)
			}
			// insert in redis
			if len(links) >= 100 {
				if err := InsertRedisListLPush("job_links", links); err != nil {
					LogError("LinkDupper", "cannot insert in redis list", err)
				} else {
					links = []string{}
				}
			}
			// next button click
			nextBtn, nextBtnErr := page.Timeout(30 * time.Second).Element(pl.NextPageBtn)
			if nextBtnErr != nil {
				LogError("LinkDupper", fmt.Sprintf("error while getting next page button from page: %s, element: %s, errorcount : %d, err:", pl.Link, pl.NextPageBtn, errcount), nextBtnErr)
				errcount++
				// frame := page.MustElements("iframe")
				// for _, f := range frame {
				// 	f.MustFrame()
				// 	checkbox, err := f.Element(`#PNSoX8 > div > label > input[type=checkbox]`)
				// 	if err != nil {
				// 		LogError(fmt.Sprintf("error while getting next page button from page: %s, element: %s, err:", pl.Link, pl.NextPageBtn), err)
				// 		continue
				// 	}
				// 	if err := checkbox.Click("left", 1); err != nil {
				// 		LogError(fmt.Sprintf("error while clicking next page button from page: %s, element: %s, err:", pl.Link, pl.NextPageBtn), err)
				// 		continue
				// 	}
				// }

			}
			nextBtnErr = nextBtn.Click("left", 1)
			if nextBtnErr != nil {
				LogError("LinkDupper", fmt.Sprintf("error while clicking next page button from page: %s, element: %s, errorcount : %d, err:", pl.Link, pl.NextPageBtn, errcount), nextBtnErr)
				errcount++
			}
			RandTimeSleep(3)
			if errcount >= 10 {
				break
			}
		}
	}
}

// this function gets the data from the links in chan UniqueTags and stores it in postgres and redis set
func GetDataFromLink() {
	defer func() {
		if r := recover(); r != nil {
			LogError("GetDataFromLink", fmt.Sprintf("GetDataFromLink crashed because %s", r), nil)
		}
	}()
	LogError("GetDataFromLink", fmt.Sprintf("running GetDataFromLink at time: %s", time.Now().String()), nil)
	// page := Browser.MustPage()
	page := stealth.MustPage(Browser)
	defer page.Close()

	var JobD struct {
		mu       sync.Mutex
		JobData  []types.JobListing
		JobLinks []string
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Redislinks := make(chan string, 100)
	var errorCount int
	// insert data if crashed
	go func(ctx context.Context) {
		<-ctx.Done()
		LogError("GetDataFromLink", "GetDataFromLink context done", nil)
		close(Redislinks)
		JobD.mu.Lock()
		if len(JobD.JobData) != 0 {
			if fres, err := InsertBulkDataPostgres(JobD.JobData); err != nil {
				LogError("GetDataFromLink", "cannot insert job data", err)
			} else {
				if len(fres) > 0 {
					InsertRedisSetBulk("Failde_posted_job_links", fres)
				}
				InsertRedisSetBulk("posted_job_links", JobD.JobLinks)
			}
			JobD.JobData = []types.JobListing{}
			JobD.JobLinks = []string{}
		}
		JobD.mu.Unlock()
		pending := []string{}
		for d := range Redislinks {
			pending = append(pending, d)
		}
		if len(pending) == 0 {
			return
		}
		if e := InsertRedisListLPush("job_links", pending); e != nil {
			LogError("GetDataFromLink", "cannot insert in redis list", e)
		}
	}(ctx)

	for {
		UniqueTags, err := GetRedisListRPOP("job_links", 100)
		if err != nil {
			LogError("GetDataFromLink", "cannot get from redis", err)
			return
		}
		for _, d := range UniqueTags {
			Redislinks <- string(d)
		}
		for d := range Redislinks {
			if errorCount >= 10 {
				return
			}
			link := d
			// check if link exists in redis
			if f, err := CheckRedisSetMemeber("posted_job_links", link); err != nil {
				LogError("GetDataFromLink", "cannot get from redis", err)
				errorCount++
				continue
			} else if f {
				continue
			}
			pageNErr := page.Timeout(30 * time.Second).Navigate(link)
			if pageNErr != nil {
				LogError("GetDataFromLink", fmt.Sprintf("error while navigating to page: %s", link), pageNErr)
				errorCount++
				continue
			}
			if err := page.Timeout(30 * time.Second).WaitLoad(); err != nil {
				LogError("GetDataFromLink", "error while waiting for page to be stable", err)
				errorCount++
				continue
			}
			u, err := url.Parse(link)
			if err != nil {
				LogError("GetDataFromLink", "cannot parse url", err)
				errorCount++
				continue
			}
			homeDomain := u.Scheme + "://" + u.Host
			sMap, ok := ScrapeMap[homeDomain]
			if !ok {
				LogError("GetDataFromLink", fmt.Sprintf("cannot get ScrapeMap for %s", homeDomain), nil)
				errorCount++
				continue
			}

			// scrape elements on page
			da := ScrapperElements(page, sMap)
			da.JobURL = link
			da.ApplicationDeadline = da.JobPostingDate.AddDate(0, 0, 30)
			da.CreatedAt = time.Now()
			da.UpdatedAt = time.Now()
			JobD.mu.Lock()
			JobD.JobData = append(JobD.JobData, da)
			JobD.JobLinks = append(JobD.JobLinks, link)
			JobD.mu.Unlock()
			if len(JobD.JobData) == 100 {
				JobD.mu.Lock()
				if fres, err := InsertBulkDataPostgres(JobD.JobData); err != nil {
					LogError("GetDataFromLink", "cannot insert job data", err)
				} else {
					if len(fres) > 0 {
						InsertRedisSetBulk("Failde_posted_job_links", fres)
					}
					InsertRedisSetBulk("posted_job_links", JobD.JobLinks)
				}
				JobD.JobData = []types.JobListing{}
				JobD.JobLinks = []string{}
				JobD.mu.Unlock()
			}
			RandTimeSleep(3)
		}
	}
}

// this function gets the data from the links in chan UniqueTags and stores it in postgres and redis set
func UpdateDataFromLink() {
	defer func() {
		if r := recover(); r != nil {
			LogError("UpdateDataFromLink", fmt.Sprintf("updateDataFromLink crashed because %s", r), nil)
		}
	}()

	// page := Browser.MustPage()
	page := stealth.MustPage(Browser)
	defer page.Close()

	r, err := GetManyDocPostgres("SELECT * FROM job_listings WHERE updated_at > now() - interval '7 days'", nil)
	if err != nil {
		LogError("UpdateDataFromLink", "Unable to get data from Postgres, check logs", err)
		return
	}
	defer r.Close()
	var val []types.JobListing
	for r.Next() {
		var res types.JobListing
		if err := r.Scan(&res.ID, &res.JobTitle, &res.CompanyName, &res.CompanyURL, &res.JobDescription, &res.JobType, &res.Location, &res.RemoteOption, &res.SalaryMin, &res.SalaryMax, &res.ExperienceMin, &res.ExperienceMax, &res.EducationRequirements, &res.Skills, &res.Benefits, &res.JobPostingDate, &res.ApplicationDeadline, &res.JobURL, &res.CreatedAt, &res.UpdatedAt, &res.IsActive); err != nil {
			LogError("UpdateDataFromLink", "cannot decode doc in postgres", err)
			continue
		}
		val = append(val, res)
	}
	del := func(val int, str string) {
		if err := DeleteDocPostgres("DELETE FROM job_listings WHERE id = $1", val); err != nil {
			LogError("UpdateDataFromLink", "cannot delete doc in postgres", err)
		}
		if _, err := DeleteRedisSetMemeber("posted_job_links", str); err != nil {
			LogError("UpdateDataFromLink", "cannot delete from redis", err)
		}
	}

	for _, jdata := range val {
		link := jdata.JobURL
		pageNErr := page.Timeout(30 * time.Second).Navigate(link)
		if pageNErr != nil {
			LogError("UpdateDataFromLink", fmt.Sprintf("error while navigating to page: %s", link), pageNErr)
			del(jdata.ID, link)
			continue
		}
		if err := page.Timeout(30 * time.Second).WaitLoad(); err != nil {
			LogError("UpdateDataFromLink", "error while waiting for page to be stable", err)
			del(jdata.ID, link)
			continue
		}
		u, err := url.Parse(link)
		if err != nil {
			LogError("UpdateDataFromLink", "cannot parse url", err)
			continue
		}
		homeDomain := u.Scheme + "://" + u.Host
		sMap, ok := ScrapeMap[homeDomain]
		if !ok {
			LogError("UpdateDataFromLink", fmt.Sprintf("cannot get ScrapeMap for %s", homeDomain), nil)
			continue
		}

		// scrape elements on page
		da := ScrapperElements(page, sMap)
		if reflect.DeepEqual(da, jdata) {
			if e := UpdateDocPostgres("UPDATE job_listings SET updated_at = now() WHERE id = $1", jdata.ID); e != nil {
				LogError("UpdateDataFromLink", "cannot update doc in postgres", e)
			}
		} else {
			da.ID = jdata.ID
			da.CreatedAt = jdata.CreatedAt
			da.UpdatedAt = time.Now()
			da.ApplicationDeadline = jdata.ApplicationDeadline.AddDate(0, 0, 10)
			UpdateDocPostgres("UPDATE job_listings SET (job_title, company_name, company_url, job_description, job_type, location, remote_option, salary_min, salary_max, experience_min, experience_max, education_requirements, skills, benefits, job_posting_date, application_deadline, job_url, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)  WHERE id = $20", []interface{}{da.JobTitle, da.CompanyName, da.CompanyURL, da.JobDescription, da.JobType, da.Location, da.RemoteOption, da.SalaryMin, da.SalaryMax, da.ExperienceMin, da.ExperienceMax, da.EducationRequirements, da.Skills, da.Benefits, da.JobPostingDate, da.ApplicationDeadline, da.JobURL, da.CreatedAt, da.UpdatedAt, da.ID})
		}

	}
}

func SetField(obj interface{}, fieldName string, value interface{}) error {
	// Get the pointer to the struct
	structValue := reflect.ValueOf(obj).Elem()
	structField := structValue.FieldByName(fieldName)

	if !structField.IsValid() {
		return fmt.Errorf("no such field: %s in struct", fieldName)
	}

	if !structField.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldName)
	}

	// Get the value to set and assign it to the field
	val := reflect.ValueOf(value)

	if structField.Type() != val.Type() {
		return fmt.Errorf("provided value type doesn't match field type")
	}

	structField.Set(val)
	return nil
}

func CleanText(text string, cleaner string) string {
	regexp := regexp.MustCompile(cleaner)
	text = regexp.ReplaceAllString(text, "")
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, "\n", ". ")
	return text
}

func CleanUrl(l string, home_url string) string {
	re := regexp.MustCompile(`(?P<Protocol>https?)://(?P<Domain>[^/]+)`)

	str, err := url.Parse(l)
	if err != nil {
		return "" // Invalid URL; return empty string
	}

	// If the URL does not have a hostname, join with the home_url
	if len(str.Hostname()) == 0 {
		ourl, oerr := url.JoinPath(home_url, str.Path)
		if oerr != nil {
			// Log the error and return empty string
			LogError("CleanUrl", "ourl", oerr)
			return ""
		}
		// Check if the cleaned URL matches the regex
		if !re.MatchString(ourl) {
			return ""
		}
		if str.RawQuery != "" {
			ourl = ourl + "?" + str.RawQuery
		}
		return ourl
	}

	// Ensure the URL uses HTTPS or HTTP
	if str.Scheme != "https" && str.Scheme != "http" {
		str.Scheme = "https"
	}

	// Parse the home URL to extract the hostname
	// h, err := url.Parse(home_url)
	// if err != nil || str.Hostname() != h.Hostname() {
	// 	return "" // Return empty if hostnames do not match or if home_url is invalid
	// }

	// Match the URL against the regex
	if !re.MatchString(str.String()) {
		return ""
	}
	if str.RawQuery != "" {
		return str.String() + "?" + str.RawQuery
	}
	return "" // Return empty if the URL does not match
}

func RandTimeSleep(i int) {
	i = i * 1000
	rand := rand.New((rand.NewSource(time.Now().UnixNano())))
	r := rand.Intn(i) + 100
	t := time.Duration(r)
	time.Sleep(time.Millisecond * t)
}
