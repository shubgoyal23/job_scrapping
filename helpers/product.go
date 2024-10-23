package helpers

import (
	"encoding/json"
	"fmt"
	"nScrapper/types"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func ProductScrapper(page *rod.Page, scrapeFeilds []string) {
	err := page.WaitIdle(2 * time.Minute)
	info, _ := page.Info()
	erl := info.URL
	if err != nil {
		LogError(fmt.Sprintf("ProductScrapper: page failde to load: %s", erl), err)
	}

	elements := NaukriElements(page)
	Insert(elements)

	// for _, scrapeFeild := range scrapeFeilds {
	// 	feild, err := page.Element(scrapeFeild)
	// 	if err != nil {
	// 		LogError(fmt.Sprintf("ProductScrapper: feild not found: %s, on page: %s", scrapeFeild, erl), err)
	// 	}
	// 	feildText, err := feild.Text()
	// 	if err != nil {
	// 		LogError(fmt.Sprintf("ProductScrapper: feild not found: %s, on page: %s", scrapeFeild, erl), err)
	// 	}
	// 	fmt.Printf("%s: %s\n", scrapeFeild, feildText)
	// }
}

func NaukriElements(page *rod.Page) (jobD types.JobListing) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			LogError("naukri.com "+fmt.Sprintf("begin handler crashed because %s", r), nil)
		}
	}()

	// Change first element (Job Title) and continue to apply logic for subsequent elements
	if titleElement, err := page.Timeout(5 * time.Second).Element(".styles_jd-header-title__rZwM1"); err != nil {
		LogError("cannot get job title", err)
	} else {
		if tt, te := titleElement.Text(); te != nil {
			LogError("cannot get job title", te)
		} else {
			jobD.JobTitle = tt
		}
	}

	// After the first element is changed, continue with the rest, including the Timeout for each element.
	if companyNameElement, err := page.Timeout(5 * time.Second).Element(".styles_jd-header-comp-name__MvqAI > a"); err == nil {
		if cn, ce := companyNameElement.Text(); ce == nil {
			jobD.CompanyName = cn
		} else {
			LogError("cannot get company name", ce)
		}
	}

	if companyURLElement, err := page.Timeout(5 * time.Second).Element(".styles_jd-header-comp-name__MvqAI > a"); err == nil {
		if cu, ce := companyURLElement.Attribute("href"); ce == nil {
			jobD.CompanyURL = string(*cu)
		} else {
			LogError("cannot get company url", ce)
		}
	}

	if jobDescriptionElement, err := page.Timeout(5 * time.Second).Element("#root > div > main > div.styles_jdc__content__EZJMQ > div.styles_left-section-container__btAcB > section.styles_job-desc-container__txpYf > div:nth-child(2)"); err == nil {
		if jd, je := jobDescriptionElement.Text(); je == nil {
			jobD.JobDescription = jd
		} else {
			LogError("cannot get job description", je)
		}
	}

	if jobTypeElement, err := page.Timeout(5 * time.Second).Element("#root > div > main > div.styles_jdc__content__EZJMQ > div.styles_left-section-container__btAcB > section.styles_job-desc-container__txpYf > div:nth-child(2) > div.styles_other-details__oEN4O > div:nth-child(4) > span > span"); err == nil {
		if jt, je := jobTypeElement.Text(); je == nil {
			jobD.JobType = jt
		} else {
			LogError("cannot get job type", je)
		}
	}

	if jobLocationElement, err := page.Timeout(5 * time.Second).Element(".styles_jhc__loc___Du2H"); err == nil {
		if jl, je := jobLocationElement.Text(); je == nil {
			jobD.Location = jl
		} else {
			LogError("cannot get job location", je)
		}
	}

	if remoteElement, err := page.Timeout(5 * time.Second).Element(".styles_jhc__wfhmode-link__aHmrK"); err == nil {
		if _, re := remoteElement.Text(); re == nil {
			jobD.RemoteOption = true
		}
	}

	if salaryElement, err := page.Timeout(5 * time.Second).Element(".styles_jhc__salary__jdfEC"); err == nil {
		if salaryText, se := salaryElement.Text(); se == nil {
			if strings.Contains(salaryText, "-") {
				str := strings.Split(salaryText, "-")
				if min, err := strconv.Atoi(str[0]); err == nil {
					jobD.SalaryMin = float64(min)
				}
				if min, err := strconv.Atoi(str[1]); err == nil {
					jobD.SalaryMax = float64(min)
				}

			}
		} else {
			LogError("cannot get salary", se)
		}
	}

	if jobExperienceMinElement, err := page.Timeout(5 * time.Second).Element(".styles_jhc__exp__k_giM"); err == nil {
		if jel, je := jobExperienceMinElement.Text(); je == nil {
			if strings.Contains(jel, "years") {
				str := strings.Split(jel, " ")
				if min, err := strconv.Atoi(str[0]); err == nil {
					jobD.ExperienceMin = min
				}
				if max, err := strconv.Atoi(str[2]); err == nil {
					jobD.ExperienceMax = max
				}
			}
		} else {
			LogError("cannot get job experience", je)
		}
	}

	if jobEducationRequirementsElement, err := page.Timeout(5 * time.Second).Element(".styles_education__KXFkO"); err == nil {
		if jeE, je := jobEducationRequirementsElement.Text(); je == nil {
			jobD.EducationRequirements = jeE
		} else {
			LogError("cannot get job education requirements", je)
		}
	}

	if jobSkillsElement, err := page.Timeout(5 * time.Second).Element("#root > div > main > div.styles_jdc__content__EZJMQ > div.styles_left-section-container__btAcB > section.styles_job-desc-container__txpYf > div:nth-child(2) > div.styles_JDC__dang-inner-html__h0K4t > ul:nth-child(23)"); err == nil {
		if js, se := jobSkillsElement.Text(); se == nil {
			jobD.Skills = strings.Split(js, "\n")
		} else {
			LogError("cannot get job skills", se)
		}
	}

	if jobBenefitsElement, err := page.Timeout(5 * time.Second).Element(".styles_jhc__benefits__jdfEC"); err == nil {
		if jb, be := jobBenefitsElement.Text(); be == nil {
			jobD.Benefits = strings.Split(jb, "\n")
		} else {
			LogError("cannot get job benefits", be)
		}
	}

	// Todo

	// if jobApplicationDeadlineElement, err := page.Timeout(5 * time.Second).Element(".styles_jhc__application-deadline__jdfEC"); err == nil {
	// 	if jad, je := jobApplicationDeadlineElement.Text(); je == nil {
	// 		jobD.ApplicationDeadline = jad
	// 	} else {
	// 		LogError("cannot get job application deadline", je)
	// 	}
	// }

	if jobCreatedAtElement, err := page.Timeout(5 * time.Second).Element("#job_header > div.styles_jhc__bottom__DrTmB > div.styles_jhc__jd-stats__KrId0 > span:nth-child(1) > span"); err == nil {
		if jc, ce := jobCreatedAtElement.Text(); ce == nil {
			jc = strings.ToLower(jc)
			if strings.Contains(jc, "ago") {
				str := strings.Split(jc, " ")
				age := str[0]
				if agoTago, e := strconv.Atoi(age); e == nil {
					jobD.JobPostingDate = time.Now().AddDate(0, 0, agoTago*-1)
				} else {
					jobD.JobPostingDate = time.Now()
				}
			}
		} else {
			LogError("cannot get job created at", ce)
		}
	}

	if jobJobURLElement, err := page.Timeout(5 * time.Second).Element(".styles_jhc__job-url__jdfEC"); err == nil {
		if jj, je := jobJobURLElement.Text(); je == nil {
			jobD.JobURL = jj
		} else {
			LogError("cannot get job url", je)
		}
	}
	info, err := page.Info()
	if err != nil {
		LogError("cannot get page info", err)
	}
	jobD.JobURL = info.URL
	jobD.CreatedAt = time.Now()

	// PrintToJson(jobD)
	return jobD
}

