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

			elapsed = append(elapsed, float64(time.Since(last).Milliseconds()))
			fmt.Print("*")
		}(ctx, i)
	}
	wg.Wait()

	timeElapsed := time.Since(start)

	if len(elapsed) < 1 {
		fmt.Println("\ntime elapsed(ms):", timeElapsed)
		fmt.Println("success:", success, "/", iteration)
		fmt.Println("failed:", failed, "/", iteration)
		fmt.Println("affected rows:", queryRows)
		fmt.Println("rows per second:", float32(queryRows*1000)/float32(timeElapsed.Milliseconds()))
		return
	}

	histogram, distribution := report.GetReport(elapsed)
	if err := template.PrintReport(&report.Report{
		Name:  "mysqlslap",
		Count: uint64(iteration),

		Total:   timeElapsed,
		Slowest: time.Millisecond * time.Duration(elapsed[0]),
		Fastest: time.Millisecond * time.Duration(elapsed[len(elapsed)-1]),
		Average: time.Millisecond * time.Duration(report.AverageFloat64(elapsed)),
		Rps:     float64(iteration*1000) / float64(timeElapsed.Milliseconds()),

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
