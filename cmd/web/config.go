package main

import (
	"database/sql"
	"log"
	"sync"
	"working-with-concurrency-final-project/data"

	"github.com/alexedwards/scs/v2"
)

type Config struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	Infolog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Models   data.Models
}
