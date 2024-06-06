package worker

import (
	"context"
	"database/sql"
	"log"
	"sync"

	"test.com/db"
	"test.com/model"
)

func DispatchWorker(worker int, db *sql.DB, data chan []model.IMDBMovie, wg *sync.WaitGroup) {
	for i := 1; i <= worker; i++ {
		go func(workerIdx int, db *sql.DB, data chan []model.IMDBMovie, wg *sync.WaitGroup) {
			for d := range data {
				doTheJob(workerIdx, db, d, wg)
				wg.Done()
			}

		}(i, db, data, wg)
	}
}

func doTheJob(workerIndex int, sq *sql.DB, data []model.IMDBMovie, wg *sync.WaitGroup) {
	wg.Add(1)
	for _, d := range data {
		conn, err := sq.Conn(context.Background())

		if err != nil {
			log.Println(err)
			break
		}

		raw, fields := db.CreateRawInsert(d)
		_, err = conn.ExecContext(context.Background(), raw, fields...)

		if err != nil {
			log.Println(err)
			break
		}

		err = conn.Close()

		if err != nil {
			log.Println(err)
			break
		}

		log.Println("Worker", workerIndex, "Berhasil insert data")
	}
}
