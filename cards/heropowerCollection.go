package cards

import "github.com/google/uuid"

func Recruiter() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      1,
		Name:      "Recruiter",
		Attack:    1,
		Health:    1,
		MaxHealth: 1,
		Rarity:    Basic,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/c/cc/Silver_Hand_Recruit_full.jpg",
		Tags:      []string{Minion, Paladin},
	}
}

func SummonRecruiter() *Card {
	return &Card{
		Id:     uuid.NewString(),
		Mana:   2,
		Name:   "Summon Recruiter",
		Rarity: Basic,
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/c/cc/Silver_Hand_Recruit_full.jpg",
		Tags:   []string{Heropower, Paladin},
		Events: map[string]Event{
			EventHeroPower: SummonMinionEvent(Recruiter),
		},
	}
}

func DrawCardFor2Life() *Card {
	return &Card{
		Id:     uuid.NewString(),
		Mana:   2,
		Name:   "Life tap",
		Rarity: Basic,
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/1/1d/Life_Tap_full.jpg",
		Tags:   []string{Heropower, Warlock},
		Events: map[string]Event{
			EventHeroPower: func(ctx *EventContext) error {
				ctx.Target = ctx.Source.Player.Hero
				ctx.Board.DealDamage(ctx.Source, ctx.Target, 2)
				ctx.Board.DrawCard(ctx.This, ctx.Source.Player, 0)
				return nil
			},
		},
	}
}

var HeropowerCollection = []func() *Card{
	Recruiter,
	SummonRecruiter,
}
