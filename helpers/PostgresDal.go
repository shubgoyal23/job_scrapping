package helpers

import (
	"context"
	"fmt"
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
	if _, err := postgresConn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS job_listings (id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, job_title VARCHAR(255) NOT NULL, company_name VARCHAR(255) NOT NULL, company_url VARCHAR(255), job_description TEXT NOT NULL, job_type VARCHAR(225), location VARCHAR(255) NOT NULL, remote_option BOOLEAN DEFAULT FALSE, salary_min DECIMAL(10, 2), salary_max DECIMAL(10, 2), experience_min INT, experience_max INT, education_requirements TEXT[], skills TEXT[], benefits TEXT[], job_posting_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, application_deadline TIMESTAMP, job_url VARCHAR(255) NOT NULL UNIQUE, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, is_active BOOLEAN DEFAULT TRUE);"); err != nil {
		LogError("cannot create table", err)
		return err
	}
	return nil
}

// insert bulk data in postgres
func InsertBulkDataPostgres(val []types.JobListing) ([]string, error) {
	batch := &pgx.Batch{}
	failedRecords := []string{}
	for _, v := range val {
		batch.Queue("Insert into job_listings (job_title, company_name, company_url, job_description, job_type, location, remote_option, salary_min, salary_max, experience_min, experience_max, education_requirements, skills, benefits, job_posting_date, application_deadline, job_url, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)", v.JobTitle, v.CompanyName, v.CompanyURL, v.JobDescription, v.JobType, v.Location, v.RemoteOption, v.SalaryMin, v.SalaryMax, v.ExperienceMin, v.ExperienceMax, v.EducationRequirements, v.Skills, v.Benefits, v.JobPostingDate, v.ApplicationDeadline, v.JobURL, v.CreatedAt, v.UpdatedAt)
	}
	res := postgresConn.SendBatch(context.Background(), batch)
	defer res.Close()
	for _, v := range val {
		_, err := res.Exec()
		if err != nil {
			LogError(fmt.Sprintf("Failed to insert record %s", v.JobURL), err)
			failedRecords = append(failedRecords, v.JobURL)
		}
	}
	return failedRecords, nil
}

func GetManyDocPostgres(query string) ([]types.JobListing, error) {
	var val []types.JobListing
	r, err := postgresConn.Query(context.Background(), query)
	if err != nil {
		LogError("cannot get doc in postgres", err)
		return nil, err
	}
	defer r.Close()
	if r.Err() != nil {
		LogError("cannot get doc in postgres", r.Err())
		return nil, r.Err()
	}
	for r.Next() {
		var res types.JobListing
		if err := r.Scan(&res.ID, &res.JobTitle, &res.CompanyName, &res.CompanyURL, &res.JobDescription, &res.JobType, &res.Location, &res.RemoteOption, &res.SalaryMin, &res.SalaryMax, &res.ExperienceMin, &res.ExperienceMax, &res.EducationRequirements, &res.Skills, &res.Benefits, &res.JobPostingDate, &res.ApplicationDeadline, &res.JobURL, &res.CreatedAt, &res.UpdatedAt); err != nil {
			LogError("cannot decode doc in postgres", err)
			return nil, err
		}
		val = append(val, res)
	}
	return val, nil
}

func DeleteDocPostgres(query string, val int) error {
	_, err := postgresConn.Exec(context.Background(), query, val)
	if err != nil {
		LogError("cannot delete doc in postgres", err)
		return err
	}
	return nil
}

func UpdateDocPostgres(query string, val interface{}) error {
	_, err := postgresConn.Exec(context.Background(), query, val)
	if err != nil {
		LogError("cannot update doc in postgres", err)
		return err
	}
	return nil
}
