package cards

import (
	"strconv"
	"strings"
)

const (
	// Classes
	Neutral = "neutral"
	Druid   = "druid"
	Hunter  = "hunter"
	Mage    = "mage"
	Paladin = "paladin"
	Priest  = "priest"
	Rogue   = "rogue"
	Shaman  = "shaman"
	Warlock = "warlock"
	Warrior = "warrior"

	// Rarities
	Basic     = "basic"
	Common    = "common"
	Rare      = "rare"
	Epic      = "epic"
	Legendary = "legendary"

	// Card types
	Minion    = "minion"
	Spell     = "spell"
	Hero      = "hero"
	Heropower = "heropower"
	Weapon    = "weapon"

	// Tribes
	Pirate = "pirate"
	Murloc = "murloc"
	Beast  = "beast"

	// Keywords
	Charge       = "charge"
	Rush         = "rush"
	Battlecry    = "battlecry"
	Windfury     = "windfury"
	MegaWindfury = "megaWindfury"
	Deathrattle  = "deathrattle"
	Taunt        = "taunt"
	BloodPayment = "bloodPayment"
	Reborn       = "reborn"
	DivineShield = "divineShield"
	Secret       = "secret"
	Spellpower   = "spellpower"
	EndOfTurn    = "endOfTurn"
)

type Card struct {
	Id           string
	Player       *Player
	SpellDamage  int
	Mana         int
	Name         string
	Attack       int
	Health       int
	MaxHealth    int
	Rarity       string
	Text         string
	Image        string
	Tags         []string
	Tribe        string
	Events       map[string]Event
	Sleeping     bool
	AttacksLeft  int
	Targets      func(*Board) []string
	Enchantments []*Enchantment
}

func (c *Card) AddEnchantment(e *Enchantment) {
	c.Health += e.Health
	c.Enchantments = append(c.Enchantments, e)
}

// Returns the formatted text with spell power, info for example
func (c *Card) GetFormattedText() string {
	text := c.Text
	text = strings.ReplaceAll(text, "%SpellDamage%", strconv.Itoa(c.SpellDamage))
	return text
}

// Returns the attack value with all enchantments if there are any
func (c *Card) GetAttack() int {
	attack := c.Attack
	for _, e := range c.Enchantments {
		attack += e.Attack
	}
	return attack
}

// Returns the max health value with all enchantments if there are any
func (c *Card) GetMaxHealth() int {
	maxHealth := c.MaxHealth
	for _, e := range c.Enchantments {
		maxHealth += e.Health
	}
	return maxHealth
}

func (c *Card) GetTag(tag string) int {
	for i, t := range c.GetTags() {
		if t == tag {
			return i
		}
	}
	return -1
}

func (c *Card) GetTags() []string {
	tags := c.Tags
	for _, e := range c.Enchantments {
		tags = append(tags, e.Tags...)
	}
	return tags
}

func (c *Card) HasTag(tag string) bool {
	for _, t := range c.GetTags() {
		if t == tag {
			return true
		}
	}
	return false
}

func (c *Card) DelTag(tag string) {
	for i, t := range c.GetTags() {
		if t == tag {
			c.Tags = append(c.Tags[:i], c.Tags[i+1:]...)
		}
	}
}

func (c *Card) GetManaCost() int {
	cost := c.Mana
	for _, e := range c.Enchantments {
		cost += e.ManaCost
	}
	if cost < 0 {
		return 0
	}
	return cost
}

func (c *Card) GetEnch(id string) *Enchantment {
	for _, e := range c.Enchantments {
		if e.Id == id {
			return e
		}
	}
	return nil
}

func (c *Card) DelEnchId(id string) {
	for i, e := range c.Enchantments {
		if e.Id == id {
			c.Enchantments = append(c.Enchantments[:i], c.Enchantments[i+1:]...)
		}
	}
}
