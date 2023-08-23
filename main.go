package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"tvformat/backtest"
)

var (
	help    bool
	inFile  string
	outFile string
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&inFile, "i", "", "input file")
	flag.StringVar(&outFile, "o", "", "output file")
}

func main() {
	flag.Parse()

	if help || inFile == "" || outFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	lines, err := readLines(inFile)
	if err != nil {
		log.Fatal(err)
	}
	backtest := backtest.New("ETHUSDT")
	for i := len(lines) - 1; i >= 0; i-- {
		cells := strings.Split(lines[i], ",")
		if len(cells) == 14 && cells[13] != "" {
			switch cells[1] {
			case "多头进场":
				entryTime, err := time.ParseInLocation("2006-01-02 15:04", cells[3], time.Now().Location())
				if err != nil {
					log.Fatal(err)
				}
				entryPrice, err := strconv.ParseFloat(cells[4], 64)
				if err != nil {
					log.Fatal(err)
				}
				backtest.Open("LONG", entryTime.UnixMilli(), entryPrice)
			case "多头出场":
				entryTime, err := time.ParseInLocation("2006-01-02 15:04", cells[3], time.Now().Location())
				if err != nil {
					log.Fatal(err)
				}
				closePrice, err := strconv.ParseFloat(cells[4], 64)
				if err != nil {
					log.Fatal(err)
				}
				backtest.Close("LONG", entryTime.UnixMilli(), closePrice)
			case "空头进场":
				entryTime, err := time.ParseInLocation("2006-01-02 15:04", cells[3], time.Now().Location())
				if err != nil {
					log.Fatal(err)
				}
				entryPrice, err := strconv.ParseFloat(cells[4], 64)
				if err != nil {
					log.Fatal(err)
				}
				backtest.Open("SHORT", entryTime.UnixMilli(), entryPrice)
			case "空头出场":
				entryTime, err := time.ParseInLocation("2006-01-02 15:04", cells[3], time.Now().Location())
				if err != nil {
					log.Fatal(err)
				}
				closePrice, err := strconv.ParseFloat(cells[4], 64)
				if err != nil {
					log.Fatal(err)
				}
				backtest.Close("SHORT", entryTime.UnixMilli(), closePrice)
			}
		}
	}
	backtest.Output(outFile)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
