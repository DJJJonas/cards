package cards

import (
	"github.com/google/uuid"
)

/*
 * Card creation rules:
 * - The first tag in the card is the type of the card (Spell, Minion, etc).
 */

func TheCoin() *Card {
	return &Card{
		Mana:   0,
		Name:   "The Coin",
		Rarity: Basic,
		Text:   "Gain one Mana Crystal this turn only.",
		Image:  "/imgs/TheCoin.png",
		Tags:   []string{Spell, Neutral},
		Events: map[string]Event{
			EventSpellCast: GainXManaCrystalEvent(1),
		},
	}
}

func ElvenArcher() *Card {
	return &Card{
		Mana:      1,
		Name:      "Elven Archer",
		Attack:    1,
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
			for _, c := range b.AllCharacters() {
				options = append(options, c.Id)
			}
			return options
		},
	}
}

func RockBottom() *Card {
	return &Card{
		Mana:   1,
		Name:   "Rock Bottom",
		Rarity: Rare,
		Text:   "Summon a 1/1 Murloc, then <b>Dredge</b>. If it's also a Murloc, summon one more.",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/7/74/Rock_Bottom_full.jpg",
		Tags:   []string{Spell, Warlock},
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
					ctx.Board.SummonMinion(ctx.This, murloc)
				}
				return nil
			},
		},
	}
}

func AzsharanScavenger() *Card {
	return &Card{
		Mana:      2,
		Name:      "Azsharan Scavenger",
		Attack:    2,
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
							Name:   "+1/+1",
							Attack: 1,
							Health: 1,
						})
					}
				}
				for _, m := range p.Minions {
					if m.Tribe == Murloc {
						for _, c := range p.Hand {
							if c.Tribe == Murloc {
								ctx.Target = c
								ctx.Board.EnchantCard(c, &Enchantment{
									Name:   "+1/+1",
									Attack: 1,
									Health: 1,
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
		Mana:      2,
		Name:      "Voidgill",
		Attack:    3,
		MaxHealth: 2,
		Rarity:    Rare,
		Text:      "<b>Deathrattle:</b> Give all Murlocs in your hand +1/+1.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/e/e3/Voidgill_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				for _, c := range ctx.Target.Player.Hand {
					if c.Tribe == Murloc {
						ctx.Board.EnchantCard(c, &Enchantment{
							Name:   "+1/+1",
							Attack: 1,
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
		Mana:      3,
		Name:      "Bloodscent Vilefin",
		Attack:    3,
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
		Mana:      4,
		Name:      "Seadevil Stinger",
		Attack:    4,
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
		Mana:      8,
		Name:      "Gigafin",
		Attack:    7,
		MaxHealth: 4,
		Rarity:    Legendary,
		Text:      "<b>Colossal +1. Battlecry:</b> Devour all enemy minions. <b>Deathrattle:</b> Spit them back out.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/2/26/Gigafin_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock, Deathrattle},
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
				ctx.Board.SummonMinion(ctx.This, maw)
				return nil
			},
		},
	}
}

func MurlocTinyfin() *Card {
	return &Card{
		Mana:      0,
		Name:      "Murloc Tinyfin",
		Attack:    1,
		MaxHealth: 1,
		Rarity:    Common,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/f/f5/Murloc_Tinyfin_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral},
	}
}

func Murmy() *Card {
	return &Card{
		Mana:      1,
		Name:      "Murmy",
		Attack:    1,
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
		Mana:      2,
		Name:      "Bluegill Warrior",
		Attack:    2,
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
		Mana:      2,
		Name:      "Crabrider",
		Attack:    1,
		MaxHealth: 4,
		Text:      "<b>Rush. Windfury.</b>",
		Rarity:    Common,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/5/59/Crabrider_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Rush, Windfury},
	}
}

func LushwaterScout() *Card {
	return &Card{
		Mana:      2,
		Name:      "Lushwater Scout",
		Attack:    1,
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
		Mana:      3,
		Name:      "Murloc Warleader",
		Attack:    3,
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
		Mana:      3,
		Name:      "Twin-fin Fin Twin",
		Attack:    2,
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
		Mana:      4,
		Name:      "Old Murk-Eye",
		Attack:    2,
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
		Mana:      5,
		Name:      "Gorloc Ravager",
		Attack:    4,
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
						Id:     ctx.This.Id,
						Health: -ctx.Target.GetMaxHealth() + 1,
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
		Mana:      3,
		Name:      "Aldor Peacekeeper",
		Attack:    3,
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
					Attack: -ctx.Target.GetAttack() + 1,
				})
				return nil
			},
		},
	}
}

