package ninu

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/elanq/msql"
	"github.com/jedib0t/go-pretty/table"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	oneDay         = 24 * time.Hour
	tableName      = "transactions"
	defaultTimeout = 1 * time.Second
)

func AddTransaction(payload string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	input, err := ReadTransaction(payload)
	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02")
	query := msql.Insert(tableName,
		msql.SQLField{"date": now},
		msql.SQLField{"category": input.Category},
		msql.SQLField{"amount": input.Amount},
	)

	err = PostgreDB.Insert(ctx, query)
	select {
	case <-ctx.Done():
		return errors.New("insert timeout")
	default:
	}
	return err
}

func ShowMonthlyTransaction() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	startMonth, endMonth := monthRange()
	query := msql.Select("amount").From(tableName).Where(
		msql.SQLField{"date": formatDate(startMonth)}.Gte(),
		msql.SQLField{"date": formatDate(endMonth)}.Lt(),
	)
	fmt.Printf("showing transaction from %v to %v\n", formatDate(startMonth), formatDate(endMonth))
	results, err := PostgreDB.FindAll(ctx, query)
	if err != nil {
		return "", err
	}

	msg := "Belum ada data"
	p := message.NewPrinter(language.Indonesian)
	if len(results) > 0 {
		sum := sumResults(results)
		msg = `Jumlah transaksi bulan ini Rp %v`
		msg = p.Sprintf(msg, sum)
	}
	select {
	case <-ctx.Done():
		return "", errors.New("select timeout")
	default:
	}
	return msg, nil
}

func ShowWeeklyTransaction() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	startWeek, endWeek := weekRange()
	query := msql.Select("amount").From(tableName).Where(
		msql.SQLField{"date": formatDate(startWeek)}.Gte(),
		msql.SQLField{"date": formatDate(endWeek)}.Lt(),
	)
	fmt.Printf("showing transaction from %v to %v\n", formatDate(startWeek), formatDate(endWeek))
	results, err := PostgreDB.FindAll(ctx, query)
	if err != nil {
		return "", err
	}

	msg := "Belum ada data"
	p := message.NewPrinter(language.Indonesian)
	if len(results) > 0 {
		sum := sumResults(results)
		msg = `Jumlah transaksi minggu ini Rp %v`
		msg = p.Sprintf(msg, sum)
	}
	select {
	case <-ctx.Done():
		return "", errors.New("select timeout")
	default:
	}
	return msg, nil
}

func ShowTodayTransaction() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	now := formatDate(time.Now())
	query := msql.Select("amount").
		From("transactions").
		Where(
			msql.SQLField{"date": now},
		)

	results, err := PostgreDB.FindAll(ctx, query)
	if err != nil {
		return "", err
	}

	msg := "Belum ada data"
	p := message.NewPrinter(language.Indonesian)
	if len(results) > 0 {
		sum := sumResults(results)
		msg = `Jumlah transaksi hari ini Rp %v`
		msg = p.Sprintf(msg, sum)
	}
	select {
	case <-ctx.Done():
		return "", errors.New("select timeout")
	default:
	}
	return msg, nil
}

func ShowAllTransaction() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	first, last := monthRange()
	//	now := formatDate(time.Now())
	//	tomorrow := formatDate(time.Now().Add(24 * time.Hour))
	query := msql.Select("id", "category", "amount").
		From("transactions").
		Where(
			msql.SQLField{"date": first}.Gte(),
			msql.SQLField{"date": last}.Lt(),
		)

	results, err := PostgreDB.FindAll(ctx, query)
	if err != nil {
		return "", err
	}

	msg := buildTransactionMessage(results)
	fmt.Println(msg)
	select {
	case <-ctx.Done():
		return "", errors.New("select timeout")
	default:
	}
	return msg, nil
}

func sumResults(results []interface{}) int64 {
	sum := int64(0)
	for _, v := range results {
		row, ok := v.(Row)
		if !ok {
			continue
		}
		sum += assertInt64(row["amount"])
	}

	return sum
}

func buildTransactionMessage(results []interface{}) string {
	if len(results) > 0 {
		t := table.NewWriter()

		p := message.NewPrinter(language.Indonesian)
		t.AppendHeader(table.Row{"id", "category", "amount"})
		for _, res := range results {
			row, ok := res.(Row)
			if !ok {
				continue
			}
			amount := p.Sprintf("%v", row["amount"])
			t.AppendRow(table.Row{row["id"], row["category"], amount})
		}
		return t.Render()
	}

	return "Belum ada data"
}

func assertInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int32:
		return int64(val)
	case int64:
		return int64(val)
	}

	return 0
}

func weekRange() (time.Time, time.Time) {
	today := time.Now()
	switch today.Weekday() {
	case time.Monday:
		return today, today.Add(oneDay * 6)
	case time.Tuesday:
		return today.Add(-oneDay), today.Add(oneDay * 5)
	case time.Wednesday:
		return today.Add(-oneDay * 2), today.Add(oneDay * 4)
	case time.Thursday:
		return today.Add(-oneDay * 3), today.Add(oneDay * 3)
	case time.Friday:
		return today.Add(-oneDay * 4), today.Add(oneDay * 2)
	case time.Saturday:
		return today.Add(-oneDay * 5), today.Add(oneDay * 1)
	case time.Sunday:
		return today.Add(-oneDay * 6), today
	}
	return time.Now(), time.Now()
}

func monthRange() (time.Time, time.Time) {
	today := time.Now()
	firstMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	lastMonth := firstMonth.AddDate(0, 1, -1)

	return firstMonth, lastMonth

}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}
