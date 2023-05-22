package cards

import (
	"github.com/google/uuid"
)

/*
 * Card creation rules:
 * - The first tag in the card is the type of the card (Spell, Minion, etc).
 */

func Coin() *Card {
	return &Card{
		Id:     uuid.NewString(),
		Mana:   0,
		Name:   "Coin",
		Rarity: Basic,
		Text:   "Gain one Mana Crystal this turn only.",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/a/a9/The_Coin_full.jpg",
		Tags:   []string{Spell, Neutral},
		Events: map[string]Event{
			EventSpellCast: GainXManaCrystalEvent(1),
		},
	}
}

func ElvenArcher() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      1,
		Name:      "Elven Archer",
		Attack:    1,
		Health:    1,
		MaxHealth: 1,
		Rarity:    Basic,
		Text:      "<b>Battlecry:</b> Deal 1 damage.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/2/2f/Elven_Archer_full.jpg",
		Tags:      []string{Minion, Neutral},
		Events: map[string]Event{
			EventBattlecry: DealXDamageEvent(1),
		},
		Targets: func(b *Board) []string {
			options := []string{}
			for _, c := range b.characters() {
				options = append(options, c.Id)
			}
			return options
		},
	}
}

func RockBottom() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      1,
		Name:      "Rock Bottom",
		Attack:    0,
		Health:    0,
		MaxHealth: 0,
		Rarity:    Rare,
		Text:      "Summon a 1/1 Murloc, then <b>Dredge</b>. If it's also a Murloc, summon one more.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/7/74/Rock_Bottom_full.jpg",
		Tags:      []string{Spell, Warlock},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				p := ctx.Source.Player
				murloc := RockBottomMurloc()
				murloc.Player = p
				ctx.Board.SummonMinion(ctx.This, murloc)
				card := ctx.Board.Dredge(p)
				if card == nil {
					return nil
				}
				if card.Tribe == Murloc {
					murloc := RockBottomMurloc()
					murloc.Player = p
					ctx.Board.SummonMinion(ctx.This, murloc)
				}
				return nil
			},
		},
	}
}

func AzsharanScavenger() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      2,
		Name:      "Azsharan Scavenger",
		Attack:    2,
		Health:    3,
		MaxHealth: 3,
		Rarity:    Common,
		Text:      "<b>Battlecry:</b> Put a 'Sunken Scavenger' on the bottom of your deck.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/f/f8/Azsharan_Scavenger_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				card := SunkenScavenger()
				card.Player = ctx.This.Player
				p := ctx.This.Player
				// TODO: create method to add to deck and assign player
				p.Deck = append(p.Deck, card)
				return nil
			},
		},
	}
}

func ChumBucket() *Card {
	return &Card{
		Id:     uuid.NewString(),
		Mana:   2,
		Name:   "Chum Bucket",
		Rarity: Epic,
		Text:   "Give all Murlocs in your hand +1/+1. Repeat for each Murloc you control.",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/6/61/Chum_Bucket_full.jpg",
		Tags:   []string{Spell, Warlock},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				p := ctx.This.Player
				for _, c := range p.Hand {
					if c.Tribe == Murloc {
						ctx.Target = c
						ctx.Board.EnchantCard(c, &Enchantment{
							Id:        uuid.NewString(),
							Name:      "+1/+1",
							Attack:    1,
							MaxHealth: 1,
						})
					}
				}
				for _, m := range p.Minions {
					if m.Tribe == Murloc {
						for _, c := range p.Hand {
							if c.Tribe == Murloc {
								ctx.Target = c
								ctx.Board.EnchantCard(c, &Enchantment{
									Id:        uuid.NewString(),
									Name:      "+1/+1",
									Attack:    1,
									MaxHealth: 1,
								})
							}
						}
					}
				}
				return nil
			},
		},
	}
}

func Voidgill() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      2,
		Name:      "Voidgill",
		Attack:    3,
		Health:    2,
		MaxHealth: 2,
		Rarity:    Rare,
		Text:      "<b>Deathrattle:</b> Give all Murlocs in your hand +1/+1.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/e/e3/Voidgill_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				for _, c := range ctx.Target.Player.Hand {
					if c.Tribe == Murloc {
						ctx.Board.EnchantCard(c, &Enchantment{
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

func BloodscentVilefin() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      3,
		Name:      "Bloodscent Vilefin",
		Attack:    3,
		Health:    4,
		MaxHealth: 4,
		Rarity:    Rare,
		Text:      "<b>Battlecry: Dredge.</b> If it's a Murloc, change its Cost to Health instead of Mana.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/9/9a/Bloodscent_Vilefin_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				p, b := ctx.This.Player, ctx.Board
				card := ctx.Board.Dredge(p)
				if card != nil && card.Tags[0] == Minion && card.Tribe == Murloc {
					b.EnchantCard(card, &Enchantment{
						Tags: []string{BloodPayment},
					})
				}
				return nil
			},
		},
	}
}

