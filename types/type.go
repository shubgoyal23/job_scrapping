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
	ExperienceMin         int       `json:"experience_min"`         // Minimum experience required (in years)
	ExperienceMax         int       `json:"experience_max"`         // Maximum experience required (in years)
	EducationRequirements string    `json:"education_requirements"` // Required education qualifications
	Skills                []string  `json:"skills"`                 // Required skills
	Benefits              []string  `json:"benefits"`               // Benefits offered by the company
	JobPostingDate        time.Time `json:"job_posting_date"`       // Date and time when the job was posted
	ApplicationDeadline   time.Time `json:"application_deadline"`   // Deadline for applications
	JobURL                string    `json:"job_url"`                // URL to the job listing on the website
	CreatedAt             time.Time `json:"created_at"`             // Record creation timestamp
	UpdatedAt             time.Time `json:"updated_at"`             // Record update timestamp
}
