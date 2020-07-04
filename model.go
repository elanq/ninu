package ninu

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

var (
	ErrParamsNotEnough = errors.New("Parameter is not enough for this request")
)

type InputTransaction struct {
	Date     time.Time
	Category string
	Amount   int64
}

//ReadTransaction reads message from telegram with following format
// /add Makanan 20000
// /add is the telegram command which won't included in the input
// Makanan is the category of the transaction
// 20000 is the amount of the transaction
func ReadTransaction(message string) (*InputTransaction, error) {
	values := strings.Split(message, " ")
	if len(values) == 2 {
		category := values[0]
		amount, err := strconv.ParseInt(values[1], 10, 64)
		if err != nil {
			return nil, err

		}
		return &InputTransaction{
			Date:     time.Now(),
			Category: category,
			Amount:   amount,
		}, nil
	}
	return nil, ErrParamsNotEnough
}

func (i *InputTransaction) ToValueRange() *sheets.ValueRange {
	var values []interface{}
	//IDK. there's no elegant way to do this
	values = append(values, i.Date.Format("02/01/2006"))
	values = append(values, i.Category)
	values = append(values, i.Amount)

	return &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values: [][]interface{}{
			values,
		},
	}
}
