package run

import (
	"context"
	"database/sql"
	"fmt"
	"mysqlslap/db"
	"mysqlslap/report"
	"mysqlslap/template"
	"sync"
	"time"
)

func Run(ctx context.Context,
	client *sql.DB,
	sqls []string,
	concurrency int,
	iteration int, result bool) {

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

	elapsed := make([]float64, 0, iteration)
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(int(iteration))
	for i := range queue {
		go func(ctx context.Context, i int) {
			defer wg.Done()
			last := time.Now()
			rows, err := db.ExecQuery(ctx, client, sqls[i%sqlLen])
			if err != nil {
				failed++
			} else {
				success++
			}
			queryRows += rows
			if result {
				fmt.Println(sqls[i%sqlLen], err)
			}

			elapsed = append(elapsed, time.Since(last).Seconds())
			fmt.Print("*")
		}(ctx, i)
	}
	wg.Wait()

	timeElapsed := time.Since(start)

	if len(elapsed) < 1 {
		fmt.Println("\ntime elapsed :", timeElapsed)
		fmt.Println("success:", success, "/", iteration)
		fmt.Println("failed:", failed, "/", iteration)
		fmt.Println("affected rows:", queryRows)
		fmt.Println("rows per second:", float64(queryRows)/timeElapsed.Seconds())
		return
	}

	histogram, distribution := report.GetReport(elapsed)
	if err := template.PrintReport(&report.Report{
		Name:  "mysqlslap",
		Count: uint64(iteration),

		Total:   timeElapsed,
		Fastest: time.Duration(elapsed[0] * float64(time.Second)),
		Slowest: time.Duration(elapsed[len(elapsed)-1] * float64(time.Second)),
		Average: time.Duration(report.AverageFloat64(elapsed) * float64(time.Second)),
		Rps:     float64(iteration) / timeElapsed.Seconds(),

		Histogram:           histogram,
		LatencyDistribution: distribution,

		StatusCodeDist: map[string]int{
			"success": success,
			"failed":  failed,
		},
		ErrorDist: map[string]int{},
	}); err != nil {
		fmt.Println(err)
	}
}
