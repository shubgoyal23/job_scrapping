package types

import (
	"time"
)

// JobListing represents a job listing in the database.
type JobListing struct {
	ID                    int       `json:"id"`                     // Unique identifier for each job listing
	JobTitle              string    `json:"job_title"`              // Title of the job
	CompanyName           string    `json:"company_name"`           // Name of the company offering the job
	CompanyURL            string    `json:"company_url"`            // URL to the company's website
	JobDescription        string    `json:"job_description"`        // Detailed description of the job
	JobType               string    `json:"job_type"`               // Type of job (e.g., Full-time, Part-time, Contract, Internship)
	Location              string    `json:"location"`               // Location of the job (city, state, country)
	RemoteOption          bool      `json:"remote_option"`          // Is the job remote? (true/false)
	SalaryMin             float64   `json:"salary_min"`             // Minimum salary for the job
	SalaryMax             float64   `json:"salary_max"`             // Maximum salary for the job
	ExperienceMin         float64   `json:"experience_min"`         // Minimum experience required (in years)
	ExperienceMax         float64   `json:"experience_max"`         // Maximum experience required (in years)
	EducationRequirements []string  `json:"education_requirements"` // Required education qualifications
	Skills                []string  `json:"skills"`                 // Required skills
	Benefits              []string  `json:"benefits"`               // Benefits offered by the company
	JobPostingDate        time.Time `json:"job_posting_date"`       // Date and time when the job was posted
	ApplicationDeadline   time.Time `json:"application_deadline"`   // Deadline for applications
	JobURL                string    `json:"job_url"`                // URL to the job listing on the website
	CreatedAt             time.Time `json:"created_at"`             // Record creation timestamp
	UpdatedAt             time.Time `json:"updated_at"`             // Record update timestamp
}

type JobListingFeilds struct {
	JobTitle              TagField `json:"job_title" bson:"job_title"`
	CompanyName           TagField `json:"company_name" bson:"company_name"`
	CompanyURL            TagField `json:"company_url" bson:"company_url"`
	JobDescription        TagField `json:"job_description" bson:"job_description"`
	JobType               TagField `json:"job_type" bson:"job_type"`
	Location              TagField `json:"location" bson:"location"`
	RemoteOption          TagField `json:"remote_option" bson:"remote_option"`
	SalaryMin             TagField `json:"salary_min" bson:"salary_min"`
	SalaryMax             TagField `json:"salary_max" bson:"salary_max"`
	ExperienceMin         TagField `json:"experience_min" bson:"experience_min"`
	ExperienceMax         TagField `json:"experience_max" bson:"experience_max"`
	EducationRequirements TagField `json:"education_requirements" bson:"education_requirements"`
	Skills                TagField `json:"skills" bson:"skills"`
	Benefits              TagField `json:"benefits" bson:"benefits"`
	JobPostingDate        TagField `json:"job_posting_date" bson:"job_posting_date"`
	ApplicationDeadline   TagField `json:"application_deadline" bson:"application_deadline"`
	JobURL                TagField `json:"job_url" bson:"job_url"`
}

type TagField struct {
	Element         string `json:"element" bson:"element"`
	TagType         string `json:"tagtype" bson:"tagtype"`
	IsAttribute     string `json:"isattribute" bson:"isattribute"`
	AttributeTarget string `json:"attribute_target" bson:"attribute_target"`
	Cleaner         string `json:"cleaner" bson:"cleaner"`
}

type PageLinks struct {
	Link        string `json:"link" bson:"link"`
	NextPageBtn string `json:"total_pages" bson:"total_pages"`
	Element     string `json:"element" bson:"element"`
}

type JobDataScrapeMap struct {
	Homepage  string           `json:"homepage" bson:"homepage"`
	PageLinks []PageLinks      `json:"page_links" bson:"page_links"`
	JobData   JobListingFeilds `json:"job_data" bson:"job_data"`
}
