package helpers

import (
	"log"
	"nScrapper/types"
	"os"
	"time"
)

var Logger *os.File

func LogError(fu string, strings string, err error) {
	if err == nil {
		str := time.Now().Local().Format("2006-01-02 15:04:05") + ": " + fu + ": Message : " + strings + "\n"
		Logger.Write([]byte(str))
		return
	}
	str := time.Now().Local().Format("2006-01-02 15:04:05") + ": " + fu + ": Message : " + strings + ": Error : " + err.Error() + "\n"
	Logger.Write([]byte(str))
}

func InitLogger() *os.File {
	os.MkdirAll("./log", os.ModePerm)
	file, err := os.OpenFile("./log/log.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	Logger = file
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return Logger
}

func InsertNaukriMapToMongoDB() {
	var jg types.JobDataScrapeMap
	jg.Homepage = "https://www.naukri.com"
	jg.PageLinks = []types.PageLinks{{Link: "https://www.naukri.com/it-jobs-?src=gnbjobs_homepage_srch", NextPageBtn: "#lastCompMark > a:nth-child(4)", Element: "#listContainer > div.styles_job-listing-container__OCfZC > div > div > div > div > a"}}

	var jobData types.JobListingFeilds
	jobData.JobTitle = types.TagField{Element: ".styles_jd-header-title__rZwM1", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.CompanyName = types.TagField{Element: ".styles_jd-header-comp-name__MvqAI > a", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.CompanyURL = types.TagField{Element: ".styles_jd-header-comp-name__MvqAI > a", TagType: "url", Cleaner: "", AttributeTarget: "href"}
	jobData.JobDescription = types.TagField{Element: "#root > div > main > div.styles_jdc__content__EZJMQ > div.styles_left-section-container__btAcB > section.styles_job-desc-container__txpYf > div:nth-child(2)", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.JobType = types.TagField{Element: "#root > div > main > div.styles_jdc__content__EZJMQ > div.styles_left-section-container__btAcB > section.styles_job-desc-container__txpYf > div:nth-child(2) > div.styles_other-details__oEN4O > div:nth-child(4) > span > span", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.Location = types.TagField{Element: ".styles_jhc__loc___Du2H", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.RemoteOption = types.TagField{Element: ".styles_jhc__wfhmode-link__aHmrK", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.SalaryMin = types.TagField{Element: ".styles_jhc__salary__jdfEC", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.SalaryMax = types.TagField{Element: ".styles_jhc__salary__jdfEC", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.ExperienceMin = types.TagField{Element: ".styles_jhc__exp__k_giM", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.ExperienceMax = types.TagField{Element: ".styles_jhc__exp__k_giM", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.EducationRequirements = types.TagField{Element: ".styles_education__KXFkO > div.styles_details__Y424J", TagType: "[]string", Cleaner: "", AttributeTarget: ""}
	jobData.Skills = types.TagField{Element: "div.styles_key-skill__GIPn_ > div > a > span", TagType: "[]string", Cleaner: "", AttributeTarget: ""}
	jobData.Benefits = types.TagField{Element: ".styles_jhc__benefits__jdfEC", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.JobPostingDate = types.TagField{Element: "#job_header > div.styles_jhc__bottom__DrTmB > div.styles_jhc__jd-stats__KrId0 > span:nth-child(1) > span", TagType: "date", Cleaner: "", AttributeTarget: ""}
	jobData.ApplicationDeadline = types.TagField{Element: "", TagType: "date", Cleaner: "", AttributeTarget: ""}

	jg.JobData = jobData

	// Insert the data into MongoDB
	if err := InsertMongoDB(jg); err != nil {
		LogError("InsertNaukriMapToMongoDB", "cannot insert in mongodb", err)
	}
}

func InsertFounditMapToMongoDB() {
	var jg types.JobDataScrapeMap
	jg.Homepage = "https://www.foundit.in"
	jg.PageLinks = []types.PageLinks{{Link: "https://www.foundit.in/search/it-jobs", NextPageBtn: "#pagination > a > div.arrow-right", Element: "div.srpResultCard > div > div > div > div.cardHead > div > div > h3 > a"}}

	var jobData types.JobListingFeilds
	jobData.JobTitle = types.TagField{Element: "#jobDetailContainer > div > div > h1", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.CompanyName = types.TagField{Element: "#jobDetailContainer > div > div > a", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.CompanyURL = types.TagField{Element: "#jobDetailContainer > div > div > a", TagType: "url", Cleaner: "", AttributeTarget: "href"}
	jobData.JobDescription = types.TagField{Element: "#jobDescription > div", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.JobType = types.TagField{Element: "#jobDetailContainer > div > div > p:nth-child(3) > span.font-normal.text-content-secondary > a", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.Location = types.TagField{Element: "#jobDetailContainer > div > div:nth-child(2) > div > a:nth-child(2)", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.RemoteOption = types.TagField{Element: "", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.SalaryMin = types.TagField{Element: "", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.SalaryMax = types.TagField{Element: "", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.ExperienceMin = types.TagField{Element: "#jobDetailContainer > div > div > div > span:nth-child(2)", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.ExperienceMax = types.TagField{Element: "#jobDetailContainer > div > div > div > span:nth-child(2)", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.EducationRequirements = types.TagField{Element: "", TagType: "[]string", Cleaner: "", AttributeTarget: ""}
	jobData.Skills = types.TagField{Element: "#skillSectionNew > div.flex.flex-wrap > div", TagType: "[]string", Cleaner: "", AttributeTarget: ""}
	jobData.Benefits = types.TagField{Element: "", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.JobPostingDate = types.TagField{Element: "#jobDetailContainer > div > div > div > ul > li:nth-child(1) > span", TagType: "date", Cleaner: "", AttributeTarget: ""}
	jobData.ApplicationDeadline = types.TagField{Element: "", TagType: "date", Cleaner: "", AttributeTarget: ""}

	jg.JobData = jobData

	// Insert the data into MongoDB
	if err := InsertMongoDB(jg); err != nil {
		LogError("InsertFounditMapToMongoDB", "cannot insert in mongodb", err)
	}
}
func InsertIndeedMapToMongoDB() {
	var jg types.JobDataScrapeMap
	var jobData types.JobListingFeilds
	jg.Homepage = "https://in.indeed.com"
	jg.PageLinks = []types.PageLinks{{Link: "https://in.indeed.com/jobs?q=backend+developer", NextPageBtn: "a[data-testid=pagination-page-next]", Element: "h2.jobTitle > a"}}

	jobData.JobTitle = types.TagField{Element: "div.jobsearch-JobInfoHeader-title-container.css-bbq8li.eu4oa1w0 > h1 > span", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.CompanyName = types.TagField{Element: "#viewJobSSRRoot > div.fastviewjob.jobsearch-ViewJobLayout--standalone.css-r07ztj.eu4oa1w0.hydrated > div.css-1quav7f.eu4oa1w0 > div > div > div.jobsearch-JobComponent.css-u4y1in.eu4oa1w0 > div.jobsearch-InfoHeaderContainer.jobsearch-DesktopStickyContainer.css-zt53js.eu4oa1w0 > div:nth-child(1) > div.css-1moflg3.eu4oa1w0 > div > div > div > div.css-1h46us2.eu4oa1w0 > div > span > a", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.CompanyURL = types.TagField{Element: "#viewJobSSRRoot > div.fastviewjob.jobsearch-ViewJobLayout--standalone.css-r07ztj.eu4oa1w0.hydrated > div.css-1quav7f.eu4oa1w0 > div > div > div.jobsearch-JobComponent.css-u4y1in.eu4oa1w0 > div.jobsearch-InfoHeaderContainer.jobsearch-DesktopStickyContainer.css-zt53js.eu4oa1w0 > div:nth-child(1) > div.css-1moflg3.eu4oa1w0 > div > div > div > div.css-1h46us2.eu4oa1w0 > div > span > a", TagType: "url", Cleaner: "", AttributeTarget: "href"}
	jobData.JobDescription = types.TagField{Element: "#jobDescriptionText", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.JobType = types.TagField{Element: "#jobDetailsSection > div > div.js-match-insights-provider-36vfsm.eu4oa1w0 > div.js-match-insights-provider-h05mm8.e37uo190 > div:nth-child(2) > div > div > ul > li > div > div > div:nth-child(1)", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.Location = types.TagField{Element: "#viewJobSSRRoot > div.fastviewjob.jobsearch-ViewJobLayout--standalone.css-r07ztj.eu4oa1w0.hydrated > div.css-1quav7f.eu4oa1w0 > div > div > div.jobsearch-JobComponent.css-u4y1in.eu4oa1w0 > div.jobsearch-InfoHeaderContainer.jobsearch-DesktopStickyContainer.css-zt53js.eu4oa1w0 > div:nth-child(1) > div.css-1moflg3.eu4oa1w0 > div > div > div > div.css-waniwe.eu4oa1w0 > div", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.RemoteOption = types.TagField{Element: "", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.SalaryMin = types.TagField{Element: "#salaryInfoAndJobType > span", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.SalaryMax = types.TagField{Element: "#salaryInfoAndJobType > span", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.ExperienceMin = types.TagField{Element: "", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.ExperienceMax = types.TagField{Element: "", TagType: "range", Cleaner: "[^0-9.-]+", AttributeTarget: ""}
	jobData.EducationRequirements = types.TagField{Element: "", TagType: "[]string", Cleaner: "", AttributeTarget: ""}
	jobData.Skills = types.TagField{Element: "", TagType: "[]string", Cleaner: "", AttributeTarget: ""}
	jobData.Benefits = types.TagField{Element: "", TagType: "string", Cleaner: "", AttributeTarget: ""}
	jobData.JobPostingDate = types.TagField{Element: "", TagType: "date", Cleaner: "", AttributeTarget: ""}
	jobData.ApplicationDeadline = types.TagField{Element: "", TagType: "date", Cleaner: "", AttributeTarget: ""}

	jg.JobData = jobData

	// Insert the data into MongoDB
	if err := InsertMongoDB(jg); err != nil {
		LogError("InsertIndeedMapToMongoDB", "cannot insert in mongodb", err)
	}
}
