from BaseClasses import MultiWorld
import Options
from apworld import PseudoregaliaWorld
from apworld.rules import PseudoregaliaRulesHelpers
from apworld.rules_normal import PseudoregaliaNormalRules
from apworld.rules_hard import PseudoregaliaHardRules
from apworld.rules_expert import PseudoregaliaExpertRules
from apworld.rules_lunatic import PseudoregaliaLunaticRules
import apworld.options as options
from lib.difficulty import NORMAL_INDEX, NORMAL_OBSCURE_INDEX, HARD_INDEX, HARD_OBSCURE_INDEX, EXPERT_INDEX, LUNATIC_INDEX

def options_from_difficulty(difficulty: int) -> options.PseudoregaliaOptions:
    logic_level = 1
    if difficulty in {HARD_INDEX, HARD_OBSCURE_INDEX}:
        logic_level = 2
    elif difficulty == EXPERT_INDEX:
        logic_level = 3
    elif difficulty == LUNATIC_INDEX:
        logic_level = 4

    obscure = 1
    if difficulty in {NORMAL_INDEX, HARD_INDEX}:
        obscure = 0

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
        options.LogicLevel(logic_level),
        options.ObscureLogic(obscure),
        options.ProgressiveBreaker(0),
        options.ProgressiveSlide(0),
        options.SplitSunGreaves(0),
        Options.DeathLink(0),
    )

def rules_from_difficulty(multiworld: MultiWorld, difficulty: int) -> PseudoregaliaRulesHelpers:
    world = PseudoregaliaWorld(multiworld, 1)
    world.options = options_from_difficulty(difficulty)
    if difficulty in {NORMAL_INDEX, NORMAL_OBSCURE_INDEX}:
        return PseudoregaliaNormalRules(world)
    if difficulty in {HARD_INDEX, HARD_OBSCURE_INDEX}:
        return PseudoregaliaHardRules(world)
    if difficulty == EXPERT_INDEX:
        return PseudoregaliaExpertRules(world)
    if difficulty == LUNATIC_INDEX:
        return PseudoregaliaLunaticRules(world)
