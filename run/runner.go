package run

import (
	"context"
	"database/sql"
	"fmt"
	"mysqlslap/db"
	"sync"
	"time"
)

func Run(ctx context.Context,
	client *sql.DB,
	sqls []string,
	concurrency int,
	iteration int) {

	if len(sqls) < 1 {
		fmt.Println("should input at least one sql")
		return
	}
	sqlLen := len(sqls)

	queue := make(chan int, concurrency)
	go func() {
		for i := 0; i < iteration; i++ {
			queue <- i
		}
		close(queue)
	}()

	success := 0
	failed := 0
	var queryRows int64

	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(int(iteration))
	for i := range queue {
		go func(ctx context.Context, i int) {
			defer wg.Done()
			rows, err := db.ExecQuery(ctx, client, sqls[i%sqlLen])
			fmt.Print("*")
			if err != nil {
				failed++
			} else {
				success++
			}
			queryRows += rows
		}(ctx, i)
	}
	wg.Wait()

	timeElapsed := time.Since(start)
	fmt.Println("\ntime elapsed(ms):", timeElapsed)
	fmt.Println("success:", success, "/", iteration)
	fmt.Println("failed:", failed, "/", iteration)
	fmt.Println("affected rows:", queryRows)
	fmt.Println("qps:", float32(iteration*1000)/float32(timeElapsed.Milliseconds()))
	fmt.Println("rows per second:", float32(queryRows*1000)/float32(timeElapsed.Microseconds()))
}
