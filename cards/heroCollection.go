package cards

import "github.com/google/uuid"

func GuldanHero() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Name:      "Guldan",
		Health:    30,
		MaxHealth: 30,
		Rarity:    Basic,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/2/2f/Gul%27dan_full.jpg",
		Tags:      []string{Hero, Warlock},
	}
}

func UtherHero() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Name:      "Uther",
		Health:    30,
		MaxHealth: 30,
		Rarity:    Basic,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/2/2f/Uther_full.jpg",
		Tags:      []string{Hero, Paladin},
	}
}

var HeroCollection = []func() *Card{
	GuldanHero,
	UtherHero,
}
