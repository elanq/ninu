package ninu

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"google.golang.org/api/sheets/v4"
)

var (
	sheetID      = os.Getenv("SPREADSHEET_ID")
	sheetService *sheets.Service
)

const (
	spreadSheetName      = "ELANQIST0609_1137757232"
	lastTransactionEnd   = "sheet:transaction:end"
	lastTransactionStart = "sheet:transaction:start"
)

func sheetClient() (*sheets.Service, error) {
	if sheetService != nil {
		return sheetService, nil
	}

	return sheets.New(Client())
}

func lastStartCell() (string, string) {
	raw, err := Cache.Get(lastTransactionStart)
	if err == nil && raw != nil {
		cell := strings.Split(string(raw), "")
		return cell[0], cell[1]
	}
	return "I", "1"
}

func lastEndCell() (string, string) {
	raw, err := Cache.Get(lastTransactionEnd)
	if err == nil && raw != nil {
		cell := strings.Split(string(raw), "")
		return cell[0], cell[1]
	}
	return "K", "1"
}

func transactionSheet() string {
	startCol, startRow := lastStartCell()
	endCol, endRow := lastEndCell()

	sheetRange := fmt.Sprintf("%v!%v%v:%v%v", spreadSheetName, startCol, startRow, endCol, endRow)
	return sheetRange
}

func updateCells(updateRows int64) error {
	startCol, startRow := lastStartCell()
	endCol, endRow := lastEndCell()
	if startRow != endRow {
		return errors.New("Start and End cell row not in same number")
	}

	intStartRow, err := strconv.ParseInt(startRow, 10, 64)
	if err != nil {
		return err
	}
	intStartRow += updateRows

	intEndRow, err := strconv.ParseInt(endRow, 10, 64)
	if err != nil {
		return err
	}
	intEndRow += updateRows
	//should have error handling here
	Cache.Set(lastTransactionStart,
		[]byte(fmt.Sprintf("%v%v", startCol, intStartRow)))
	Cache.Set(lastTransactionEnd,
		[]byte(fmt.Sprintf("%v%v", endCol, intEndRow)))

	return nil
}

//ExportTransactions export transactions into excel and delete existing if process is completed
func ExportTransactions() error {
	//	client, err := sheetClient()
	//	if err != nil {
	//		return nil
	//	}
	//
	//	valueRange := in.ToValueRange()
	//	sheetRange := transactionSheet()
	//	response, err := client.Spreadsheets.Values.
	//		Update(sheetID, sheetRange, valueRange).
	//		ValueInputOption("RAW").Do()
	//	if err != nil {
	//		return err
	//	}
	//
	//	return updateCells(response.UpdatedRows)
	return nil
}
