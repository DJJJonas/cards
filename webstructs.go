package main

import "cards/cards"

type Enchantment struct {
	Id          string   `json:"id"`
	Mana        int      `json:"mana"`
	Name        string   `json:"name"`
	Text        string   `json:"text"`
	Attack      int      `json:"attack"`
	MaxHealth   int      `json:"maxHealth"`
	SpellDamage int      `json:"spellDamage"`
	Tags        []string `json:"tags"`
}

type Card struct {
	Id           string         `json:"id"`
	Mana         int            `json:"mana"`
	Name         string         `json:"name"`
	Attack       int            `json:"attack"`
	Health       int            `json:"health"`
	MaxHealth    int            `json:"maxHealth"`
	Rarity       string         `json:"rarity"`
	Text         string         `json:"text"`
	Image        string         `json:"image"`
	Type         string         `json:"type"`
	Tags         []string       `json:"tags"`
	Tribe        string         `json:"tribe"`
	Sleeping     bool           `json:"sleeping"`
	CanAttack    bool           `json:"canAttack"`
	Targets      []string       `json:"targets"`
	Enchantments []*Enchantment `json:"enchantments"`
}

type Event struct {
	Type   string `json:"type"`
	Turn   int    `json:"turn"`
	HeroId string `json:"heroId"`
	Source *Card  `json:"source"`
	Target *Card  `json:"target"`
}

type Board struct {
	MyTurn     bool     `json:"myTurn"`
	TurnCount  int      `json:"turnCount"`
	LastEvents []*Event `json:"lastEvents"`

	MyHero       *Card   `json:"myHero"`
	MyHeroPower  *Card   `json:"myHeroPower"`
	MyWeapon     *Card   `json:"myWeapon"`
	MyHand       []*Card `json:"myHand"`
	MyMinions    []*Card `json:"myMinions"`
	MyMana       int     `json:"myMana"`
	MyMaxMana    int     `json:"myMaxMana"`
	MyMaxMaxMana int     `json:"myMaxMaxMana"`
	MyDeckSize   int     `json:"myDeckSize"`

	EnemyHero       *Card   `json:"enemyHero"`
	EnemyHeroPower  *Card   `json:"enemyHeroPower"`
	EnemyWeapon     *Card   `json:"enemyWeapon"`
	EnemyHandSize   int     `json:"enemyHandSize"`
	EnemyMinions    []*Card `json:"enemyMinions"`
	EnemyMana       int     `json:"enemyMana"`
	EnemyMaxMana    int     `json:"enemyMaxMana"`
	EnemyMaxMaxMana int     `json:"enemyMaxMaxMana"`
	EnemyDeckSize   int     `json:"enemyDeckSize"`
}

func TranslateEvents(b *cards.Board, es []*cards.HistoricEvent) []*Event {
	translatedEvents := []*Event{}
	for _, e := range es {
		translatedEvents = append(translatedEvents, TranslateEvent(b, e))
	}
	return translatedEvents
}

func TranslateEvent(b *cards.Board, e *cards.HistoricEvent) *Event {
	return &Event{
		Type:   e.Type,
		Turn:   e.Turn,
		HeroId: e.Player.Hero.Id,
		Source: TranslateCard(b, e.Source),
		Target: TranslateCard(b, e.Target),
	}
}

func TranslateEnchantment(e *cards.Enchantment) *Enchantment {
	return &Enchantment{
		Id:          e.Id,
		Mana:        e.ManaCost,
		Name:        e.Name,
		Attack:      e.Attack,
		MaxHealth:   e.MaxHealth,
		SpellDamage: e.SpellDamage,
		Tags:        e.Tags,
		Text:        e.GetFormattedText(),
	}
}

func TranslateEnchantments(es []*cards.Enchantment) []*Enchantment {
	translatedEnchantments := []*Enchantment{}
	for _, e := range es {
		translatedEnchantments = append(translatedEnchantments, TranslateEnchantment(e))
	}
	return translatedEnchantments
}

func TranslateCard(b *cards.Board, c *cards.Card) *Card {
	if c == nil {
		return nil
	}
	var targets []string
	if c.Targets != nil {
		targets = c.Targets(b)
	}
	return &Card{
		Id:           c.Id,
		Mana:         c.Mana,
		Name:         c.Name,
		Attack:       c.Attack,
		Health:       c.Health,
		MaxHealth:    c.MaxHealth,
		Rarity:       c.Rarity,
		Text:         c.GetFormattedText(),
		Image:        c.Image,
		Type:         c.Tags[0],
		Tags:         c.Tags,
		Tribe:        c.Tribe,
		Sleeping:     c.Sleeping,
		CanAttack:    c.AttacksLeft > 0 && (c.GetAttack() > 0 || c.HasTag(cards.Hero)),
		Targets:      targets,
		Enchantments: TranslateEnchantments(c.Enchantments),
	}
}

func TranslateCards(b *cards.Board, cs []*cards.Card) []*Card {
	translatedCards := []*Card{}
	for _, c := range cs {
		translatedCards = append(translatedCards, TranslateCard(b, c))
	}
	return translatedCards
}

func TranslateBoard(b *cards.Board, playerId byte) *Board {
	enemyId := 1 - playerId
	return &Board{
		MyTurn:     playerId == b.PlayerTurn,
		TurnCount:  b.TurnCount,
		LastEvents: TranslateEvents(b, b.LastEvents),

		MyHero:       TranslateCard(b, b.Players[playerId].Hero),
		MyHeroPower:  TranslateCard(b, b.Players[playerId].HeroPower),
		MyWeapon:     TranslateCard(b, b.Players[playerId].Weapon),
		MyHand:       TranslateCards(b, b.Players[playerId].Hand),
		MyMinions:    TranslateCards(b, b.Players[playerId].Minions),
		MyMana:       b.Players[playerId].Mana,
		MyMaxMana:    b.Players[playerId].MaxMana,
		MyMaxMaxMana: b.Players[playerId].MaxMaxMana,
		MyDeckSize:   len(b.Players[playerId].Deck),

		EnemyHero:       TranslateCard(b, b.Players[enemyId].Hero),
		EnemyHeroPower:  TranslateCard(b, b.Players[enemyId].HeroPower),
		EnemyWeapon:     TranslateCard(b, b.Players[enemyId].Weapon),
		EnemyHandSize:   len(b.Players[enemyId].Hand),
		EnemyMinions:    TranslateCards(b, b.Players[enemyId].Minions),
		EnemyMana:       b.Players[enemyId].Mana,
		EnemyMaxMana:    b.Players[enemyId].MaxMana,
		EnemyMaxMaxMana: b.Players[enemyId].MaxMaxMana,
		EnemyDeckSize:   len(b.Players[enemyId].Deck),
	}
}
