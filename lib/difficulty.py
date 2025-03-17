NORMAL_INDEX = 0
NORMAL_OBSCURE_INDEX = 1
HARD_INDEX = 2
HARD_OBSCURE_INDEX = 3
EXPERT_INDEX = 4
LUNATIC_INDEX = 5

DIFFICULTIES = 6

NORMAL_STR = "normal"
NORMAL_OBSCURE_STR = "normal + obscure"
HARD_STR = "hard"
HARD_OBSCURE_STR = "hard + obscure"
EXPERT_STR = "expert"
LUNATIC_STR = "lunatic"

def convert_index_to_str(difficulty_index: int) -> str:
    if difficulty_index == NORMAL_INDEX:
        return NORMAL_STR
    if difficulty_index == NORMAL_OBSCURE_INDEX:
        return NORMAL_OBSCURE_STR
    if difficulty_index == HARD_INDEX:
        return HARD_STR
    if difficulty_index == HARD_OBSCURE_INDEX:
        return HARD_OBSCURE_STR
    if difficulty_index == EXPERT_INDEX:
        return EXPERT_STR
    if difficulty_index == LUNATIC_INDEX:
        return LUNATIC_STR
    raise Exception(f"unrecognized difficulty index {difficulty_index}")
