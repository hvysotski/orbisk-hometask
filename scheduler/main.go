package main

import (
	"api"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
)

const (
	DefaultDatabaseName = "db.sqlite"
)

var (
	DataDir = getEnvDefault("DATA_DIR", "../data")
)

type Job struct {
	Name     string
	Schedule string
	Query    string
	Handler  func(result interface{}) error
}

type Scheduler struct {
	cron *cron.Cron
	db   *sql.DB
	jobs []Job
}

func NewScheduler(dbPath string) (*Scheduler, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		cron: cron.New(),
		db:   db,
	}, nil
}

func (s *Scheduler) AddJob(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) Start() {
	for _, job := range s.jobs {
		j := job
		s.cron.AddFunc(j.Schedule, func() {
			result, err := s.executeQuery(j.Query)
			if err != nil {
				log.Printf("Error executing job %s: %v", j.Name, err)
				return
			}
			err = j.Handler(result)
			if err != nil {
				log.Printf("Error handling result for job %s: %v", j.Name, err)
			}
		})
	}
	s.cron.Start()
}

func (s *Scheduler) executeQuery(query string) (interface{}, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var result interface{}
		err = rows.Scan(&result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("No results")
}

func getEnvDefault(name string, defaultValue string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return defaultValue
}

func main() {

	scheduler, err := NewScheduler(filepath.Join(DataDir, "db.sqlite"))
	if err != nil {
		log.Fatalf("Error creating scheduler: %v", err)
	}

	scheduler.AddJob(Job{
		Name:     "Average Weight",
		Schedule: "* * * * *", // week 50 of 2022
		Query:    "SELECT AVG(weight) AS average_weight FROM registration WHERE timestamp>= '2022-12-12' AND timestamp < '2022-12-19'",
		Handler: func(result interface{}) error {
			fmt.Printf("Average weight: %v\n", result)

			err := main.PostResult(map[string]interface{}{
				"Average weight": result,
				"Start Date":     "2022-12-12",
				"End Date":       "2022-12-19",
			})

			if err != nil {
				log.Printf("Error posting result: %v", err)
				return err
			}

			return nil
		},
	})

	scheduler.AddJob(Job{
		Name:     "Daily Registrations",
		Schedule: "* * * * *",
		Query:    "SELECT COUNT(*) AS total_registrations FROM registration WHERE DATE(timestamp) = '2022-12-17'", // date - yesterday(2022-12-17)
		Handler: func(result interface{}) error {
			fmt.Printf("Total registrations: %v\n", result)

			err := main.PostResult(map[string]interface{}{
				"Total Registrations": result,
				"Date":                "2022-12-12",
			})

			if err != nil {
				log.Printf("Error posting result: %v", err)
				return err
			}

			return nil
		},
	})

	scheduler.Start()

	// Block until Ctrl+C
	c := make(chan os.Signal)
	//nolint:govet,staticcheck
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
