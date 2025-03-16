package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	mycsv "github.com/highrow623/pseudo-ap-scripts/go/csv"
	"github.com/highrow623/pseudo-ap-scripts/go/logic"
)

const (
	tricksCsvFilename = "sheets/Pseudoregalia Tricks - Tricks.csv"
	tagsCsvFilename   = "sheets/Pseudoregalia Tricks - Tags.csv"

	tricksJsonFilename    = "tricks/tricks.json"
	tricksMinJsonFilename = "tricks/tricks.min.json"
)

func main() {
	tricksRows, err := readCsvFile(tricksCsvFilename)
	check(err)

	tagsRows, err := readCsvFile(tagsCsvFilename)
	check(err)

	logic, err := parseRows(tricksRows, tagsRows)
	check(err)

	err = outputLogic(logic, tricksJsonFilename, true)
	check(err)

	err = outputLogic(logic, tricksMinJsonFilename, false)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readCsvFile(filename string) ([]mycsv.Row, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return mycsv.RowsFromRecords(records), nil
}

func parseRows(tricksRows, tagsRows []mycsv.Row) (logic.LogicTricks, error) {
	entranceTricks := make(map[string][]logic.Trick)
	locationTricks := make(map[string][]logic.Trick)
	for i, row := range tricksRows {
		loadout, err := logic.NewLoadout(row)
		if err != nil {
			return logic.LogicTricks{}, fmt.Errorf("failed to build trick from row %d: %w", i+2, err)
		}

		trick, err := logic.NewTrick(row, loadout)
		if err != nil {
			return logic.LogicTricks{}, fmt.Errorf("failed to build trick from row %d: %w", i+2, err)
		}

		name, isLocationTrick, err := logic.GetName(row)
		if err != nil {
			return logic.LogicTricks{}, fmt.Errorf("failed to build trick from row %d: %w", i+2, err)
		}

		if isLocationTrick {
			locationTricks[name] = append(locationTricks[name], trick)
		} else {
			entranceTricks[name] = append(entranceTricks[name], trick)
		}
	}

	tagHierarchy, err := logic.ParseTags(tagsRows)
	if err != nil {
		return logic.LogicTricks{}, err
	}

	return logic.LogicTricks{
		EntranceTricks: entranceTricks,
		LocationTricks: locationTricks,
		TagHierarchy:   tagHierarchy,
	}, nil
}

func outputLogic(logic logic.LogicTricks, filename string, indent bool) error {
	tricksFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer tricksFile.Close()

	enc := json.NewEncoder(tricksFile)
	enc.SetEscapeHTML(false)
	if indent {
		enc.SetIndent("", "    ")
	}

	err = enc.Encode(logic)
	if err != nil {
		return fmt.Errorf("failed to write json to file %s: %s", filename, err)
	}

	return nil
}
