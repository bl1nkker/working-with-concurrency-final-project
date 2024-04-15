package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
)

type Config struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	Infolog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
}
