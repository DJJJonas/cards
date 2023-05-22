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
		Image:     "https://www.fishkeepingworld.com/wp-content/uploads/2018/02/17-Most-Popular-Freshwater-Fish-Article-Banner.png",
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
		Image:     "https://i.pinimg.com/originals/0a/97/8b/0a978bcf51182a0d441358ce56118532.jpg",
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
							Id:        uuid.NewString(),
							Name:      "+1/+1",
							Attack:    1,
							MaxHealth: 1,
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
		Image:     "http://4.bp.blogspot.com/-PBPdLSUjYNY/UbG0DP746sI/AAAAAAAAHcs/cXlMp1zgrG8/s1600/beautifuk+fish+wallpaper+10.jpg",
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
