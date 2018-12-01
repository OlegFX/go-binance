package binance

import (
	"github.com/buger/jsonparser"
	"github.com/shopspring/decimal"
)

var pathInWsPartialBookDepth = [][]string{
	{"lastUpdateId"},
	{"bids"},
	{"asks"},
}

func parseWsPartialBookDepthEvent(data []byte, levels int) *WsPartialBookDepthEvent  {
	bookDepth := &WsPartialBookDepthEvent{}

	jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error){
		switch idx {
		case 0: // lastUpdateId
			bookDepth.LastUpdateId, _ = jsonparser.GetInt(value)
		case 1: // bids
			bookDepth.Bids = parseBidsAsks(value, levels)
		case 2: // asks
			bookDepth.Asks = parseBidsAsks(value, levels)
		}
	}, pathInWsPartialBookDepth...)

	return bookDepth
}

func parseBidsAsks(data []byte, levels int) [][2]decimal.Decimal {
	asksBids := make([][2]decimal.Decimal, levels, levels)
	i := 0
	var arr [2]decimal.Decimal
	var b int

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if i > levels - 1 {
			return
		}
		arr = [2]decimal.Decimal{}
		b = 0

		jsonparser.ArrayEach(value, func(d []byte, dataType jsonparser.ValueType, offset int, err error) {
			if b > 1 {
				return
			}
			arr[b], _ = decimal.NewFromString(string(d))
			b++
		});

		asksBids[i] = arr
		i++
	});

	return asksBids
}
