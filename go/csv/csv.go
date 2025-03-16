package csv

import (
	"strconv"
	"strings"
)

type Row struct {
	raw           []string
	headerIndexes map[string]int
}

func (row Row) GetString(header string) (string, bool) {
	headerIndex, ok := row.headerIndexes[header]
	if !ok {
		return "", false
	}

	if headerIndex >= len(row.raw) {
		return "", false
	}

	return row.raw[headerIndex], true
}

func (row Row) GetStringSlice(header string, sep string) ([]string, bool) {
	rawCell, ok := row.GetString(header)
	if !ok {
		return nil, ok
	}

	if rawCell == "" {
		return []string{}, true
	}

	return strings.Split(rawCell, sep), true
}

func (row Row) GetBool(header string) (bool, bool) {
	rawCell, ok := row.GetString(header)
	if !ok {
		return false, ok
	}

	return rawCell == "TRUE", true
}

func (row Row) GetInt(header string) (int, bool) {
	rawCell, ok := row.GetString(header)
	if !ok {
		return 0, ok
	}

	if rawCell == "" {
		return 0, true
	}

	num, err := strconv.ParseInt(rawCell, 10, 64)
	if err != nil {
		return 0, false
	}

	return int(num), true
}

func RowsFromRecords(records [][]string) []Row {
	if len(records) == 0 {
		return nil
	}

	headerIndexes := make(map[string]int)
	for headerIndex, header := range records[0] {
		headerIndexes[header] = headerIndex
	}

	var rows []Row
	for i := 1; i < len(records); i++ {
		rawRow := records[i]
		rows = append(rows, Row{raw: rawRow, headerIndexes: headerIndexes})
	}
	return rows
}
