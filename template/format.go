package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"mysqlslap/report"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	barChar = "∎"
)

func jsonify(v interface{}, pretty bool) string {
	d, _ := json.Marshal(v)
	if !pretty {
		return string(d)
	}

	var out bytes.Buffer
	err := json.Indent(&out, d, "", "  ")
	if err != nil {
		return string(d)
	}

	return out.String()
}

func formatNanoUnit(d time.Duration) string {
	v := d.Nanoseconds()
	if v < 10000 {
		return fmt.Sprintf("%+v ns", v)
	}

	valMs := float64(v) / 1000000.0
	if valMs < 1000 {
		return fmt.Sprintf("%4.2f ms", valMs)
	}

	return fmt.Sprintf("%4.2f s", float64(valMs)/1000.0)
}

func formatMilli(duration float64) string {
	return fmt.Sprintf("%4.2f", duration*1000)
}

func formatDate(d time.Time) string {
	return d.Format("Mon Jan _2 2006 @ 15:04:05")
}

func formatSeconds(duration float64) string {
	return fmt.Sprintf("%4.2f", duration)
}

func formatPercent(num int, total uint64) string {
	p := float64(num) / float64(total)
	return fmt.Sprintf("%.2f", p*100)
}

func histogram(buckets []report.Bucket) string {
	maxMark := 0.0
	maxCount := 0
	for _, b := range buckets {
		if v := b.Mark; v > maxMark {
			maxMark = v
		}
		if v := b.Count; v > maxCount {
			maxCount = v
		}
	}

	formatMark := func(mark float64) string {
		return fmt.Sprintf("%.3f", mark*1000)
	}
	formatCount := func(count int) string {
		return fmt.Sprintf("%v", count)
	}

	maxMarkLen := len(formatMark(maxMark))
	maxCountLen := len(formatCount(maxCount))
	res := new(bytes.Buffer)
	for i := 0; i < len(buckets); i++ {
		// Normalize bar lengths.
		var barLen int
		if maxCount > 0 {
			barLen = (buckets[i].Count*40 + maxCount/2) / maxCount
		}
		markStr := formatMark(buckets[i].Mark)
		countStr := formatCount(buckets[i].Count)
		res.WriteString(fmt.Sprintf(
			"  %s%s [%v]%s |%v\n",
			markStr,
			strings.Repeat(" ", maxMarkLen-len(markStr)),
			countStr,
			strings.Repeat(" ", maxCountLen-len(countStr)),
			strings.Repeat(barChar, barLen),
		))
	}

	return res.String()
}

func formatMarkMs(m float64) string {
	m = m * 1000.0

	if m < 1 {
		return fmt.Sprintf("'%4.4f ms'", m)
	}

	return fmt.Sprintf("'%4.2f ms'", m)
}

func formatStatusCode(statusCodeDist map[string]int) string {
	padding := 3
	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 0, 0, padding, ' ', 0)
	for status, count := range statusCodeDist {
		// bytes.Buffer can be assumed to not fail on write
		_, _ = fmt.Fprintf(w, "  [%+s]\t%+v responses\t\n", status, count)
	}
	// bytes.Buffer can be assumed to not fail on write
	_ = w.Flush()
	return buf.String()
}

func formatErrorDist(errDist map[string]int) string {
	padding := 3
	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 0, 0, padding, ' ', 0)
	for status, count := range errDist {
		// bytes.Buffer can be assumed to not fail on write
		_, _ = fmt.Fprintf(w, "  [%+v]\t%+s\t\n", count, status)
	}
	// bytes.Buffer can be assumed to not fail on write
	_ = w.Flush()
	return buf.String()
}

var tmplFuncMap = template.FuncMap{
	"formatMilli":      formatMilli,
	"formatSeconds":    formatSeconds,
	"histogram":        histogram,
	"jsonify":          jsonify,
	"formatMark":       formatMarkMs,
	"formatPercent":    formatPercent,
	"formatStatusCode": formatStatusCode,
	"formatErrorDist":  formatErrorDist,
	"formatDate":       formatDate,
	"formatNanoUnit":   formatNanoUnit,
}

func PrintReport(rp *report.Report) error {
	buf := &bytes.Buffer{}
	templ := template.Must(template.New("tmpl").Funcs(tmplFuncMap).Parse(defaultTmpl))
	if err := templ.Execute(buf, *rp); err != nil {
		return err
	}

	return rp.Print(buf.String())
}
