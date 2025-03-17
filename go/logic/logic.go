package logic

import (
	"fmt"

	"github.com/highrow623/pseudo-ap-scripts/go/csv"
)

const (
	headerDreamBreaker   = "Dream Breaker"
	headerStrikebreak    = "Strikebreak"
	headerSoulCutter     = "Soul Cutter"
	headerSunsetter      = "Sunsetter"
	headerSlide          = "Slide"
	headerSolarWind      = "Solar Wind"
	headerAscendantLight = "Ascendant Light"
	headerClings         = "Clings"
	headerKicks          = "Kicks"
	headerSmallKeys      = "Small Keys"

	headerTrickID = "Trick ID"
	headerTags    = "Tags"

	headerLocation        = "Location"
	headerRegion          = "Region"
	headerConnectedRegion = "Connected Region"

	headerTag       = "Tag"
	headerChildTags = "Child Tags"
)

type LogicTricks struct {
	EntranceTricks map[string][]Trick  `json:"entrance_tricks"`
	LocationTricks map[string][]Trick  `json:"location_tricks"`
	TagHierarchy   map[string][]string `json:"tag_hierarchy"`
}

type Trick struct {
	ID      string   `json:"id"`
	Loadout Loadout  `json:"loadout"`
	Tags    []string `json:"tags,omitempty"`
}

type Loadout struct {
	DreamBreaker   bool `json:"dream_breaker,omitempty"`
	Strikebreak    bool `json:"strikebreak,omitempty"`
	SoulCutter     bool `json:"soul_cutter,omitempty"`
	Sunsetter      bool `json:"sunsetter,omitempty"`
	Slide          bool `json:"slide,omitempty"`
	SolarWind      bool `json:"solar_wind,omitempty"`
	AscendantLight bool `json:"ascendant_light,omitempty"`
	Clings         int  `json:"clings,omitempty"`
	Kicks          int  `json:"kicks,omitempty"`
	SmallKeys      bool `json:"small_keys,omitempty"`
}

func NewTrick(row csv.Row, loadout Loadout) (Trick, error) {
	trickID, ok := row.GetString(headerTrickID)
	if !ok {
		return Trick{}, headerError(headerTrickID)
	}

	tags, ok := row.GetStringSlice(headerTags, ", ")
	if !ok {
		return Trick{}, headerError(headerTags)
	}

	return Trick{
		ID:      trickID,
		Loadout: loadout,
		Tags:    tags,
	}, nil
}

func NewLoadout(row csv.Row) (Loadout, error) {
	dreamBreaker, ok := row.GetBool(headerDreamBreaker)
	if !ok {
		return Loadout{}, headerError(headerDreamBreaker)
	}
	strikebreak, ok := row.GetBool(headerStrikebreak)
	if !ok {
		return Loadout{}, headerError(headerStrikebreak)
	}
	soulCutter, ok := row.GetBool(headerSoulCutter)
	if !ok {
		return Loadout{}, headerError(headerSoulCutter)
	}
	sunsetter, ok := row.GetBool(headerSunsetter)
	if !ok {
		return Loadout{}, headerError(headerSunsetter)
	}
	slide, ok := row.GetBool(headerSlide)
	if !ok {
		return Loadout{}, headerError(headerSlide)
	}
	solarWind, ok := row.GetBool(headerSolarWind)
	if !ok {
		return Loadout{}, headerError(headerSolarWind)
	}
	ascendantLight, ok := row.GetBool(headerAscendantLight)
	if !ok {
		return Loadout{}, headerError(headerAscendantLight)
	}
	clings, ok := row.GetInt(headerClings)
	if !ok {
		return Loadout{}, headerError(headerClings)
	}
	kicks, ok := row.GetInt(headerKicks)
	if !ok {
		return Loadout{}, headerError(headerKicks)
	}
	smallKeys, ok := row.GetBool(headerSmallKeys)
	if !ok {
		return Loadout{}, headerError(headerSmallKeys)
	}

	return Loadout{
		DreamBreaker:   dreamBreaker,
		Strikebreak:    strikebreak,
		SoulCutter:     soulCutter,
		Sunsetter:      sunsetter,
		Slide:          slide,
		SolarWind:      solarWind,
		AscendantLight: ascendantLight,
		Clings:         clings,
		Kicks:          kicks,
		SmallKeys:      smallKeys,
	}, nil
}

func (loadout Loadout) BitRep() int {
	var bitRep int
	mask := 1

	markIfTrue := func(condition bool) {
		if condition {
			bitRep = bitRep | mask
		}
		mask = mask << 1
	}

	markIfTrue(loadout.DreamBreaker)
	markIfTrue(loadout.Strikebreak)
	markIfTrue(loadout.SoulCutter)
	markIfTrue(loadout.Sunsetter)
	markIfTrue(loadout.Slide)
	markIfTrue(loadout.SolarWind)
	markIfTrue(loadout.AscendantLight)
	for i := 1; i <= 6; i++ {
		markIfTrue(loadout.Clings >= i)
	}
	for i := 1; i <= 4; i++ {
		markIfTrue(loadout.Kicks >= i)
	}
	markIfTrue(loadout.SmallKeys)

	return bitRep
}

func GetName(row csv.Row) (name string, isLocation bool, err error) {
	var ok bool
	name, ok = row.GetString(headerLocation)
	if !ok {
		return "", false, headerError(headerLocation)
	}

	if name != "" {
		isLocation = true
		return
	}

	region, ok := row.GetString(headerRegion)
	if !ok {
		return "", false, headerError(headerRegion)
	}

	connectedRegion, ok := row.GetString(headerConnectedRegion)
	if !ok {
		return "", false, headerError(headerConnectedRegion)
	}

	name = fmt.Sprintf("%s -> %s", region, connectedRegion)
	return
}

func ParseTags(rows []csv.Row) (map[string][]string, error) {
	tagHierarchy := make(map[string][]string)
	for _, row := range rows {
		tag, ok := row.GetString(headerTag)
		if !ok {
			return nil, headerError(headerTag)
		}

		if tag == "" {
			continue
		}

		childTags, ok := row.GetStringSlice(headerChildTags, ", ")
		if !ok {
			return nil, headerError(headerChildTags)
		}

		tagHierarchy[tag] = childTags
	}
	return tagHierarchy, nil
}

func headerError(header string) error {
	return fmt.Errorf("failed to get header %s", header)
}
