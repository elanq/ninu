package ninu

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/elanq/msql"
	"github.com/jackc/pgx"
)

var (
	PostgreDB DB
)

type Row map[string]interface{}

func (r Row) Serialize() string {
	res, err := json.Marshal(r)
	if err != nil {
		return "{}"
	}

	return string(res)
}

func (r Row) String() string {
	var str string
	for k, v := range r {
		str += fmt.Sprintf("%v= %v ", k, v)
	}
	return str
}

func InitPostgre() {
	PostgreDB = NewPostgreDB()
}

const defaultDB = "ninu"

type DB interface {
	Insert(ctx context.Context, query *msql.SQLQuery) error
	FindAll(ctx context.Context, query *msql.SQLQuery) ([]interface{}, error)
}

type Postgre struct {
	conn *pgx.Conn
}

func NewPostgreDB() DB {
	var url string
	var conf pgx.ConnConfig
	var err error

	if os.Getenv("DATABASE_URL") != "" {
		url = os.Getenv("DATABASE_URL")
		conf, err = pgx.ParseURI(url)
		if err != nil {
			return nil
		}
	} else {
		conf = pgx.ConnConfig{
			Host:     os.Getenv("PG_HOST"),
			User:     os.Getenv("PG_USERNAME"),
			Password: os.Getenv("PG_PASSWORD"),
			Database: defaultDB,
		}
	}

	conn, err := pgx.Connect(conf)
	if err != nil {
		panic(err)
	}
	return &Postgre{
		conn: conn,
	}
}

func (p *Postgre) Insert(ctx context.Context, query *msql.SQLQuery) error {
	sql, args, err := query.Generate()
	if err != nil {
		return err
	}
	sql = decorateQuery(sql)

	options := &pgx.QueryExOptions{}
	_, err = p.conn.ExecEx(ctx, sql, options, args...)
	return err
}

func (p *Postgre) FindAll(ctx context.Context, query *msql.SQLQuery) ([]interface{}, error) {
	sql, args, err := query.Generate()
	if err != nil {
		return nil, err
	}
	sql = decorateQuery(sql)
	fmt.Println(args)

	options := &pgx.QueryExOptions{}
	rows, err := p.conn.QueryEx(ctx, sql, options, args...)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0)
	cols := rows.FieldDescriptions()
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		row := Row{}
		for i, col := range cols {
			row[col.Name] = values[i]
		}
		results = append(results, row)
	}
	return results, nil
}

func decorateQuery(sql string) string {
	placeholderCount := strings.Count(sql, "?")
	for i := 1; i <= placeholderCount; i += 1 {
		placeholder := fmt.Sprintf("$%v", i)
		sql = strings.Replace(sql, "?", placeholder, 1)
	}

	return sql
}
