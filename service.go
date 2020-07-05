package ninu

import (
	"context"
	"errors"
	"time"

	"github.com/elanq/msql"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
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

func ShowTodayTransaction() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	now := time.Now().Format("2006-01-02")
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
