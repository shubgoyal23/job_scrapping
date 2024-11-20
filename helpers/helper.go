package helpers

import (
	"log"
	"nScrapper/types"
	"os"
	"time"
)

var Logger *os.File

func LogError(strings string, err error) {
	if err == nil {
		str := time.Now().Local().Format("2006-01-02 15:04:05") + ": " + strings + "\n"
		Logger.Write([]byte(str))
		return
	}
	str := time.Now().String() + ": " + strings + " " + err.Error() + "\n"
	Logger.Write([]byte(str))
}

func InitLogger() *os.File {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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
		LogError("cannot insert in mongodb", err)
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
		LogError("cannot insert in mongodb", err)
	}
}
