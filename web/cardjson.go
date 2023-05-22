package web

import "cards/cards"

type CardJSON struct {
	Id        string
	Mana      int
	Name      string
	Attack    int
	Health    int
	MaxHealth int
	Rarity    string
	Text      string
	Image     string
	Tribe     string
}

func CardToJSON(card *cards.Card) *CardJSON {
	return &CardJSON{
		Id:        card.Id,
		Mana:      card.Mana,
		Name:      card.Name,
		Attack:    card.Attack,
		Health:    card.Health,
		MaxHealth: card.MaxHealth,
		Rarity:    card.Rarity,
		Text:      card.GetFormattedText(),
		Image:     card.Image,
		Tribe:     card.Tribe,
	}
}