func SeadevilStinger() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      4,
		Name:      "Seadevil Stinger",
		Attack:    4,
		Health:    2,
		MaxHealth: 2,
		Rarity:    Rare,
		Text:      "<b>Battlecry:</b> The next Murloc you play this turn costs Health instead of Mana.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/e/e8/Seadevil_Stinger_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				h, b := ctx.This.Player.Hero, ctx.Board
				enchid := uuid.NewString()
				blpmid := uuid.NewString()
				for _, c := range h.Player.Hand {
					b.EnchantCard(c, &Enchantment{
						Id:   blpmid,
						Tags: []string{BloodPayment},
					})
				}
				b.EnchantCard(h, &Enchantment{
					Id: enchid,
					Events: map[string]Event{
						EventAfterAddToHand: func(ctx *EventContext) error {
							if ctx.Source.Player != ctx.This.Player {
								return nil
							}
							c := ctx.Target
							b.EnchantCard(c, &Enchantment{
								Id:   blpmid,
								Tags: []string{BloodPayment},
							})
							return nil
						},
						EventAfterSummon: func(ctx *EventContext) error {
							if ctx.Target.Player != ctx.This.Player {
								return nil
							}
							if ctx.Target.Tribe == Murloc {
								for _, c := range h.Player.Hand {
									c.DelEnchId(blpmid)
								}
								ctx.This.DelEnchId(enchid)
								return nil
							}
							return nil
						},
						EventEndOfTurn: func(ctx *EventContext) error {
							for _, c := range h.Player.Hand {
								c.DelEnchId(blpmid)
							}
							ctx.This.DelEnchId(enchid)
							return nil
						},
					},
				})
				return nil
			},
		},
	}
}

func Gigafin() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      8,
		Name:      "Gigafin",
		Attack:    7,
		Health:    4,
		MaxHealth: 4,
		Rarity:    Legendary,
		Text:      "<b>Colossal +1. Battlecry:</b> Devour all enemy minions. <b>Deathrattle:</b> Spit them back out.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/2/26/Gigafin_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				p := ctx.Source.Player
				opp := ctx.Board.getOpponent(p)
				oppminions := append([]*Card{}, opp.Minions...)
				opp.Minions = []*Card{}
				enchId := uuid.NewString()
				ench := &Enchantment{
					Id: enchId,
					Events: map[string]Event{
						EventDeathrattle: func(ctx *EventContext) error {
							for _, c := range oppminions {
								ctx.Board.SummonMinion(ctx.This, c)
							}
							return nil
						},
					},
				}
				ctx.Board.EnchantCard(ctx.Source, ench)
				maw := GenerateGigafinMaw(ctx.Source.Id, enchId)
				maw.Player = p
				ctx.Board.SummonMinion(ctx.This, maw)
				return nil
			},
		},
	}
}

func MurlocTinyfin() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      0,
		Name:      "Murloc Tinyfin",
		Attack:    1,
		Health:    1,
		MaxHealth: 1,
		Rarity:    Common,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/f/f5/Murloc_Tinyfin_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral},
	}
}

func Murmy() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      1,
		Name:      "Murmy",
		Attack:    1,
		Health:    1,
		MaxHealth: 1,
		Text:      "<b>Reborn</b>",
		Rarity:    Common,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/a/ad/Murmy_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Reborn},
	}
}

func BluegillWarrior() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      2,
		Name:      "Bluegill Warrior",
		Attack:    2,
		Health:    1,
		MaxHealth: 1,
		Text:      "Charge",
		Rarity:    Basic,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/f/f2/Bluegill_Warrior_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Charge},
	}
}

func Crabrider() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      2,
		Name:      "Crabrider",
		Attack:    2,
		Health:    2,
		MaxHealth: 2,
		Text:      "<b>Rush. Windfury.</b>",
		Rarity:    Common,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/5/59/Crabrider_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Rush, Windfury},
	}
}

func LushwaterScout() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      2,
		Name:      "Lushwater Scout",
		Attack:    1,
		Health:    3,
		MaxHealth: 3,
		Rarity:    Common,
		Text:      "After you summon a Murloc, give it +1 Attack and <b>Rush.</b>",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/6/6b/Lushwater_Scout_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral},
		Events: map[string]Event{
			EventAfterSummon: func(ctx *EventContext) error {
				card := ctx.Source
				if card.Player != ctx.This.Player || card == ctx.This || card.Tribe != Murloc {
					return nil
				}
				ench := &Enchantment{
					Id:     uuid.NewString(),
					Attack: 1,
					Tags:   []string{Rush},
				}
				ctx.Board.EnchantCard(card, ench)
				return nil
			},
		},
	}
}

