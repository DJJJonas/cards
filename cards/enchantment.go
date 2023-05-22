package cards

import (
	"strconv"
	"strings"
)

type Enchantment struct {
	Id          string
	ManaCost    int
	Name        string
	Attack      int
	MaxHealth   int
	SpellDamage int
	Text        string
	Tags        []string
	Events      map[string]Event
}

func (e *Enchantment) GetFormattedText() string {
	text := e.Text
	text = strings.ReplaceAll(text, "%SpellDamage%", strconv.Itoa(e.SpellDamage))
	return text
}