func Consecration() *Card {
	return &Card{
		Mana:   4,
		Name:   "Consecration",
		Rarity: Basic,
		Text:   "Deal 2 damage to all enemies.",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/b/b4/Consecration_full.jpg",
		Tags:   []string{Spell, Paladin},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				op := ctx.Board.getOpponent(ctx.This.Player)
				ctx.Board.DealDamage(ctx.This, op.Hero, 2)
				for _, c := range op.Minions {
					ctx.Board.DealDamage(ctx.This, c, 2)
				}
				return nil
			},
		},
	}
}

func DivineFavor() *Card {
	return &Card{
		Mana:   3,
		Name:   "Divine Favor",
		Rarity: Rare,
		Text:   "Draw cards until you have as many in hand as your opponent.",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/1/19/Divine_Favor_full.jpg",
		Tags:   []string{Spell, Paladin},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				cardDiff := len(ctx.Board.getOpponent(ctx.This.Player).Hand) - len(ctx.This.Player.Hand)
				for i := 0; i < cardDiff; i++ {
					ctx.Board.DrawCard(ctx.This, ctx.This.Player, 0)
				}
				return nil
			},
		},
	}
}

func SwordOfJustice() *Card {
	return &Card{
		Mana:      3,
		Name:      "Sword of Justice",
		Rarity:    Epic,
		Attack:    1,
		MaxHealth: 5,
		Text:      "After you summon a minion, give it +1/+1 and this loses 1 Durability.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/4/4e/Sword_of_Justice_full.jpg",
		Tags:      []string{Weapon, Paladin},
		Events: map[string]Event{
			EventAfterSummon: func(ctx *EventContext) error {
				if ctx.Target.Player != ctx.This.Player {
					return nil
				}
				ctx.Target.AddEnchantment(&Enchantment{
					Id:     ctx.This.Id,
					Attack: 1,
					Health: 1,
				})
				ctx.Board.LoseDurability(ctx.This, ctx.This, 1)
				return nil
			},
		},
	}
}

func BlessingOfKings() *Card {
	return &Card{
		Mana:   4,
		Name:   "Blessing of Kings",
		Rarity: Basic,
		Text:   "Give a minion +4/+4. (+4 Attack/+4 Health)",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/7/72/Blessing_of_Kings_full.jpg",
		Tags:   []string{Spell, Paladin},
		Targets: func(b *Board) []string {
			targets := []string{}
			for _, c := range b.AllMinionCards() {
				targets = append(targets, c.Id)
			}
			return targets
		},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				ctx.Board.EnchantCard(ctx.Target, &Enchantment{
					Id:     ctx.This.Id,
					Attack: 4,
					Health: 4,
				})
				return nil
			},
		},
	}
}

func KeeperOfUldaman() *Card {
	return &Card{
		Mana:      4,
		Name:      "Keeper of Uldaman",
		Rarity:    Common,
		Attack:    3,
		MaxHealth: 3,
		Text:      "<b>Battlecry:</b> Set a minion's Attack and Health to 3.",
		Image:     "https://i.pinimg.com/originals/b9/84/2a/b9842add813d82af043d15303ecca5d9.png",
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
					Attack: -ctx.Target.GetAttack() + 3,
					Health: -ctx.Target.GetMaxHealth() + 3,
				})
				return nil
			},
		},
	}
}

func StandAgainstDarkness() *Card {
	return &Card{
		Mana:   4,
		Name:   "Stand Against Darkness",
		Rarity: Common,
		Text:   "Summon five 1/1 Silver Hand Recruits.",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/f/ff/Stand_Against_Darkness_full.jpg",
		Tags:   []string{Spell, Paladin},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				for i := 0; i < 5; i++ {
					m := SilverHandRecruiter()
					ctx.Board.SummonMinion(ctx.This, m)
				}
				return nil
			},
		},
	}
}

func TruesilverChampion() *Card {
	return &Card{
		Mana:      4,
		Name:      "Truesilver Champion",
		Rarity:    Basic,
		Attack:    4,
		MaxHealth: 2,
		Text:      "Whenever your hero attacks, restore 3 Health to it.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/9/99/Truesilver_Champion_full.jpg",
		Tags:      []string{Weapon, Paladin},
		Events: map[string]Event{
			EventAfterAttack: func(ctx *EventContext) error {
				if ctx.Source.HasTag(Hero) {
					ctx.Board.Heal(ctx.This, ctx.This.Player.Hero, 3)
				}
				return nil
			},
		},
	}
}

func LayOnHands() *Card {
	return &Card{
		Mana:   8,
		Name:   "Lay on Hands",
		Rarity: Epic,
		Text:   "Restore 8 Health. Draw 3 cards.",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/b/bf/Lay_on_Hands_full.jpg",
		Tags:   []string{Spell, Paladin},
		Targets: func(b *Board) []string {
			targets := []string{}
			for _, c := range b.AllCharacters() {
				targets = append(targets, c.Id)
			}
			return targets
		},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				ctx.Board.Heal(ctx.This, ctx.Target, 8)
				for i := 0; i < 3; i++ {
					ctx.Board.DrawCard(ctx.This, ctx.This.Player, 0)
				}
				return nil
			},
		},
	}
}

