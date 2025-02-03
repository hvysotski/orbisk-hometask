package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

      "github.com/robfig/cron/v3"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DefaultDatabaseName = "db.sqlite"
)


var (
	DataDir = getEnvDefault("DATA_DIR", "../data")
)

type Job struct {
    Name string,
    Schedule string,
    Query string,
    Handler func(result interface{}) error
}

type Scheduler struct {
    cron *cron.Cron
    db *sql.DB
    jobs []Job
}

func NewScheduler(dbPath string) (*Scheduler, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    return &Scheduler{
        cron: cron.New(),
        db: db,
    }, nil
}

func (s* Scheduler) AddJob(job Job) {
    s.jobs = append(s.jobs, job)
}

func (s* Scheduler) Start() {
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
        var result interface {}
        err = rows.Scan(&result)
        if err != nil {
            return nil, err
        }
        return result, nil
    }
    return nil, fmt.Errorf("No results")
}

// getEnvDefault return the value of the environment variable specified by name, or the defaultValue if not set
func getEnvDefault(name string, defaultValue string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return defaultValue
}

func main() {

// 	db, err := sql.Open("sqlite3", filepath.Join(DataDir, DefaultDatabaseName))
//
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	defer db.Close()
//
// 	var count int
// 	err = db.QueryRow("select count(*) from registration").Scan(&count)
//
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	fmt.Printf("number of registrations: %d\n", count)

    scheduler, err := NewScheduler(filepath.Join(DataDir, DefaultDatabaseName)
    if err != nil {
        log.Fatalf("Error creating scheduler: %v", err)
    }

    scheduler.AddJob(Job{
        Name:      "Average Weight",
        Schedule:  "0 13 * * *", // daily at 1PM
        Query:     "SELECT AVG(weight) FROM registrations WHERE week = 50 AND year = 2022",
        Handler: func(result interface{}) error {
            fmt.Printf("Average weight: %v\n", result)

            // Here you would post the result to the API

            return nil
        },
    })

    scheduler.AddJob(Job{
        Name:     "Daily Registrations",
        Schedule: "0 13 * * *", // Run daily at 1 PM
        Query:    "SELECT COUNT (*) FROM registrations WHERE date = '2022-12-22'", // date==yesterday
        Handler: func(result interface{}) error {
            fmt.Printf("Total registrations: %v\n", result)

            // Here you would post the result to the API
            return nil
        },
    })

    scheduler.Start()

    select{}


	// Block until Ctrl+C
	c := make(chan os.Signal)
	//nolint:govet,staticcheck
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
