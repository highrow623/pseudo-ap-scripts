from BaseClasses import MultiWorld
import Options
from apworldrework import PseudoregaliaWorld
from apworldrework.rules import PseudoregaliaRules
import apworldrework.options as options
from lib.difficulty import NORMAL_OBSCURE_INDEX, HARD_INDEX, HARD_OBSCURE_INDEX, EXPERT_INDEX, LUNATIC_INDEX

def options_from_difficulty(difficulty: int) -> options.PseudoregaliaOptions:
    tags = []
    if difficulty == NORMAL_OBSCURE_INDEX:
        tags = ["obscure"]
    elif difficulty == HARD_INDEX:
        tags = ["hard"]
    elif difficulty == HARD_OBSCURE_INDEX:
        tags = ["hard", "obscure"]
    elif difficulty == EXPERT_INDEX:
        tags = ["expert"]
    elif difficulty == LUNATIC_INDEX:
        tags = ["lunatic"]

    return options.PseudoregaliaOptions(
        Options.ProgressionBalancing(50),
        Options.Accessibility(0),
        Options.LocalItems([]),
        Options.NonLocalItems([]),
        Options.StartInventory({}),
        Options.StartHints([]),
        Options.StartLocationHints([]),
        Options.ExcludeLocations([]),
        Options.PriorityLocations([]),
        Options.ItemLinks([]),
        options.TrickTags(tags),
        options.IncludeTrickIDs([]),
        options.ExcludeTrickIDs([]),
        options.ProgressiveBreaker(0),
        options.ProgressiveSlide(0),
        options.SplitSunGreaves(0),
        Options.DeathLink(0),
    )

def rules_from_difficulty(multiworld: MultiWorld, difficulty: int) -> PseudoregaliaRules:
    world = PseudoregaliaWorld(multiworld, 1)
    world.options = options_from_difficulty(difficulty)
    return PseudoregaliaRules(world)
 