// func NaukriElemets(page *rod.Page) (jobD types.JobListingFeilds) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			fmt.Println("Recovered in f", r)
// 			LogError("naukri.com "+fmt.Sprintf("begin handler crashed because %s", r), nil)
// 		}
// 	}()

// 	if titleElement, err := page.Timeout(60 * time.Second).Element(".styles_jd-header-title__rZwM1"); err != nil {
// 		LogError("cannot get job title", err)
// 	} else {
// 		if tt, te := titleElement.Text(); te != nil {
// 			LogError("cannot get job title", te)
// 		} else {
// 			jobD.JobTitle = tt
// 		}
// 	}

// 	companyNameElement, _ := page.Timeout(60 * time.Second).Element(".styles_jd-header-comp-name__MvqAI > a")
// 	cn, ce := companyNameElement.Text()
// 	if ce != nil {
// 		LogError("cannot get company name", ce)
// 	}
// 	jobD.CompanyName = cn

// 	companyURLElement, _ := page.Timeout(60 * time.Second).Element(".styles_jd-header-comp-name__MvqAI > a")
// 	cu, ce := companyURLElement.Attribute("href")
// 	if ce != nil {
// 		LogError("cannot get company url", ce)
// 	}
// 	jobD.CompanyURL = string(*cu)

// 	jobDescriptionElement, _ := page.Timeout(60 * time.Second).Element("#root > div > main > div.styles_jdc__content__EZJMQ > div.styles_left-section-container__btAcB > section.styles_job-desc-container__txpYf > div:nth-child(2)")
// 	jd, je := jobDescriptionElement.Text()
// 	if je != nil {
// 		LogError("cannot get job description", je)
// 	}
// 	jobD.JobDescription = jd

