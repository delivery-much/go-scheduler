<p align="center"><img src="assets/gopher.png" width="350"></p>

<h1 align="center">
  go-scheduler!
</h1>

<h3 align="center">
  your personal scheduling assistant
</h3>

- [Overview](#overview)
- [Setup](#setup)
- [Configuring the library](#configuring-the-library)
  - [Start the library providing a mongo connection](#start-the-library-providing-a-mongo-connection)
  - [Start the library providing your own Database](#start-the-library-providing-your-own-database)
  - [Providing a Logger to the library](#providing-a-logger-to-the-library)
  - [Configuring the library location](#configuring-the-library-location)
- [Scheduling jobs](#scheduling-jobs)
  - [1. Create your job function](#1-create-your-job-function)
  - [2. Define the job](#2-define-the-job)
  - [3. Schedule your job](#3-schedule-your-job)
    - [In](#in)
    - [On](#on)
    - [Every](#every)
- [Manually handling jobs](#manually-handling-jobs)
  - [Listing jobs manually](#listing-jobs-manually)
  - [Handling jobs](#handling-jobs)

## Overview
The **go-scheduler** library is a highly customizable tool that empowers developers to schedule, persist, and manage job schedules effortlessly.

With **go-scheduler**, developers can define functions and schedule them to execute at specified times.

The **go-scheduler** library is:
- `Customizable` -> The library provides numerous configuration options to better align with your requirements.
- `Persistent` -> Job schedules are persistent, ensuring that even if your application encounters issues, the schedule information remains intact.
- `Friendly` -> **go-scheduler** prioritizes a user-friendly API for managing jobs, allowing developers to create clean, scalable, and readable code!



Ex.:

```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

// create your job functions
func myJobFunction(job *scheduler.Job) (err error) {
	fmt.Println("JOB TRIGGERED!!!!")

  // your business logic goes here...
	return
}

func main() {
  // init the library
  scheduler.Init(/* you configuration goes here */)

  // define the job
  scheduler.Define("myJob", myJobFunction)

  // then you can schedule in a specific time duration...
  oneHour := time.Hour
  scheduler.In(oneHour).Do("myJob")

  // or in a specific date...
  nextWeek := time.Now().Add(time.Hour*24*7)
  scheduler.On(nextWeek).Do("myJob")

  // or you can even schedule recurrent jobs!
  scheduler.Every("monday at 16:07").Do("myJob")
}
```


## Setup
To download **go-scheduler** and add it to your project, just run:

```shell
$ go get github.com/delivery-much/go-scheduler
```

And you're good to Go!

## Configuring the library

When initiating the **go-scheduler** library, developers must provide a `Config` struct with the library necessary configuration.

```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

func main() {
  // init the library
  scheduler.Init(scheduler.Config{
    // your configuration
  })
}
```

The `Config` struct has the following values:

- `DB` -> Represents a user created struct that implements the `JobDatabase` interface.
Users should provide this value when they want to have full controll over the
actions that the library will execute in its database, or if they whant to use other database other than mongoDB.
If no user created DB is specified, the `MongoDB` value should be provided.
(See [Start the library providing your own Database](#start-the-library-providing-your-own-database) section for more)

- `MongoDB` -> Represents the configuration values that the library need to start a job DB on a mongoDB connection.
This configuration uses the original mongoDB driver to do so.
When providing this configuration, the library will have access to the user's mongoDB connection, since it will be responsible for managing jobs.
If this value is not specified, the `DB` value should be provided.
(See [Start the library providing a mongo connection](#start-the-library-providing-a-mongo-connection) section for more)

- `Logger` -> Represents a user created struct that implements the `Logger` interface.
This logger will be used by the library to log information whenever necessary.
If no logger is specified, the library will log nothing.
(See [Providing a Logger to the library](#providing-a-logger-to-the-library) section for more)

- `ProcessingRate` -> Represents the rate that the library will process jobs.
If the value is not specified, the default rate is **1 minute**.

- `Location` -> Represents the location that the library should use when generating time values.
If the value is not specified, the default location is **UTC**.
(See [Configuring the library location](#configuring-the-library-location) section for more)

- `DeleteOnDone` -> Defines if, when a job is done, the job should be deleted from the database.
If the value is not specified, the default value is **false**.

- `DeleteOnDone` -> Defines if, when a job is canceled, the job should be deleted from the database.
If the value is not specified, the default value is **false**.

> The only really required values are either the `MongoDB` or the `DB` value, the other values are optional.
> However, is highly recommended to provide the `Logger` value, so you can keep track of the job process.

### Start the library providing a mongo connection

To make developers lives easier, **go-scheduler** library allows them to provide a mongoDB connection directly when starting the library.

In this case, **go-scheduler** will be responsible for accessing and manipulating the job database, and your only worry will be to define and schedule your jobs.

The configuration should look something like this:
```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
  "go.mongodb.org/mongo-driver/mongo"
)

func main() {
  var myConnection *mongo.Client
  // ... connect to mongoDB

  // init the library
  scheduler.Init(scheduler.Config{
    MongoDB: &scheduler.MongoJobDBConfig{
	    Conn: myConnection,
	    DbName: "MyDB",
    	CollName: "MyCollection",
    },
  })
}
```


- `Conn` -> Its a pointer to your mongoDB connection! 
**go-scheduler** expects a `Client` struct from the [default mongo driver](go.mongodb.org/mongo-driver/mongo).
If no value is specified, the library initiation will fail.

- `DbName` -> Its the database name that the library should use to save jobs.
If no value is specified, the default is `"go-scheduler"`.

- `CollName` -> Its the collection name that the library should use to save jobs.
If no value is specified, the default is `"scheduler-jobs"`.

> ⚠️ **DISCLAIMER:** **go-scheduler** will **NOT** access any other database or collection than the ones that the user specified.


### Start the library providing your own Database

In case you don't want to give **go-scheduler** access to your database for safety reasons, or you don't want to work with mongo, there's no problem!!
You can also start the library with your own **custom database**.

In this case, the developer has total control over the database access, and is responsible for saving, reading and deleting jobs.

The database you provide must be a struct that implements the `JobDatabase` interface:
```go
type JobDatabase interface {
	InitJobDB() error
	ListExpiredSchedules() ([]*Job, error)
	SaveJob(j Job) error
	List(f Finder) ([]*Job, error)
	DeleteJob(j Job) error
}
```

Where:
- `InitJobDB` -> Its a function that will be called at the beginning of the library instantiation. 
It should start the job database and make it ready to read and write jobs.
In this method you can create your database indexes, run your migrations, or anything you want to do so your database is ready to manage jobs.

- `ListExpiredSchedules` -> Its a function that will be called when the library requests the jobs that should run.
It should search the database for any job schedules that are expired, and return a list of pointers to those jobs.

- `SaveJob` -> Its a function that will be called when the library needs to save a job on the database.
It should receive a job struct, and "upsert" it in the database. (If it's a new job, should insert a new job, if its an existent job, should update the existent job).

- `DeleteJob` -> Its a function that will be called when the library needs to delete a job.
It should receive a job struct, and remove it completely from the database.

- `List` -> Its a function that will be called when the developer wants to list jobs outside of the normal schedule flow.
  (See [Listing jobs manually](#listing-jobs-manually) section for more).

  This method receives a `Finder` struct, and should use the values inside the finder to list jobs in the database.
  Currently, the `Finder` allows developers to find jobs by `name`, `status` or by the extra `data` that was provided when the job was scheduled
  (See [Schedule your job](#3-schedule-your-job) section for more).
  ```go
  type Finder struct {
	  Status string
	  Name   string
	  Data   map[string]any
  }
  ```

  The developer should parse the `Finder` struct correctly into a filter that better suits the database used, and return the jobs accordingly.

It's very important to note that, to ensure the correct execution of the library, it's imperative that the `Job` struct is saved and read correctly from the database.
The job struct goes as such:
```go
type Job struct {
	ID                string
	ScheduleType      ScheduleType
	Status            ScheduleStatus
	NextRunAt         time.Time
	LastRunAt         *time.Time
	ScheduleString    string
	ScheduleLimitDate *time.Time
	Name              string
	Data              map[string]any
}
```
Every job field available should be mapped and saved correctly in the database.

When you have your database implementation ready, you can provide it when starting the library!

```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

type myDB struct { }

func (db *myDB) InitJobDB() error {
  // ... your implementation
  return nil
}

func (db *myDB) ListExpiredSchedules() ([]*Job, error) {
  // ... your implementation
  return []*Job{}, nil
}

func (db *myDB) SaveJob(j Job) error {
  // ... your implementation
  return nil
}

func (db *myDB)	DeleteJob(j Job) error {
  // ... your implementation
  return nil
}

func (db *myDB)	List(f scheduler.Finder) ([]*Job, error) {
  // ... your implementation
  return []*Job{}, nil
}

func main() {
  db := &myDB{}

  // init the library with your custom database
  scheduler.Init(scheduler.Config{
    DB: db,
  })
}
```


### Providing a Logger to the library

When dealing with jobs that are running in the background of your application, it's **extremely** important to keep track of failures during the job execution.

For that reason, the **go-scheduler** library allows developers to provide a logger when the library starts.

The logger should be a struct that implements the `Logger` interface:
```go
type Logger interface {
	Error(message string)
	Errorf(format string, a ...any)
}
```

The `Logger` interface is fairly simple, it has two methods:

- `Error` -> Which should log an error message

- `Errorf` -> Which should format the message and log it on error level

Developers should implement the `Logger` interface as they please, and provide it when the library is initiated:

```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

type myLogger struct { }

func (logger *myLogger) Error(message string) {
  // your logging goes here
}

func (logger *myLogger) Errorf(format string, a ...any)  {
  // your logging goes here
}

func main() {
  logger := &myLogger{}

  // init the library with your custom logger
  scheduler.Init(scheduler.Config{
    Logger: logger,
  })
}
```
Whenever a job execution fails, the library will use this logger struct to log an informative error message.

If no logger is specified, the library will not log any errors.

### Configuring the library location

When dealing with dates, it's important to have control over the timezone which the dates are in.

Very often the **go-scheduler** library will instantiate dates, and it's important for the library to know which timezone the developer wants to use.

Because of that, developers can provide the `Location` that the library should use when initiating the library:

```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

func main() {
  // init the library
  scheduler.Init(scheduler.Config{
    Location: "America/Sao_Paulo",
  })
}
```

The `Location` its a string field in the library configuration that developers can provide to specify which timezone should be used when instantiating dates.

If no value is specified, the timezone will be set to `UTC`.

> ⚠️ **DISCLAIMER:** It's crucial to note that the provided string must be a valid location, and the machine running the library should have the specified location available.
>
> If the location is invalid or unavailable, the `Init` function will return an error.


## Scheduling jobs

After you have [Configured the library](#configuring-the-library), you are all set to define and schedule jobs!

The following sections are a step-by-step guide for defining and handling jobs.

### 1. Create your job function

The first step is to create the logic flow that should be executed when the job triggers.

For that, developers should create job functions!

The job functions are functions that should implement the `JobFunc` contract:
```go
type JobFunc func(*Job) (err error)
```

They should receive a pointer to a job, do all the necessary logic, and then return an error if anything went wrong.

- If the job function returns an error, the library will set the job status as `FAILED`.

- If the job function returns no error, the library will set the job status as `DONE`.

Since the job function receives a pointer to the `Job` struct, developers can also alter the job struct itself if necessary, and the library will save it on the database with those changes.

> ⚠️ **DISCLAIMER:** Given the fact above, it's not recommended for the developer to alter key values of the job (like the job `ScheduleType` or `Status` for instance), since it might break the library flow.

### 2. Define the job

After you have [Created your job function](#1-create-your-job-function) you can define your job.

Defining jobs is easy, all you need to do is call the `Define` function,
passing the job name and the function that should be executed when a job with that name is triggered.

It's very important to note that, different from job schedules, job definitions are **not persisted**,
and should be defined everytime that your application is executed.

Ex.:
```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

func MyJobFunc(j *Job) (err error) {
  // your business logic goes here
  return nil
}

func main() {
  // init the library
  scheduler.Init(scheduler.Config{
    // your configuration
  })

  scheduler.Define("myJobName", MyJobFunc)
}
```

After that, everytime that a job with the name `"myJobName"` is triggered, the `MyJobFunc` will be called!


### 3. Schedule your job

After you have succesfully [Defined your job](#2-define-the-job), you can schedule it!

To schedule jobs, developers should use one of the schedule functions available, 
and then subsequently call the `Do` function so that the job can be saved correctly in the database.

```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

// create your job functions
func myJobFunction(job *scheduler.Job) (err error) {
  // your business logic goes here...
	return
}

func main() {
  // init the library
  scheduler.Init(scheduler.Config{
    // your configuration
  })

  // define the job
  scheduler.Define("myJobName", MyJobFunc)

  myExtraData := map[string]any{
    "myIntField": 42,
    "myDateField": time.Now(),
  }
  oneHour := time.Hour

  err := scheduler.In(oneHour).Do("myJob", myExtraData)
  if err != nil {
    fmt.Error("Failed to schedule job!")
  }
}
```

The `Do` function receives the job name that was [previously defined](#2-define-the-job),
and an optional map with any extra data that the user wants to save with the job.

There are three main function that the developer can use to schedule jobs using the **go-scheduler** library:

#### In

The `In` function will schedule a job to be executed **once** in the provided duration:

```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

// create your job functions
func myJobFunction(job *scheduler.Job) (err error) {
  // your business logic goes here...
	return
}

func main() {
  // init the library
  scheduler.Init(scheduler.Config{
    // your configuration
  })

  // define the job
  scheduler.Define("myJobName", MyJobFunc)

  // then you can schedule in a specific time duration...
  oneHour := time.Hour
  scheduler.In(oneHour).Do("myJobName")
}
```

#### On

The `On` function will schedule a job to be executed **once** in the provided date:

```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

// create your job functions
func myJobFunction(job *scheduler.Job) (err error) {
  // your business logic goes here...
	return
}

func main() {
  // init the library
  scheduler.Init(scheduler.Config{
    // your configuration
  })

  // define the job
  scheduler.Define("myJobName", MyJobFunc)

  nextWeek := time.Now().Add(time.Hour*24*7)
  scheduler.On(nextWeek).Do("myJobName")
}
```

> ⚠️ **DISCLAIMER:** It's crucial to note that the provided date should be in the same timezone as it was configured in the library instantiation.
(See [Configuring the library location](#configuring-the-library-location) section for more)


#### Every

The `Every` function will schedule a job to be executed **repeatedly**, given a duration string.

This duration string accepts four formats:

- A time interval string (Ex.: `"minute"`, `"2 months"`, `"6 years"`);

- A time string in HH:MM format (Ex.: `"11:27"`);

  In this scenario, the job will be scheduled to run every day at the specified hour.
  If the specified hour has already passed on the day the job is being defined, the job will be scheduled for the next day.

- A weekday string (Ex.: `"monday"`, `"friday"`);

  In this scenario, the job will be scheduled to run every week on the specified weekday, beginning at the first minute of the day (`00:01`).
  If the specified weekday has already passed in the current week during the job definition, the job will be scheduled for the next week.
  If the job is being scheduled on the specified weekday, it will be scheduled for the next week.

- A weekday and time string (Ex.: `"monday at 12:00"`, `"friday at 15:08"`).
  
  In this scenario, the job will be scheduled to run every week on the specified weekday, beginning at the specified hour.
  If the specified weekday and/or hour has already passed in the present time during the job definition, the job will be scheduled for the next week.

Everytime that the job runs successfully, it will be re-scheduled as a `PENDING` job.

If the job fails, the job status will be set as `FAILED` and **will not be re-scheduled**.

Also, developers can use the `Until` function to provide a limit date to the recurrent job.
If the job has a limit, the job will be set as `DONE` if it executes correctly and the limit date arrived.

Ex.:
```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

// create your job functions
func myJobFunction(job *scheduler.Job) (err error) {
  // your business logic goes here...
	return
}

func main() {
  // init the library
  scheduler.Init(scheduler.Config{
    // your configuration
  })

  // define the job
  scheduler.Define("myJobName", MyJobFunc)

  // schedule a recurrent job with no limit
  scheduler.Every("monday at 13:00").Do("myJobName") // this job will run forever, unless it fails or its manually canceled or deleted

  // schedule a recurrent job with a limit
  nextWeek := time.Now().Add(time.Hour*24*7)
  scheduler.Every("hour").Until(nextWeek).Do("myJobName") // this job will run until next week.
}
```

## Manually handling jobs

The **go-scheduler** library takes care of the majority of job handling for you, but there may be instances where developers want to manage specific jobs outside the regular job flow.

If you wish to cancel jobs, delete jobs, or manually change their status, we've got you covered!

### Listing jobs manually

Developers can list jobs manually using the `List` function!

Ex.:
```go
import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

func main() {
  jobs, err := scheduler.List(scheduler.Finder{
    Name:   "MyJobName",
    Status: "PENDING",
    Data: map[string]any{
      "my_data_field": "myDataValue",
    },
  })
}
```

The `List` function receives a `Finder` struct.
This struct is used to pass parameters to the listing action.

Currently, the `Finder` struct allows developers to find jobs by `name`, `status` or by the extra `data` that was provided when the job was scheduled
(See [Schedule your job](#3-schedule-your-job) section for more).
```go
type Finder struct {
  Status string
  Name   string
  Data   map[string]any
}
```

The `List` method then returns a list of pointer to jobs that where found given the filter, and an error if anything went wrong.


### Handling jobs

After you have your jobs that you want to handle, the `Job` struct provides a series of function that can help you manage it.

- `Done` -> Will set the job status as `DONE` and save it on the database;

  If the library was configured as `DeleteOnDone`, the job will be deleted instead.

- `Cancel` -> Will set the job status as `CANCELED` and save it on the database;

  If the library was configured as `DeleteOnCancel`, the job will be deleted instead.

- `Fail` -> Will set the job status as `FAILED` and save it on the database;

- `Delete` -> Will delete the job from the database;

Ex.:
```go

import (
  "time"
  "fmt"

  "github.com/delivery-much/go-scheduler"
)

func main() {
  // First you list jobs with the necessary filter...
  jobs, err := scheduler.List(scheduler.Finder{
    Name:   "MyJobName",
    Status: "PENDING",
    Data: map[string]any{
      "my_data_field": "myDataValue",
    },
  })
  if err != nil {
    fmt.Error("Failed to manage job!")
  }

  for _, j := range jobs {
    // ... then you can set it as done ...
    err = j.Done()
    if err != nil {
      fmt.Error("Failed to manage job!")
    }

    // ... or set it as failed ...
    err = j.Fail()
    if err != nil {
      fmt.Error("Failed to manage job!")
    }

    // ... or cancel it ...
    err = j.Cancel()
    if err != nil {
      fmt.Error("Failed to manage job!")
    }

    // ... or event delete it!
    err = j.Delete()
    if err != nil {
      fmt.Error("Failed to manage job!")
    }
  }
}
```

> 