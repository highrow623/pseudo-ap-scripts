# compare logic rework against current logic after my changes
# TODO support testing progressive_breaker, progressive_slide, split_sun_greaves

from typing import Dict, List, Set
from BaseClasses import MultiWorld
from apworldrework.tricks import Loadout, state_to_loadout, loadout_to_bit_rep, bit_rep_to_summary
import lib.apworld_helper as apw
import lib.apworldreworkd_helper as apwr
from lib.difficulty import DIFFICULTIES, convert_index_to_str
from lib.state import build_states

multiworld = MultiWorld(1)
all_apw_rules = [apw.rules_from_difficulty(multiworld, difficulty) for difficulty in range(DIFFICULTIES)]
all_apwr_rules = [apwr.rules_from_difficulty(multiworld, difficulty) for difficulty in range(DIFFICULTIES)]
states = build_states(multiworld)

# entrance -> difficulty -> apw_passes -> Loadouts
entrance_errors: Dict[str, Dict[int, Dict[bool, List[Loadout]]]] = {}
# location -> difficulty -> apw_passes -> Loadouts
location_errors: Dict[str, Dict[int, Dict[bool, List[Loadout]]]] = {}

skip_locations: Set[str] = {
    "Listless Library - Sun Greaves 1",
    "Listless Library - Sun Greaves 2",
    "Listless Library - Sun Greaves 3",
}

def log_error(
        errors: Dict[str, Dict[int, Dict[bool, List[Loadout]]]],
        name: str,
        difficulty: int,
        apw_passes: bool,
        loadout: Loadout):
    if name not in errors:
        errors[name] = {difficulty:{apw_passes:[loadout]}}
        return
    if difficulty not in errors[name]:
        errors[name][difficulty] = {apw_passes:[loadout]}
        return
    if apw_passes not in errors[name][difficulty]:
        errors[name][difficulty][apw_passes] = [loadout]
        return

    # we keep the loadout list reduced as we go
    loadouts: List[Loadout] = [loadout]
    loadout_bit_rep = loadout_to_bit_rep(loadout)
    for known_loadout in errors[name][difficulty][apw_passes]:
        known_bit_rep = loadout_to_bit_rep(known_loadout)
        if loadout_bit_rep & known_bit_rep == known_bit_rep:
            # the new loadout is covered by a known one, so we can ignore the new one
            return
        if loadout_bit_rep & known_bit_rep == loadout_bit_rep:
            # the new loadout covers the known one, so we can ignore the known one
            continue
        loadouts.append(known_loadout)
    errors[name][difficulty][apw_passes] = loadouts

for difficulty in range(DIFFICULTIES):
    apw_rules = all_apw_rules[difficulty]
    apwr_rules = all_apwr_rules[difficulty]

    entrances = set(apw_rules.region_rules.keys()).union(set(apwr_rules.entrance_rules.keys()))
    for entrance in entrances:
        for state in states:
            apw_passes = entrance not in apw_rules.region_rules or apw_rules.region_rules[entrance](state)
            apwr_passes = entrance not in apwr_rules.entrance_rules or apwr_rules.entrance_rules[entrance](state)
            if apw_passes != apwr_passes:
                loadout = state_to_loadout(state, 1, 7)
                log_error(entrance_errors, entrance, difficulty, apw_passes, loadout)

    locations = set(apw_rules.location_rules.keys()).union(set(apwr_rules.location_rules.keys()))
    for location in locations:
        if location in skip_locations:
            continue
        for state in states:
            apw_passes = location not in apw_rules.location_rules or apw_rules.location_rules[location](state)
            apwr_passes = location not in apwr_rules.location_rules or apwr_rules.location_rules[location](state)
            if apw_passes != apwr_passes:
                loadout = state_to_loadout(state, 1, 7)
                log_error(location_errors, location, difficulty, apw_passes, loadout)

with open("results/compare.txt", "w") as f:
    for entrance in entrance_errors:
        f.write(f"entrance {entrance}:\n")
        for difficulty in entrance_errors[entrance]:
            f.write(f"  difficulty {convert_index_to_str(difficulty)}:\n")
            for apw_passes, loadouts in entrance_errors[entrance][difficulty].items():
                f.write(f"    {"apworld" if apw_passes else "sheet"} passes with:\n")
                for loadout in loadouts:
                    loadout_summary = bit_rep_to_summary(loadout_to_bit_rep(loadout))
                    f.write(f"      {loadout_summary}\n")

    for location in location_errors:
        f.write(f"location {location}:\n")
        for difficulty in location_errors[location]:
            f.write(f"  difficulty {convert_index_to_str(difficulty)}:\n")
            for apw_passes, loadouts in location_errors[location][difficulty].items():
                f.write(f"    {"apworld" if apw_passes else "sheet"} passes with:\n")
                for loadout in loadouts:
                    loadout_summary = bit_rep_to_summary(loadout_to_bit_rep(loadout))
                    f.write(f"      {loadout_summary}\n")
