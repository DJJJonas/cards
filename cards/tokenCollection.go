package cards

import "github.com/google/uuid"

func RockBottomMurloc() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      1,
		Name:      "Murloc",
		Attack:    1,
		Health:    1,
		MaxHealth: 1,
		Rarity:    Basic,
		Tribe:     Murloc,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/6/6e/Coldlight_Lurker_full.jpg",
		Tags:      []string{Minion, Neutral},
	}
}

func SunkenScavenger() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      2,
		Name:      "Sunken Scavenger",
		Attack:    2,
		Health:    3,
		MaxHealth: 3,
		Text:      "<b>Battlecry:</b> Give your other Murlocs +1/+1 <i>(wherever they are)</i>.",
		Rarity:    Basic,
		Tribe:     Murloc,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/1/1b/Sunken_Scavenger_full.jpg",
		Tags:      []string{Minion, Neutral},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				cards := ctx.Source.Player.Deck
				cards = append(cards, ctx.Source.Player.Hand...)
				cards = append(cards, ctx.Source.Player.Minions...)
				for _, card := range cards {
					if card.Tribe == Murloc && card.Id != ctx.Source.Id {
						ctx.Target = card
						ctx.Board.EnchantCard(card, &Enchantment{
							Id:     uuid.NewString(),
							Name:   "+1/+1",
							Attack: 1,
							Health: 1,
						})
					}
				}
				return nil
			},
		},
	}
}

func GenerateGigafinMaw(gigafinId, enchId string) *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      2,
		Name:      "Gigafin Maw",
		Attack:    4,
		Health:    7,
		MaxHealth: 7,
		Text:      "<b>Taunt. Deathrattle:</b> Permanently destroy all minions inside Gigafin.",
		Rarity:    Legendary,
		Tribe:     Murloc,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/2/26/Gigafin_full.jpg",
		Tags:      []string{Minion, Neutral, Taunt},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				for _, c := range ctx.This.Player.Minions {
					if c.Id == gigafinId {
						for i, e := range c.Enchantments {
							if e.Id == enchId {
								c.Enchantments = append(c.Enchantments[:i], c.Enchantments[i+1:]...)
								break
							}
						}
					}
				}
				return nil
			},
		},
	}
}

var TokenCollection = []func() *Card{
	RockBottomMurloc,
}
