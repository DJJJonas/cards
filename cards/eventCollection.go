/*
 * Collection with event factories.
 * All functions return an Event.
 */
package cards

import (
	"log"

	"github.com/google/uuid"
)

func LogEvent(msg string) Event {
	return func(ctx *EventContext) error {
		log.Printf("%v\n", msg)
		return nil
	}
}

func GainXManaCrystalEvent(count int) Event {
	return func(ctx *EventContext) error {
		p := ctx.Source.Player
		if p.Mana == p.MaxMaxMana {
			return nil
		}
		p.Mana += count
		return nil
	}
}

func DealXDamageEvent(x int) Event {
	return func(ctx *EventContext) error {
		ctx.Board.DealDamage(ctx.Source, ctx.Target, x)
		return nil
	}
}

func SummonMinionEvent(minion func() *Card) Event {
	return func(ctx *EventContext) error {
		m := minion()
		m.Player = ctx.Source.Player
		ctx.Board.SummonMinion(ctx.This, m)
		return nil
	}
}

func GivePlusXYToMinion(x, y int) Event {
	return func(ctx *EventContext) error {
		ctx.Board.EnchantCard(ctx.Target, &Enchantment{
			Id:     uuid.NewString(),
			Name:   "+1/+1",
			Attack: x,
			Health: y,
			Text:   "",
		})
		return nil
	}
}

func DelEnchatmentFromAllies(id string) Event {
	return func(ctx *EventContext) error {
		for _, c := range append(ctx.Target.Player.Minions, []*Card{ctx.Target.Player.Hero, ctx.Target.Player.HeroPower, ctx.Target.Player.Weapon}...) {
			if c.Id == id {
				c.DelEnchId(id)
			}
		}
		return nil
	}
}