func Ashbringer() *Card {
	return &Card{
		Mana:      5,
		Name:      "Ashbringer",
		Rarity:    Epic,
		Attack:    5,
		MaxHealth: 3,
		Image:     "https://wow.gamepedia.com/media/wow.gamepedia.com/a/a6/Ashbringer_TCG.jpg",
		Tags:      []string{Weapon, Paladin},
	}
}

func TirionFordring() *Card {
	return &Card{
		Mana:      8,
		Name:      "Tirion Fordring",
		Attack:    6,
		MaxHealth: 6,
		Rarity:    Legendary,
		Text:      "<b>Divine Shield, Taunt Deathrattle:</b> Equip a 5/3 Ashbringer.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/6/63/Tirion_Fordring_full.jpg",
		Tags:      []string{Minion, Paladin, DivineShield, Taunt, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				w := Ashbringer()
				w.Player = ctx.This.Player
				ctx.Board.EquipWeapon(w)
				return nil
			},
		},
	}
}

func BloodmageThalnos() *Card {
	return &Card{
		Id:          uuid.NewString(),
		Mana:        2,
		Name:        "Bloodmage Thalnos",
		SpellDamage: 1,
		Attack:      1,
		MaxHealth:   1,
		Rarity:      Legendary,
		Text:        "<b>Spell Damage +1. Deathrattle:</b> Draw a card.",
		Image:       "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/e/ed/Bloodmage_Thalnos_full.jpg",
		Tags:        []string{Minion, Neutral, Spellpower, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				ctx.Board.DrawCard(ctx.This, ctx.This.Player, 0)
				return nil
			},
		},
	}
}

func LootHoarder() *Card {
	return &Card{
		Mana:      2,
		Name:      "LootHoarder",
		Attack:    2,
		MaxHealth: 1,
		Rarity:    Common,
		Text:      "<b>Deathrattle:</b> Draw a card.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/d/d6/Loot_Hoarder_full.jpg",
		Tags:      []string{Minion, Neutral, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				ctx.Board.DrawCard(ctx.This, ctx.This.Player, 0)
				return nil
			},
		},
	}
}

func SpawnOfNZoth() *Card {
	return &Card{
		Mana:      3,
		Name:      "Spawn of N'Zoth",
		Attack:    2,
		MaxHealth: 2,
		Rarity:    Common,
		Text:      "<b>Deathrattle:</b> Give your minions +1/+1.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/7/73/Spawn_of_N%27Zoth_full.jpg",
		Tags:      []string{Minion, Neutral, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				for _, m := range ctx.This.Player.Minions {
					m.AddEnchantment(&Enchantment{
						Id:     ctx.This.Id,
						Name:   "Spawn of N'Zoth",
						Attack: 1,
					})
				}
				return nil
			},
		},
	}
}

func PutridSlime() *Card {
	return &Card{
		Mana:      1,
		Attack:    1,
		MaxHealth: 2,
		Rarity:    Basic,
		Text:      "<b>Taunt</b>",
		Image:     "https://gamepedia.cursecdn.com/hearthstone_gamepedia/thumb/f/fe/Slime_full.png/800px-Slime_full.png",
		Tags:      []string{Minion, Neutral, Taunt},
	}
}

func SludgeBelcher() *Card {
	return &Card{
		Mana:      5,
		Name:      "Sludge Belcher",
		Attack:    3,
		MaxHealth: 5,
		Rarity:    Common,
		Text:      "<b>Taunt Deathrattle:</b> Summon a 1/2 Slime with Taunt.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/0/01/Sludge_Belcher_full.jpg",
		Tags:      []string{Minion, Neutral, Taunt, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				m := PutridSlime()
				ctx.Board.SummonMinion(ctx.This, m)
				return nil
			},
		},
	}
}

func NZothTheCorruptor() *Card {
	return &Card{
		Mana:      10,
		Name:      "N'Zoth, the Corruptor",
		Attack:    5,
		MaxHealth: 7,
		Rarity:    Legendary,
		Text:      "<b>Battlecry:</b> Summon your <b>Deathrattle</b> minions that died this game.",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/1/13/N%27Zoth%2C_the_Corruptor_full.jpg",
		Tags:      []string{Minion, Neutral, Battlecry},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				for _, m := range ctx.This.Player.Graveyard {
					if !m.HasTag(Deathrattle) {
						continue
					}
					ress := *m
					ressptr := &ress
					ressptr.Id = ""
					ressptr.Enchantments = nil
					ctx.Board.SummonMinion(ctx.This, ressptr)
				}
				return nil
			},
		},
	}
}

var CardCollection = []func() *Card{
	TheCoin,
	ElvenArcher,
}
