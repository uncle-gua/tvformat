package backtest

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	ErrHold  = errors.New("position hold")
	ErrEmpty = errors.New("position empty")
)

type TradeList struct {
	PositionSide string
	EntryTime    int64
	EntryPrice   float64
	CloseTime    int64
	ClosePrice   float64
}

type BackTest struct {
	Symbol       string
	PositionSide string
	EntryTime    int64
	EntryPrice   float64
	TradeList    []TradeList
}

func New(symbol string) *BackTest {
	return &BackTest{
		Symbol:       symbol,
		PositionSide: "",
		EntryTime:    0,
		EntryPrice:   math.NaN(),
		TradeList:    make([]TradeList, 0),
	}
}

func (t *BackTest) Open(positionSide string, entryTime int64, entryPrice float64) error {
	if t.PositionSide == positionSide {
		return ErrHold
	}
	if t.PositionSide != "" {
		t.Close(t.PositionSide, entryTime, entryPrice)
	}

	t.EntryTime = entryTime
	t.PositionSide = positionSide
	t.EntryPrice = entryPrice

	return nil
}

func (t *BackTest) Close(positionSide string, closeTime int64, closePrice float64) error {
	if t.PositionSide == "" {
		return ErrEmpty
	}

	t.TradeList = append(t.TradeList, TradeList{
		PositionSide: t.PositionSide,
		EntryTime:    t.EntryTime,
		EntryPrice:   t.EntryPrice,
		CloseTime:    closeTime,
		ClosePrice:   closePrice,
	})
	t.PositionSide = ""
	t.EntryTime = 0
	t.EntryPrice = math.NaN()

	return nil
}

func (t *BackTest) Output(filename string) error {
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "Symbol")
	f.SetCellValue("Sheet1", "B1", "Position Side")
	f.SetCellValue("Sheet1", "C1", "Entry Time")
	f.SetCellValue("Sheet1", "D1", "Entry Price")
	f.SetCellValue("Sheet1", "E1", "Close Time")
	f.SetCellValue("Sheet1", "F1", "Close Price")
	for i := 0; i < len(t.TradeList); i++ {
		trade := t.TradeList[i]
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), t.Symbol)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), trade.PositionSide)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), time.UnixMilli(trade.EntryTime).Format("2006-01-02 15:04:05"))
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), fmt.Sprintf("%.2f", trade.EntryPrice))
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), time.UnixMilli(trade.CloseTime).Format("2006-01-02 15:04:05"))
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), fmt.Sprintf("%.2f", trade.ClosePrice))
	}
	if err := f.SaveAs(filename); err != nil {
		return err
	}
	return f.Close()
}
