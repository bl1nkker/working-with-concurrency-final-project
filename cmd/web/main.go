package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"working-with-concurrency-final-project/data"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "1234"

func main() {
	// Connect to database
	db := initDB()

	// Create sessions
	session := initSession()

	// Create loggers

	// This line creates a new logger named infoLog. It writes log messages to the standard output (os.Stdout).
	// The prefix for each log message is set to "INFO\t", followed by a tab character (\t).
	// The logging format includes the date and time of each log message (log.Ldate|log.Ltime).
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// This line creates another logger named errorLog. It also writes log messages to the standard output (os.Stdout).
	// The prefix for each log message is set to "ERROR\t", followed by a tab character (\t).
	// The logging format includes the date and time of each log message (log.Ldate|log.Ltime), as well as the file name and line number
	// where the log message was generated (log.Lshortfile).
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create some channels

	// Create wait group
	wg := sync.WaitGroup{}

	// Set up the application config
	app := Config{
		Session:  session,
		DB:       db,
		Wait:     &wg,
		Infolog:  infoLog,
		ErrorLog: errorLog,
		Models:   data.New(db),
	}

	// Set up mail

	// listen for signals
	go app.listenForShutdown()

	// Listen for web connections
	app.serve()
}

func initDB() *sql.DB {
	conn := connectToDB()

	if conn == nil {
		log.Panic("Can't connect to database")
	}
	return conn
}

func connectToDB() *sql.DB {
	counts := 0
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
		} else {
			log.Println("Connected to database")
			return connection
		}
		if counts > 10 {
			return nil
		}
		log.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		counts++
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initSession() *scs.SessionManager {
	gob.Register(data.User{})
	// set up session
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	// This sets the lifetime of each session to 24 hours. Sessions will expire and be deleted from storage after this duration of inactivity.
	session.Lifetime = 24 * time.Hour
	// This configures the session cookie to persist across browser sessions.
	// When set to true, the session cookie will remain on the user's device even after they close their browser.
	session.Cookie.Persist = true
	// This sets the SameSite attribute of the session cookie to "Lax" mode, which restricts cookies from being sent in cross-site requests
	// initiated by third-party websites. This is a security measure to prevent certain types of attacks, such as CSRF (Cross-Site Request Forgery).
	session.Cookie.SameSite = http.SameSiteLaxMode
	// This indicates that the session cookie should only be sent over HTTPS connections, ensuring that it is transmitted securely over encrypted channels.
	session.Cookie.Secure = true
	return session
}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		// This sets the maximum number of idle connections in the pool to 10.
		// Idle connections are those that are not currently in use but are kept open for future use to avoid the overhead of creating new connections.
		MaxIdle: 10,
		// This sets up a function literal (also known as an anonymous function) to establish a new connection to the Redis server.
		// The Dial function returns a redis.Conn object (a connection to the Redis server) and an error.
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return redisPool
}

func (app *Config) serve() {
	// start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.Infolog.Println("Starting web server...")
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Config) shutdown() {
	// perform any cleanup tasks
	app.Infolog.Println("Would run cleanup tasks...")

	// block until wait group is empty
	app.Wait.Wait()

	app.Infolog.Println("Closing channels and shutting down application...")
}
