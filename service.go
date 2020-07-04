package ninu

import (
	"context"
	"errors"
	"time"

	"github.com/elanq/msql"
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
