package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/highrow623/pseudo-ap-scripts/go/logic"
)

const (
	tricksJsonFilename = "tricks/tricks.json"
	resultsFilename    = "results/validate_json.txt"
)

func main() {
	logicTricks, err := loadLogic(tricksJsonFilename)
	check(err)

	var trickErrors []string
	trickErrors = validateTrickTags(logicTricks.EntranceTricks, logicTricks.TagHierarchy, trickErrors)
	trickErrors = validateTrickTags(logicTricks.LocationTricks, logicTricks.TagHierarchy, trickErrors)

	expandTrickTags(&logicTricks)

	trickErrors = validateRules(logicTricks.EntranceTricks, trickErrors)
	trickErrors = validateRules(logicTricks.LocationTricks, trickErrors)

	err = writeToFile(trickErrors, resultsFilename)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func loadLogic(filename string) (logic.LogicTricks, error) {
	tricksFile, err := os.Open(filename)
	if err != nil {
		return logic.LogicTricks{}, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer tricksFile.Close()

	var logicTricks logic.LogicTricks
	err = json.NewDecoder(tricksFile).Decode(&logicTricks)
	if err != nil {
		return logic.LogicTricks{}, fmt.Errorf("failed to read and unmarshal file %s: %w", filename, err)
	}

	return logicTricks, nil
}

func validateTrickTags(rules map[string][]logic.Trick, tagHierarchy map[string][]string, trickErrors []string) []string {
	for entrance, tricks := range rules {
		for i, trick := range tricks {
			for _, tag := range trick.Tags {
				if _, ok := tagHierarchy[tag]; !ok {
					errMsg := fmt.Sprintf("%s: trick %d tag %s is not in tag hierarchy", entrance, i, tag)
					trickErrors = append(trickErrors, errMsg)
				}
			}
		}
	}
	return trickErrors
}

func expandTrickTags(logicTricks *logic.LogicTricks) {
	for entrance, tricks := range logicTricks.EntranceTricks {
		for i, trick := range tricks {
			trick.Tags = expandTags(trick.Tags, logicTricks.TagHierarchy)
			logicTricks.EntranceTricks[entrance][i] = trick
		}
	}
	for location, tricks := range logicTricks.LocationTricks {
		for i, trick := range tricks {
			trick.Tags = expandTags(trick.Tags, logicTricks.TagHierarchy)
			logicTricks.LocationTricks[location][i] = trick
		}
	}
}

func expandTags(tags []string, tagHierarchy map[string][]string) []string {
	expandedTags := []string{}
	tagSet := make(map[string]bool)
	for _, tag := range tags {
		expandedTags = append(expandedTags, tag)
		tagSet[tag] = true
	}

	for i := 0; i < len(expandedTags); i++ {
		tag := expandedTags[i]
		childTags, ok := tagHierarchy[tag]
		if !ok {
			continue
		}

		for _, childTag := range childTags {
			if tagSet[childTag] {
				continue
			}
			expandedTags = append(expandedTags, childTag)
			tagSet[childTag] = true
		}
	}

	return expandedTags
}

func validateRules(rules map[string][]logic.Trick, trickErrors []string) []string {
	// definitions:
	//   if tags(t) represents the tag set of a trick t and loadout(t) represents the loadout of t
	//   tags(t1) < tags(t2) means t1's tags are a subset of t2's tags, i.e. if t2 is included, t1 must also be included
	//   loadout(t1) < loadout(t2) means t1 makes less demands on player's items than t2, i.e. if t2 is doable, t1 must also be doable

	// the purpose of this script is to validate the logic tricks json
	// things to validate:
	//   every rule has at least one trick that has no tags
	//   no tricks for a rule are exactly the same (same tags, same loadout)
	//     t1 and t2 are the same <=> tags(t1) == tags(t2) and loadout(t1) == loadout(t2)
	//   no tricks are made unnecessary by another trick in the same rule
	//     t1 is made unnecessary by t2 <=> tags(t1) > tags(t2) and loadout(t1) >= loadout(t2) or tags(t1) >= tags(t2) and loaout(t1) > loadout(t2)
	for name, tricks := range rules {
		var hasDefaultTrick bool
		for i1, trick1 := range tricks {
			if len(trick1.Tags) == 0 {
				hasDefaultTrick = true
			}
			for i2 := i1 + 1; i2 < len(tricks); i2++ {
				trick2 := tricks[i2]
				// get bools for comparing tricks
				// tags(t1) <= tags(t1)
				t1TagsLess := lessThanOrEqual(trick1.Tags, trick2.Tags)
				// tags(t2) <= tags(t1)
				t2TagsLess := lessThanOrEqual(trick2.Tags, trick1.Tags)
				// tags(t1) == tags(t2)
				tagsEqual := t1TagsLess && t2TagsLess

				t1BitRep := trick1.Loadout.BitRep()
				t2BitRep := trick2.Loadout.BitRep()
				// loadout(t1) <= loadout(t2)
				t1LoadoutLess := t1BitRep&t2BitRep == t1BitRep
				// loadout(t2) <= loadout(t1)
				t2LoadoutLess := t1BitRep&t2BitRep == t2BitRep
				// loadout(t1) == loadout(t2)
				loadoutsEqual := t1BitRep == t2BitRep

				if tagsEqual && loadoutsEqual {
					trickErrors = append(trickErrors, fmt.Sprintf("%s: tricks %d, %d are equal", name, i1, i2))
				} else if t1TagsLess && t1LoadoutLess {
					trickErrors = append(trickErrors, fmt.Sprintf("%s: trick %d is made unnecessary by trick %d", name, i2, i1))
				} else if t2TagsLess && t2LoadoutLess {
					trickErrors = append(trickErrors, fmt.Sprintf("%s: trick %d is made unnecessary by trick %d", name, i1, i2))
				}
			}
		}

		if !hasDefaultTrick {
			trickErrors = append(trickErrors, fmt.Sprintf("%s: no default tricks", name))
		}
	}

	return trickErrors
}

func writeToFile(trickErrors []string, filename string) error {
	resultsFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open results file %s: %w", filename, err)
	}
	defer resultsFile.Close()

	for _, trickError := range trickErrors {
		_, err = fmt.Fprintf(resultsFile, "%s\n", trickError)
		if err != nil {
			return fmt.Errorf("failed to write to results file %s: %w", filename, err)
		}
	}

	return nil
}

// returns true if all strings in s1 are also in s2
func lessThanOrEqual(s1, s2 []string) bool {
	s2AsSet := make(map[string]bool)
	for _, s := range s2 {
		s2AsSet[s] = true
	}
	for _, s := range s1 {
		if !s2AsSet[s] {
			return false
		}
	}
	return true
}
