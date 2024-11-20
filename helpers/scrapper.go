package helpers

import (
	"fmt"
	"nScrapper/types"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var UniqueTags = make(chan string, 500) // chan to store unique tags
var Headless = false                    // to run in headless mode
var ScrapeMap = make(map[string]types.JobDataScrapeMap)

// this function collects all the data from the page
func ScrapperElements(page *rod.Page, jobDMap types.JobDataScrapeMap) types.JobListing {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			LogError("naukri.com "+fmt.Sprintf("begin handler crashed because %s", r), nil)
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
					LogError("cannot get "+fname, err)
				}
			}
		} else if fVal.TagType == "numeric" {
			if element, err := page.Timeout(5 * time.Second).Element(fVal.Element); err == nil {
				if text, err := element.Text(); err == nil {
					text = CleanText(text, fVal.Cleaner)
					SetField(&jobDetails, fname, text)
				} else {
					LogError("cannot get "+fname, err)
				}
			}
		} else if fVal.TagType == "url" {
			if element, err := page.Timeout(5 * time.Second).Element(fVal.Element); err == nil {
				if text, err := element.Attribute(fVal.AttributeTarget); err == nil {
					text1 := CleanUrl(*text, jobDMap.Homepage)
					SetField(&jobDetails, fname, text1)
				} else {
					LogError("cannot get "+fname, err)
				}
			}
		} else if fVal.TagType == "date" {
			if element, err := page.Timeout(5 * time.Second).Element(fVal.Element); err == nil {
				if text, err := element.Text(); err == nil {
					var date time.Time
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
					} else {
						date = time.Now()
					}
					SetField(&jobDetails, fname, date)
				} else {
					LogError("cannot get "+fname, err)
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
						LogError("cannot get "+fname, err)
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
					LogError("cannot get "+fname, err)
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
					LogError("cannot get "+fname, err)
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
			LogError(fmt.Sprintf("LinkDupper handler crashed because %s", r), nil)
		}
	}()
	LogError(fmt.Sprintf("running linkDuper for %s scrapper at time: %s", jobMap.Homepage, time.Now().String()), nil)

	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).Headless(Headless).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(jobMap.Homepage).MustWaitStable()

	AllTags := make(map[string]bool)
	for _, pl := range jobMap.PageLinks {
		pageNErr := page.Navigate(pl.Link)
		if pageNErr != nil {
			LogError(fmt.Sprintf("error while navigating to page: %s", pl.Link), pageNErr)
			continue
		}
		for {
			if err := page.Timeout(30 * time.Second).WaitLoad(); err != nil {
				LogError("error while waiting for page to be stable", err)
				break
			}
			aTags, aTagErr := page.Timeout(30 * time.Second).Elements(pl.Element)
			if aTagErr != nil {
				LogError(fmt.Sprintf("error while getting elements from page: %s, element: %s, err:", pl.Link, pl.Element), aTagErr)
			}
			for _, a := range aTags {
				aTag, aTagErr := a.Attribute("href")
				if aTag == nil {
					continue
				}
				if aTagErr != nil {
					LogError(fmt.Sprintf("error while getting href from page: %s, element: %s, err:", pl.Link, pl.Element), aTagErr)
				}
				lk := CleanUrl(*aTag, jobMap.Homepage)
				if lk == "" {
					continue
				}
				if _, ok := AllTags[lk]; ok {
					continue
				}

				AllTags[lk] = true
				UniqueTags <- lk
			}

			// next button click
			nextBtn, nextBtnErr := page.Timeout(30 * time.Second).Element(pl.NextPageBtn)
			if nextBtnErr != nil {
				LogError(fmt.Sprintf("error while getting next page button from page: %s, element: %s, err:", pl.Link, pl.NextPageBtn), nextBtnErr)
				break
			}
			nextBtnErr = nextBtn.Click("left", 1)
			if nextBtnErr != nil {
				LogError(fmt.Sprintf("error while clicking next page button from page: %s, element: %s, err:", pl.Link, pl.NextPageBtn), nextBtnErr)
				break
			}

		}
	}
}

// this function gets the data from the links in chan UniqueTags and stores it in postgres and redis set
func GetDataFromLink() {
	defer func() {
		if r := recover(); r != nil {
			LogError(fmt.Sprintf("GetDataFromLink crashed because %s", r), nil)
		}
	}()
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).Headless(Headless).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()
	page := browser.MustPage()

	var jobData []types.JobListing
	var jobLinks []string

	for link := range UniqueTags {
		// check if link exists in redis
		if f, err := CheckRedisSetMemeber("posted_job_links", link); err != nil || f {
			LogError("cannot get from redis", err)
			continue
		}

		pageNErr := page.Timeout(30 * time.Second).Navigate(link)
		if pageNErr != nil {
			LogError(fmt.Sprintf("error while navigating to page: %s", link), pageNErr)
			continue
		}
		if err := page.Timeout(30 * time.Second).WaitLoad(); err != nil {
			LogError("error while waiting for page to be stable", err)
			continue
		}
		u, err := url.Parse(link)
		if err != nil {
			LogError("cannot parse url", err)
			continue
		}
		homeDomain := u.Scheme + "://" + u.Host
		sMap, ok := ScrapeMap[homeDomain]
		if !ok {
			LogError(fmt.Sprintf("cannot get ScrapeMap for %s", homeDomain), nil)
			continue
		}

		// scrape elements on page
		da := ScrapperElements(page, sMap)
		da.JobURL = link
		da.ApplicationDeadline = da.JobPostingDate.AddDate(0, 0, 30)
		da.CreatedAt = time.Now()
		da.UpdatedAt = time.Now()
		jobData = append(jobData, da)
		jobLinks = append(jobLinks, link)
		if len(jobData) == 100 {
			if err := InsertBulkDataPostgres(jobData); err != nil {
				LogError("cannot insert job data", err)
			} else {
				InsertRedisSetBulk("posted_job_links", jobLinks)
			}
			jobData = []types.JobListing{}
			jobLinks = []string{}
		}
	}
	if err := InsertBulkDataPostgres(jobData); err != nil {
		LogError("cannot insert job data", err)
	} else {
		InsertRedisSetBulk("posted_job_links", jobLinks)
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
			LogError("ourl", oerr)
			return ""
		}
		// Check if the cleaned URL matches the regex
		if re.MatchString(ourl) {
			return ourl
		}
		return ""
	}

	// Ensure the URL uses HTTPS or HTTP
	if str.Scheme != "https" && str.Scheme != "http" {
		str.Scheme = "https"
	}

	// Parse the home URL to extract the hostname
	h, err := url.Parse(home_url)
	if err != nil || str.Hostname() != h.Hostname() {
		return "" // Return empty if hostnames do not match or if home_url is invalid
	}

	// Match the URL against the regex
	if re.MatchString(str.String()) {
		return str.String() // Valid and cleaned URL
	}
	return "" // Return empty if the URL does not match
}
