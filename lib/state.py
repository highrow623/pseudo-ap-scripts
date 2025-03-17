from collections import Counter
from typing import List
from BaseClasses import CollectionState, MultiWorld

def build_states(multiworld: MultiWorld) -> List[CollectionState]:
    states: List[CollectionState] = []
    items = {}
    for dream_breaker in range(2):
        items["Dream Breaker"] = dream_breaker
        for strikebreak in range(2):
            items["Strikebreak"] = strikebreak
            for soul_cutter in range(2):
                items["Soul Cutter"] = soul_cutter
                for sunsetter in range(2):
                    items["Sunsetter"] = sunsetter
                    for slide in range(2):
                        items["Slide"] = slide
                        for solar_wind in range(2):
                            items["Solar Wind"] = solar_wind
                            for ascendant_light in range(2):
                                items["Ascendant Light"] = ascendant_light
                                for cling_gem in range(2):
                                    items["Cling Gem"] = cling_gem
                                    for kicks in range(5):
                                        items["Air Kick"] = kicks
                                        for small_keys in range(2):
                                            items["Small Key"] = 7 * small_keys

                                            state = CollectionState(multiworld, False)
                                            state.prog_items = {1: Counter(items)}
                                            states.append(state)
    return states
