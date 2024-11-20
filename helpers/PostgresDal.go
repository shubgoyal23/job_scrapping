package helpers

import (
	"context"
	"nScrapper/types"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var postgresConn *pgxpool.Pool

// init postgres database
func InitPostgresDataBase() error {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		LogError("cannot connect to database", err)
		return err
	}
	postgresConn = conn
	if err := postgresConn.Ping(context.Background()); err != nil {
		LogError("cannot connect to database", err)
		return err
	}
	return nil
}

func CreateTablePostgres() error {
	if _, err := postgresConn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS job_listings (id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, job_title VARCHAR(255) NOT NULL, company_name VARCHAR(255) NOT NULL, company_url VARCHAR(255), job_description TEXT NOT NULL, job_type VARCHAR(50), location VARCHAR(255) NOT NULL, remote_option BOOLEAN DEFAULT FALSE, salary_min DECIMAL(10, 2), salary_max DECIMAL(10, 2), experience_min INT, experience_max INT, education_requirements TEXT[], skills TEXT[], benefits TEXT[], job_posting_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, application_deadline TIMESTAMP, job_url VARCHAR(255) NOT NULL UNIQUE, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"); err != nil {
		LogError("cannot create table", err)
		return err
	}
	return nil
}

// insert bulk data in postgres
func InsertBulkDataPostgres(val []types.JobListing) error {
	batch := &pgx.Batch{}
	for _, v := range val {
		batch.Queue("Insert into job_listings (job_title, company_name, company_url, job_description, job_type, location, remote_option, salary_min, salary_max, experience_min, experience_max, education_requirements, skills, benefits, job_posting_date, application_deadline, job_url, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)", v.JobTitle, v.CompanyName, v.CompanyURL, v.JobDescription, v.JobType, v.Location, v.RemoteOption, v.SalaryMin, v.SalaryMax, v.ExperienceMin, v.ExperienceMax, v.EducationRequirements, v.Skills, v.Benefits, v.JobPostingDate, v.ApplicationDeadline, v.JobURL, v.CreatedAt, v.UpdatedAt)
	}
	res := postgresConn.SendBatch(context.Background(), batch)
	if err := res.Close(); err != nil {
		LogError("cannot insert", err)
		return err
	}
	return nil
}