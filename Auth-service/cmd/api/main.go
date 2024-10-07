package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webpost = "80"
var counts int64
type Config struct {
	DB *sql.DB
	Models data.Models
}

func main() {

	log.Print("starting the authentication service ")

	// todo conect to db

	conn := ConnectToDb()

	if conn == nil {

		log.Panic("can't connect to postgres ")
	}
	// set up config

	app := Config{
		DB: conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webpost),
		Handler: app.routes(),
	}
	err :=srv.ListenAndServe()
	if err != nil{
		log.Panic(err)
	}
}

func OpneDB(dsn string) (*sql.DB,error) {

	db,	err:=sql.Open("pgx",dsn)
	if err != nil {
		return nil ,err 
	}
	err =db.Ping()
	if err != nil{
		return nil,err
	}
	return db ,nil
}

func ConnectToDb()(*sql.DB){

	dsn := os.Getenv("DSN")
log.Println("124",dsn)
	for{
		connection ,err := OpneDB(dsn)

		if err != nil{
			log.Println("Postgres not yet ready..." ,dsn)
			counts++
			}else {

			log.Println("Conneted to postgres...")
			return connection
		}

		if counts >10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}