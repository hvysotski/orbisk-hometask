# Coding Assignment - Scheduler

Consider an imaginary system that has the responsibility of executing scheduled jobs. The purpose of a job is to gather
some metrics from a database (could be any metric) and post this to an external API.

These are the important components within this repository

* **scheduler** - Boilerplate implementation for the scheduler docker container. Currently it does nothing more than executing a query on the database and showing the result. This should be extended as part of this assignment.
* **api** - The "external" API that should receive the calculated metrics.
* **data** - A directory containing the database (SQLite) and used by the API for persistence

## Description of the assignment

The goal of this assignment is to actually implement the scheduler component:

* You should create a setup that easily expandable and allows other users to add jobs that run according to a certain schedule
* The purpose of a job is the query a data source (in this case the SQLite database) and post the result to the API
* A schedule should allow the user to specify how often and when the job runs. For example, each day at 1PM.
* It should be easy for another use to add a new job to the system
* Two jobs should already be implemented:
    * The average weight for a registration for the most recent week (note: the data is up to date until 2022-12-22, so the most recent week is week 50 of 2022)
    * The total number of registrations per day
    * Note: in normal operation these jobs would run real-time. For the coding test it is OK to just hardcode a time period and pretend that this is "yesterday" or "previous week"

### Evaluation

When evaluating your solution we will look at the following aspects:

* In general; code structure, readability, extendability and maintability
* Since it's a coding test, unit testing or any automated testing is out of scope unless you prefer to work test driven

## The database
The database is a SQLite database that contains a single table which has the following definition:

```sql
CREATE TABLE "registration" (
    "id" INTEGER NOT NULL,
    "timestamp" TEXT NOT NULL,
    "weight" REAL NOT NULL,
    PRIMARY KEY("id")
)
```

## The API
The API is very basic and allows storing and retrieving of results. The endpoints:

* GET /result - Used to get all results previously posted
* POST /result - Used to post a new (JSON) result

Example usage:

```bash
# Posting a result:
curl -X POST http://localhost:5000/result --header "Content-Type: application/json" --data '{"date": "2022-12-22", "count": 42}'

# Retrieving all results 
curl -X GET http://localhost:5000/result
```

## Getting started

This assignment requires *docker* and *docker-compose* to be installed on your system.

To run the skeleton application open a terminal, go to the directory containing the files for this assignment and run the following command:

```sh
$ docker-compose up --build
```

If all goes well, after some time, you should see the following output:

```sh
...
Recreating dev_sched_api_1       ... done
Recreating dev_sched_scheduler_1 ... done
Attaching to dev_sched_scheduler_1, dev_sched_api_1
api_1        | 
api_1        |  ┌───────────────────────────────────────────────────┐ 
api_1        |  │                   Fiber v2.40.1                   │ 
api_1        |  │               http://127.0.0.1:5000               │ 
api_1        |  │       (bound on host 0.0.0.0 and port 5000)       │ 
api_1        |  │                                                   │ 
api_1        |  │ Handlers ............. 3  Processes ........... 1 │ 
api_1        |  │ Prefork ....... Disabled  PID ................. 1 │ 
api_1        |  └───────────────────────────────────────────────────┘ 
api_1        | 
scheduler_1  | number of registrations: 1002
...
```

If you do not have this installed yet, use the following instructions:

#### Ubuntu

```sh
$ sudo apt-get update
$ sudo apt-get install docker.io docker-compose
```

#### Windows

Install the Docker Toolbox using the instructions on https://docs.docker.com/toolbox/toolbox_install_windows/

## Running from outside of the docker container

When developing it is often more convenient to run your code from outside of the Docker container. To support this, the API port (5000) will be exposed to the outside.

```sh
# Running the API from outside of the docker container can be done as follows
$ cd api
$ go run .

# Running the scheduler
$ cd scheduler
$ go run .
```

## Troubleshooting
If for some reason the Kafka data gets corrupted, and sending / receiving of messages is no longer working then clear the Kafka database by running the following from the root of the codebase:

```sh
$ docker-compose rm
```
