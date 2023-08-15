package utils

import (
	"context"
	"fmt"
	"github.com/dolphindb/api-go/api"
	m "github.com/dolphindb/api-go/model"
	"github.com/zhengow/Z-Quant/model"
	"strings"
	"time"
)

func LoadData(symbol string, cols []string) []model.Price {
	host := "localhost:8848"
	db, err := api.NewSimpleDolphinDBClient(context.TODO(), host, "admin", "123456")
	dbPath := "dfs://bar"
	tbName := "min"
	selectCols := strings.Join(cols, ", ")
	if !strings.Contains(selectCols, "datetime") {
		selectCols += ", datetime"
		cols = append(cols, "datetime")
	}
	sql := fmt.Sprintf("select %s from loadTable('%s','%s') where instrument='%s'", selectCols, dbPath, tbName, symbol)
	tmp, err := db.RunScript(sql)
	if err != nil {
		panic(err)
	}
	tb := tmp.(*m.Table)
	var prices []model.Price
	values := make(map[string][]interface{})
	for _, col := range cols {
		values[col] = tb.GetColumnByName(col).Data.Value()
	}
	for i := 0; i < tb.Rows(); i++ {
		p := model.Price{}
		for _, col := range cols {
			val := values[col][i]
			switch col {
			case "datetime":
				p.Datetime = val.(time.Time)
			case "close":
				p.Close = val.(float64)
			case "open":
				p.Open = val.(float64)
			case "high":
				p.High = val.(float64)
			case "low":
				p.Low = val.(float64)
			case "volume":
				p.Volume = val.(float64)
			}
		}
		prices = append(prices, p)
	}
	return prices
}