func MurlocWarleader() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      3,
		Name:      "Murloc Warleader",
		Attack:    3,
		Health:    3,
		MaxHealth: 3,
		Rarity:    Epic,
		Text:      "Your other Murlocs have +2 Attack.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/8/82/Murloc_Warleader_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral},
		Events: map[string]Event{
			EventEndOfAction: func(ctx *EventContext) error {
				for _, c := range ctx.This.Player.Minions {
					if c.Id != ctx.This.Id && c.Tribe == Murloc && c.GetEnch(ctx.This.Id) == nil {
						ctx.Board.EnchantCard(c, &Enchantment{
							Id:     ctx.This.Id,
							Attack: 2,
						})
					}
				}
				return nil
			},
			EventDestroyMinion: func(ctx *EventContext) error {
				for _, c := range ctx.This.Player.Minions {
					if c.GetEnch(ctx.This.Id) != nil {
						c.DelEnchId(ctx.This.Id)
					}
				}
				return nil
			},
		},
	}
}

func TwinfinFinTwin() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      3,
		Name:      "Twin-fin Fin Twin",
		Attack:    2,
		Health:    1,
		MaxHealth: 1,
		Rarity:    Rare,
		Text:      "<b>Rush. Battlecry:</b> Summon a copy of this.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/4/4d/Twin-fin_Fin_Twin_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Rush},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				copy := *ctx.This
				copyptr := &copy
				copyptr.Id = uuid.NewString()
				ctx.Board.SummonMinion(ctx.This, copyptr)
				return nil
			},
		},
	}
}

func OldMurkEye() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      4,
		Name:      "Old Murk-Eye",
		Attack:    2,
		Health:    4,
		MaxHealth: 4,
		Rarity:    Legendary,
		Text:      "<b>Charge.</b> Has +1 Attack for each other Murloc on the battlefield.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/6/62/Murloc_Raid_art.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Charge},
		Events: map[string]Event{
			EventEndOfAction: func(ctx *EventContext) error {
				murlocCount := 0
				for _, c := range append(ctx.This.Player.Minions, ctx.Board.getOpponent(ctx.This.Player).Minions...) {
					if c.Tribe == Murloc && c.Id != ctx.This.Id {
						murlocCount++
					}
				}
				ench := ctx.This.GetEnch(ctx.This.Id)
				if ench == nil {
					ctx.This.AddEnchantment(&Enchantment{
						Id:     ctx.This.Id,
						Attack: murlocCount,
					})
					return nil
				}
				ench.Attack = murlocCount
				return nil
			},
		},
	}
}

func GorlocRavager() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      5,
		Name:      "Gorloc Ravager",
		Attack:    4,
		Health:    3,
		MaxHealth: 3,
		Rarity:    Common,
		Text:      "<b>Battlecry:</b> Draw 3 Murlocs.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/a/a6/Gorloc_Ravager_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				count := 3
				p := ctx.Source.Player
				for i, card := range p.Deck {
					if card.Tribe == Murloc {
						ctx.Board.DrawCard(ctx.This, p, byte(i))
						count--
					}
					if count <= 0 {
						break
					}
				}
				return nil
			},
		},
	}
}

func Redemption() *Card {
	return &Card{
		Id:     uuid.NewString(),
		Mana:   1,
		Name:   "Redemption",
		Rarity: Common,
		Text:   "Secret: When a friendly minion dies, return it to life with 1 Health.",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/a/a8/Redemption_full.jpg",
		Tags:   []string{Spell, Paladin, Secret},
		Events: map[string]Event{
			EventAfterDestroyMinion: func(ctx *EventContext) error {
				if ctx.Target.Player == ctx.This.Player {
					ctx.Target.Health = 1
					ctx.Board.SummonMinion(ctx.This, ctx.Target)
				}
				return nil
			},
		},
	}
}

func Equality() *Card {
	return &Card{
		Id:     uuid.NewString(),
		Mana:   2,
		Name:   "Equality",
		Rarity: Rare,
		Text:   "Change the Health of ALL minions to 1.",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/b/bf/Equality_full.jpg",
		Tags:   []string{Spell, Paladin},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				for _, c := range append(ctx.This.Player.Minions, ctx.Board.getOpponent(ctx.This.Player).Minions...) {
					ctx.Board.EnchantCard(c, &Enchantment{
						Id:        ctx.This.Id,
						MaxHealth: -c.MaxHealth + 1,
					})
					c.Health = 1
				}
				return nil
			},
		},
	}
}

func AldorPeacekeeper() *Card {
	return &Card{
		Id:        uuid.NewString(),
		Mana:      3,
		Name:      "Aldor Peacekeeper",
		Attack:    3,
		Health:    3,
		MaxHealth: 3,
		Rarity:    Rare,
		Text:      "Battlecry: Change an enemy minion's Attack to 1.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/0/0a/Aldor_Peacekeeper_full.jpg",
		Tags:      []string{Minion, Paladin},
		Targets: func(b *Board) []string {
			targets := []string{}
			for _, c := range b.AllMinionCards() {
				targets = append(targets, c.Id)
			}
			return targets
		},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				ctx.Board.EnchantCard(ctx.Target, &Enchantment{
					Id:     ctx.This.Id,
					Attack: -ctx.Target.Attack + 1,
				})
				return nil
			},
		},
	}
}

var CardCollection = []func() *Card{
	Coin,
	ElvenArcher,
}
