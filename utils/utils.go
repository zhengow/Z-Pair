package utils

import (
	"context"
	"fmt"
	"github.com/dolphindb/api-go/api"
	m "github.com/dolphindb/api-go/model"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	gocharts "github.com/vicanso/go-charts/v2"
	"github.com/zhengow/Z-Quant/model"
	"math"
	"os"
	"path/filepath"
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
	prices := make([]model.Price, tb.Rows())
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
		prices[i] = p
	}
	return prices
}

func Draw(factors []model.Factor) {
	datetimes := make([]time.Time, 0)
	vals := make([]opts.LineData, 0)
	for _, _factor := range factors {
		datetimes = append(datetimes, _factor.Datetime)
		vals = append(vals, opts.LineData{
			Value: _factor.Val * 100000,
		})
	}
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width: "100%",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Line",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	// Put data into instance
	line.SetXAxis(datetimes).
		AddSeries("factor", vals)
	f, _ := os.Create("line.html")
	line.Render(f)
}

func DrawSimple(factors model.BaseSeries, name string) {
	datetimes := make([]string, 0)
	vals := make([]float64, 0)
	yMax := float64(-1)
	yMin := float64(1)
	for _, _factor := range factors {
		datetimes = append(datetimes, _factor.Datetime.Format("2006/02/01"))
		vals = append(vals, _factor.Val)
		yMin = math.Min(yMin, _factor.Val)
		yMax = math.Max(yMax, _factor.Val)
	}
	values := [][]float64{vals}
	p, err := gocharts.LineRender(
		values,
		gocharts.XAxisDataOptionFunc(datetimes),
		gocharts.TitleTextOptionFunc("Line"),
		gocharts.YAxisOptionFunc(gocharts.YAxisOption{
			Min: &yMin,
			Max: &yMax,
		}),
		func(goOpt *gocharts.ChartOption) {
			goOpt.Legend.Padding = gocharts.Box{
				Top:    5,
				Bottom: 10,
			}
			goOpt.SymbolShow = gocharts.FalseFlag()
			goOpt.LineStrokeWidth = 1
			goOpt.ValueFormatter = func(f float64) string {
				return fmt.Sprintf("%.5f", f)
			}
			goOpt.Width = 1200
			goOpt.Height = 800
		},
	)

	if err != nil {
		panic(err)
	}
	buf, err := p.Bytes()
	if err != nil {
		panic(err)
	}
	writeFile := func(buf []byte) error {
		tmpPath := "./tmp"
		err := os.MkdirAll(tmpPath, 0700)
		if err != nil {
			return err
		}
		file := filepath.Join(tmpPath, fmt.Sprintf("%s.png", name))
		err = os.WriteFile(file, buf, 0600)
		if err != nil {
			return err
		}
		return nil
	}
	err = writeFile(buf)
	if err != nil {
		panic(err)
	}
}

func Divide(a, b float64) float64 {
	result := math.Floor(a / b)
	if result >= 0 {
		return result
	} else {
		return result + 1
	}
}
