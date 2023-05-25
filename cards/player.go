package cards

type Player struct {
	Image         string
	Hero          *Card
	HeroPower     *Card
	UsedHeroPower bool
	Weapon        *Card
	MaxMinions    int
	Minions       []*Card
	MaxHand       int
	Hand          []*Card
	MaxDeck       int
	Deck          []*Card
	// Cards that have been destroyed
	Graveyard []*Card
	// Cards that have been discarded
	DiscardPile []*Card
	// Cards that have been burned
	Crematorium     []*Card
	MaxSecrets      int
	Secrets         []*Card
	MaxMaxMana      int
	MaxMana         int
	Mana            int
	OverchargedMana int
}
