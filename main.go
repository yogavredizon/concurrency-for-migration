package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"test.com/db"
	"test.com/model"
	"test.com/worker"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "root"
	dbname   = "test"
)

func main() {
	start := time.Now()

	// Channel for retreiving data from mongoDB
	ch := make(chan []model.IMDBMovie)

	// Initialize synchronization
	wg := new(sync.WaitGroup)

	// Create connection to Database MonggoDB
	ctx := context.Background()

	m := db.NewMongo("IMDB")
	client, _ := m.Connect(ctx, "mongodb://localhost:27017")

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			return
		}
	}()

	// Create Connection to Database Postgre SQL
	var dsn string = fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		host, port, user, password, dbname,
	)

	conn := db.SqlConnect(dsn, "postgres")

	// Create tabel in database
	db.CreateTable(conn, model.IMDBMovie{})

	// Runnning 5 goroutines
	go worker.DispatchWorker(5, conn, ch, wg)

	// fetch data from database MongoDB
	m.SelectWithLimit(ctx, client, ch)

	// Wait until all process in goroutines are done
	wg.Wait()

	duration := time.Since(start)

	log.Println("Data inserted in", duration)

}

func WithOutConcurrent(s []model.IMDBMovie) {

	var dsn string = fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		host, port, user, password, dbname,
	)

	sql := db.SqlConnect(dsn, "postgres")
	db.CreateTable(sql, model.IMDBMovie{})

	for i, v := range s {
		raw, fields := db.CreateRawInsert(v)

		_, err := sql.Exec(raw, fields...)
		if err != nil {
			return
		}

		log.Println(i+1, "Success insert data")
	}

}