// 	jobTypeElement, _ := page.Timeout(60 * time.Second).Element("#root > div > main > div.styles_jdc__content__EZJMQ > div.styles_left-section-container__btAcB > section.styles_job-desc-container__txpYf > div:nth-child(2) > div.styles_other-details__oEN4O > div:nth-child(4) > span > span")
// 	jt, je := jobTypeElement.Text()
// 	if je != nil {
// 		LogError("cannot get job type", je)
// 	}
// 	jobD.JobType = jt

// 	jobLocationElement, _ := page.Timeout(60 * time.Second).Element(".styles_jhc__loc___Du2H")
// 	jl, je := jobLocationElement.Text()
// 	if je != nil {
// 		LogError("cannot get job location", je)
// 	}
// 	jobD.Location = jl

// 	remoteElement, _ := page.Timeout(60 * time.Second).Element(".styles_jhc__loc___Du2H")
// 	remote, re := remoteElement.Text()
// 	if re != nil {
// 		LogError("cannot get remote option", re)
// 	}
// 	jobD.RemoteOption = remote

// 	salary, err := page.Timeout(60 * time.Second).Element(".styles_jhc__salary__jdfEC")
// 	if err != nil {
// 		LogError("cannot get salary", err)
// 	}
// 	salaryText, err := salary.Text()
// 	if err != nil {
// 		LogError("cannot get salary", err)
// 	}
// 	jobD.SalaryMin = salaryText
// 	jobExperienceMinElement, _ := page.Timeout(60 * time.Second).Element(".styles_jhc__exp__k_giM")
// 	jel, err := jobExperienceMinElement.Text()
// 	if err != nil {
// 		LogError("cannot get job experience", err)
// 	}
// 	jobD.ExperienceMin = jel
// 	jobEducationRequirementsElement, _ := page.Timeout(60 * time.Second).Element(".styles_education__KXFkO")
// 	jeE, err := jobEducationRequirementsElement.Text()
// 	if err != nil {
// 		LogError("cannot get job education requirements", err)
// 	}
// 	jobD.EducationRequirements = jeE
// 	jobSkillsElement, err := page.Timeout(60 * time.Second).Element(".styles_jhc__skills__jdfEC")
// 	js, err := jobSkillsElement.Text()
// 	if err != nil {
// 		LogError("cannot get job skills", err)
// 	}
// 	jobD.Skills = js
// 	jobBenefitsElement, _ := page.Timeout(60 * time.Second).Element(".styles_jhc__benefits__jdfEC")
// 	jb, err := jobBenefitsElement.Text()
// 	if err != nil {
// 		LogError("cannot get job benefits", err)
// 	}
// 	jobD.Benefits = jb
// 	jobJobPostingDateElement, _ := page.Timeout(60 * time.Second).Element(".styles_jhc__job-posting-date__jdfEC")
// 	jjd, err := jobJobPostingDateElement.Text()
// 	if err != nil {
// 		LogError("cannot get job job posting date", err)
// 	}
// 	jobD.JobPostingDate = jjd
// 	jobApplicationDeadlineElement, _ := page.Timeout(60 * time.Second).Element(".styles_jhc__application-deadline__jdfEC")
// 	jad, err := jobApplicationDeadlineElement.Text()
// 	if err != nil {
// 		LogError("cannot get job application deadline", err)
// 	}
// 	jobD.ApplicationDeadline = jad
// 	jobCreatedAtElement, _ := page.Timeout(60 * time.Second).Element(".styles_jhc__created-at__jdfEC")
// 	jc, err := jobCreatedAtElement.Text()
// 	if err != nil {
// 		LogError("cannot get job created at", err)
// 	}
// 	jobD.CreatedAt = jc
// 	jobUpdatedAtElement, _ := page.Timeout(60 * time.Second).Element(".styles_jhc__updated-at__jdfEC")
// 	ju, err := jobUpdatedAtElement.Text()
// 	if err != nil {
// 		LogError("cannot get job updated at", err)
// 	}
// 	jobD.UpdatedAt = ju
// 	jobJobURLElement, _ := page.Timeout(60 * time.Second).Element(".styles_jhc__job-url__jdfEC")
// 	jj, err := jobJobURLElement.Text()
// 	if err != nil {
// 		LogError("cannot get job job url", err)
// 	}
// 	jobD.JobURL = jj

// 	PrintToJson(jobD)
// 	return
// }

func PrintToJson(product types.JobListing) {
	b, err := json.MarshalIndent(product, "", "\t")
	if err != nil {
		LogError("cannot marshal product", err)
	}
	file, _ := os.Create("product.json")
	file.Write(b)
	file.Close()
}
