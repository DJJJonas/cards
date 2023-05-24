package cards

import (
	"cards/utils"
	"fmt"
)

type HistoricEvent struct {
	Type   string
	Player *Player
	Turn   int
	Source *Card
	Target *Card
}

type Board struct {
	PlayerTurn        byte // 0, 1
	TurnCount         int
	Players           [2]*Player
	History           []*HistoricEvent
	LastEvents        []*HistoricEvent
	ActionChan        chan *Action
	WaitingActionChan chan int
	ActionEndChan     chan error
}

func (b *Board) Start() {
	b.PlayerTurn = 1 // So on next turn, it's player 1's turn
	for _, p := range b.Players {
		setCardPlayers(p)
		b.ShuffleDeck(p)
		b.TriggerEventsFrom(p.Deck, b.Context(p.Hero, p.Hero), EventStartOfGame)
		// TODO: mulligan action mechanic
		for i := 0; i < 5; i++ {
			p.Hand = append(p.Hand, b.DrawCardEventless(p, 0))
		}
	}
	coin := Coin()
	coin.Player = b.Players[1]
	b.Players[1].Hand = append(b.Players[1].Hand, coin)
	b.NextTurn()

gameLoop:
	for {
		p := b.Players[b.PlayerTurn]
		p.UsedHeroPower = false
	turnLoop:
		for {
			b.WaitingActionChan <- -1
			a := <-b.ActionChan
			if !IsActionValid(a) {
				b.ActionEndChan <- fmt.Errorf("invalid action: %v", a)
				continue turnLoop
			}
			b.LastEvents = []*HistoricEvent{}
			actionEnd := b.DoAction(a, p)
			b.TriggerEventsFrom(b.AllActiveCards(), b.Context(p.Hero, p.Hero), EventEndOfAction)
			b.ActionEndChan <- actionEnd
			if p, ok := b.CheckWin(); ok {
				b.WaitingActionChan <- int(p)
				break gameLoop
			}
			if a.Type == EndTurn {
				break turnLoop
			}
		}
	}
}

func (b *Board) AddHistoric(type_ string, source, target *Card) {
	e := &HistoricEvent{
		Type:   type_,
		Player: b.Players[b.PlayerTurn],
		Turn:   b.TurnCount,
		Source: source,
		Target: target,
	}
	b.LastEvents = append(b.LastEvents, e)
	b.History = append(b.History, e)
}

func (b *Board) DoAction(a *Action, p *Player) error {
	switch a.Type {
	case Play:
		card := b.getCardByIdFrom(a.SourceId, p.Hand)
		if card == nil {
			return fmt.Errorf("card not found")
		}
		target := b.getCardByIdFrom(a.TargetId, b.AllCharacters())
		if card.Targets != nil {
			validOptions := card.Targets(b)
			if !b.targetIsValid(target.Id, validOptions) {
				return fmt.Errorf("target not found")
			}
		}
		err := b.PlayFromHand(card, target, a.Position)
		if err != nil {
			return err
		}
		b.AddHistoric(Play, card, target)
	case Heropower:
		target := b.getCardByIdFrom(a.TargetId, b.AllCharacters())
		err := b.UseHeroPower(p.HeroPower, target)
		if err != nil {
			return err
		}
		b.AddHistoric(Heropower, p.HeroPower, target)
	case Attack:
		source := b.getCardByIdFrom(a.SourceId, append(p.Minions, p.Hero))
		target := b.getCardByIdFrom(a.TargetId, append(b.getOpponent(p).Minions, b.getOpponent(p).Hero))
		if source == nil || target == nil {
			return fmt.Errorf("card not found")
		}
		if source.Tags[0] == Minion {
			err := b.MinionAttack(source, target)
			if err != nil {
				return err
			}
		}
		if source.Tags[0] == Hero {
			err := b.HeroAttack(source, target)
			if err != nil {
				return err
			}
		}
	case EndTurn:
		b.AddHistoric(EndTurn, p.Hero, p.Hero)
		b.EndTurn(p)
	}
	return nil
}

func setCardPlayers(p *Player) {
	cards := p.Deck
	cards = append(cards, p.Hero)
	cards = append(cards, p.HeroPower)
	cards = append(cards, p.Minions...)
	cards = append(cards, p.Weapon)
	for _, c := range cards {
		if c != nil {
			c.Player = p
		}
	}
}

func (b *Board) OpponentHasTaunt(p *Player) bool {
	opp := b.getOpponent(p)
	cardsToCheck := append(opp.Minions, opp.Hero)
	for _, c := range cardsToCheck {
		if c.HasTag(Taunt) {
			return true
		}
	}
	return false
}

func (b *Board) EnchantCard(card *Card, enchantment *Enchantment) {
	card.AddEnchantment(enchantment)
}

func (b *Board) Dredge(p *Player) *Card {
	l := len(p.Deck)
	if l == 0 {
		return nil
	}
	dredgableCards := []*Card{}
	count := 3
	for i := l - 1; i >= 0; i-- {
		dredgableCards = append(dredgableCards, p.Deck[i])
		count--
		if count == 0 {
			break
		}
	}
	cardI := utils.RandInt(0, 2)
	card := dredgableCards[cardI]
	p.Deck = b.removeCard(cardI, p.Deck)
	p.Deck = append([]*Card{card}, p.Deck...)
	return card
}

func (b *Board) DealDamage(source, target *Card, damage int) {
	if target.HasTag(DivineShield) {
		target.DelTag(DivineShield)
		return
	}
	ctx := b.Context(source, target)
	ctx.DamageAmount = damage
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeDamage)
	target.Health -= ctx.DamageAmount
	if target.Tags[0] == Minion && target.Health <= 0 {
		b.DestroyMinion(ctx.Source, ctx.Target)
	}
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterDamage)
}

func (b *Board) Heal(source, target *Card, heal int) {
	ctx := b.Context(source, target)
	ctx.HealAmount = heal
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeHeal)
	target.Health += ctx.HealAmount
	if target.Health >= target.MaxHealth {
		target.Health = target.MaxHealth
	}
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterHeal)
}

func (b *Board) HeroAttack(source, target *Card) error {
	if (source.GetAttack() <= 0 && (source.Player.Weapon == nil || source.Player.Weapon.GetAttack() <= 0)) || source.AttacksLeft == 0 {
		return fmt.Errorf("your hero cannot attack")
	}
	if b.OpponentHasTaunt(source.Player) && !target.HasTag(Taunt) {
		return fmt.Errorf("you need to attack the taunt first")
	}
	ctx := b.Context(source, target)
	ctx.DamageAmount = source.GetAttack()
	if source.Player.Weapon != nil {
		ctx.DamageAmount += source.Player.Weapon.GetAttack()
	}
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeAttack)
	b.DealDamage(ctx.Source, ctx.Target, ctx.DamageAmount)
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterAttack)
	source.AttacksLeft--
	if target.GetAttack() > 0 {
		ctx = b.Context(ctx.Target, ctx.Source)
		ctx.DamageAmount = target.GetAttack()
		b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeAttack)
		b.DealDamage(ctx.Source, ctx.Target, ctx.DamageAmount)
		b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterAttack)
	}
	if source.Player.Weapon != nil {
		b.LoseDurability(source.Player.Hero, source.Player.Weapon, 1)
	}
	return nil
}

func (b *Board) MinionAttack(source, target *Card) error {
	hasRush := source.HasTag(Rush)
	hasCharge := source.HasTag(Charge)
	p := source.Player
	if !hasCharge && hasRush && target.Id == b.getOpponent(p).Hero.Id && b.playedThisTurn(source) {
		return fmt.Errorf("minions with rush can't attack heroes")
	}
	if source.GetAttack() <= 0 || source.Sleeping || source.AttacksLeft == 0 {
		return fmt.Errorf("card cannot attack")
	}
	if b.OpponentHasTaunt(p) && !target.HasTag(Taunt) {
		return fmt.Errorf("you need to attack the taunt first")
	}
	if !b.cardIsActiveCard(source) || !b.cardIsActiveCard(target) ||
		source.Health < 0 || target.Health < 0 || source.Attack <= 0 {
		return nil
	}
	ctx := b.Context(source, target)
	ctx.DamageAmount = source.GetAttack()
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeAttack)
	b.DealDamage(ctx.Source, ctx.Target, ctx.DamageAmount)
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterAttack)
	source.AttacksLeft--
	if target.GetAttack() <= 0 {
		return nil
	}
	ctx = b.Context(ctx.Target, ctx.Source)
	ctx.DamageAmount = target.GetAttack()
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeAttack)
	b.DealDamage(ctx.Source, ctx.Target, ctx.DamageAmount)
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterAttack)
	return nil
}

// Makes the source card destroy the minion
func (b *Board) DestroyMinion(source, minion *Card) {
	ctx := b.Context(source, minion)
	p := minion.Player
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeDestroyMinion)
	p.Minions = b.removeCardById(ctx.Target.Id, p.Minions)
	b.TriggerCardEvent(ctx.Target, ctx, EventDestroyMinion)
	b.TriggerCardEvent(ctx.Target, ctx, EventDeathrattle)
	p.Graveyard = append(p.Graveyard, ctx.Target)
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterDestroyMinion)
	if ctx.Target.HasTag(Reborn) {
		cpy := *ctx.Target
		cpy.Health = cpy.MaxHealth
		cpy.Enchantments = nil
		cpy.DelTag(Reborn)
		b.SummonMinion(ctx.Target, &cpy)
	}
}

func (b *Board) CheckWin() (byte, bool) {
	for i, p := range b.Players {
		if p.Hero.Health <= 0 {
			return byte(1 - i), true
		}
	}
	return 0, false
}

func (b *Board) UseHeroPower(hp, target *Card) error {
	if hp.AttacksLeft == 0 {
		return fmt.Errorf("already used")
	}
	if !b.PayFor(hp) {
		return fmt.Errorf("not enough mana")
	}
	var ctx *EventContext
	if hp.Targets != nil {
		if target != nil || !b.targetIsValid(target.Id, hp.Targets(b)) {
			return fmt.Errorf("target not found")
		}
		ctx = b.Context(hp, target)
	} else {
		ctx = b.Context(hp, hp)
	}
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeHeroPower)
	b.TriggerCardEvent(hp, ctx, EventHeroPower)
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterHeroPower)
	hp.AttacksLeft--
	return nil
}

func (b *Board) DrawCard(source *Card, p *Player, i byte) {
	if i >= byte(len(p.Deck)) {
		return
	}
	if len(p.Deck) < 1 {
		// TODO: fatigue
		return
	}
	card := b.DrawCardEventless(p, i)
	if len(p.Hand) >= p.MaxHand {
		p.Crematorium = append(p.Crematorium, card)
		return
	}
	b.AddHistoric(Draw, source, card)
	ctx := b.Context(card, card)
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeDrawCard)
	b.AddToHand(p, ctx.Target)
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterDrawCard)
}

func (b *Board) AddToHand(p *Player, card *Card) {
	b.TriggerEventsFrom(b.AllActiveCards(), b.Context(card, card), EventBeforeAddToHand)
	p.Hand = append(p.Hand, card)
	b.TriggerEventsFrom(b.AllActiveCards(), b.Context(card, card), EventAfterAddToHand)
}

func (b *Board) DrawCardEventless(p *Player, i byte) *Card {
	if len(p.Deck) == 0 {
		return nil
	}
	card := p.Deck[i]
	p.Deck = append(p.Deck[:i], p.Deck[i+1:]...)
	return card
}

func (b *Board) TriggerCardEvent(card *Card, ctx *EventContext, eventName string) {
	if card.Events != nil {
		ctx.This = card
		f, ok := card.Events[eventName]
		if ok {
			f(ctx)
		}
	}
	if card.Enchantments != nil && len(card.Enchantments) > 0 {
		ctx.This = card
		for _, e := range card.Enchantments {
			f, ok := e.Events[eventName]
			if !ok {
				continue
			}
			f(ctx)
			if e.HasTag(Secret) {
				b.AddHistoric(Secret, card, card)
				card.Player.Hero.DelEnchId(e.Id)
				break
			}
		}
	}
}

func (b Board) TriggerEventsFrom(cards []*Card, ctx *EventContext, eventName string) {
	for _, card := range cards {
		b.TriggerCardEvent(card, ctx, eventName)
	}
}

func (b *Board) RefreshManaCrystals(p *Player) {
	p.Mana = p.MaxMana
}

func (b *Board) ShuffleDeck(p *Player) {
	// Fisher-yard shuffle
	for i := len(p.Deck) - 1; i > 0; i-- {
		j := utils.RandInt(0, i+1)
		p.Deck[i], p.Deck[j] = p.Deck[j], p.Deck[i]
	}
}

func (b *Board) StartTurn(p *Player) {
	for _, c := range p.Minions {
		b.TriggerCardEvent(c, b.Context(p.Hero, p.Hero), EventStartOfTurn)
	}
}

func (b *Board) NextTurn() {
	b.PlayerTurn = 1 - b.PlayerTurn
	b.TurnCount++
	p := b.Players[b.PlayerTurn]
	if p.MaxMana < p.MaxMaxMana {
		p.MaxMana++
	}
	b.RefreshManaCrystals(p)
	b.RefreshAttacks(p)
	b.DrawCard(p.Hero, p, 0)
	b.StartTurn(p)
}

func (b *Board) EndTurn(p *Player) {
	for _, c := range p.Minions {
		b.TriggerCardEvent(c, b.Context(p.Hero, p.Hero), EventEndOfTurn)
	}
	b.WakeUpMinions(p)
	b.NextTurn()
}

func (b *Board) RefreshAttacks(p *Player) {
	for _, c := range append(p.Minions, p.Hero, p.HeroPower) {
		b.RefreshAttack(c)
	}
}

func (*Board) RefreshAttack(c *Card) {
	if c.HasTag(MegaWindfury) {
		c.AttacksLeft = 3
	} else if c.HasTag(Windfury) {
		c.AttacksLeft = 2
	} else {
		c.AttacksLeft = 1
	}
}

func (b *Board) WakeUpMinions(p *Player) {
	for _, c := range p.Minions {
		c.Sleeping = false
	}
}

func (b *Board) PlayFromHand(card, target *Card, pos int) error {
	p := card.Player
	// Mana payment
	if !b.PayFor(card) {
		return fmt.Errorf("not enough mana")
	}
	// Max minions rule
	if card.Tags[0] == Minion && len(p.Minions) >= p.MaxMinions {
		return fmt.Errorf("board is full")
	}
	for i, c := range p.Hand {
		if c.Id == card.Id {
			p.Hand = b.removeCard(i, p.Hand)
		}
	}
	b.PlayCard(p.Hero, card, target, pos)
	return nil
}

func (b *Board) PayFor(card *Card) bool {
	p := card.Player
	if card.HasTag(BloodPayment) {
		if p.Hero.Health < card.GetManaCost() {
			return false
		}
		p.Hero.Health -= card.GetManaCost()
		return true
	}
	if p.Mana < card.GetManaCost() {
		return false
	}
	p.Mana -= card.GetManaCost()
	return true
}

func (b *Board) PlayCard(source, card, target *Card, pos int) {
	var ctx *EventContext
	if target != nil {
		ctx = b.Context(card, target)
	} else {
		ctx = b.Context(card, card)
	}
	p := card.Player
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeCardPlay)
	switch card.Tags[0] {
	case Minion:
		b.SummonMinion(source, card)
		b.TriggerCardEvent(card, ctx, EventBattlecry)
	case Spell:
		if card.HasTag(Secret) {
			events := card.Events
			ench := &Enchantment{
				Id:     card.Id,
				Name:   card.Name,
				Text:   card.Text,
				Tags:   []string{Secret},
				Events: events,
			}
			card.Player.Hero.Enchantments = append(card.Player.Hero.Enchantments, ench)
		}
		b.TriggerCardEvent(card, ctx, EventSpellCast)
		b.insertCard(0, p.Graveyard, card)
	case Weapon:
		b.EquipWeapon(card)
		// TODO: case Hero
	}
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterCardPlay)
}

func (b *Board) SummonMinion(source, card *Card) {
	p := card.Player
	if card.Tags[0] == Minion && len(p.Minions) >= p.MaxMinions {
		return
	}
	ctx := b.Context(card, card)
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventBeforeSummon)
	p.Minions = append(p.Minions, card)
	b.TriggerCardEvent(card, ctx, EventSummon)
	b.TriggerEventsFrom(b.AllActiveCards(), ctx, EventAfterSummon)
	card.Sleeping = !(card.HasTag(Charge) || card.HasTag(Rush))
	if !card.Sleeping {
		b.RefreshAttack(card)
	}
	b.AddHistoric(Summon, source, card)
}

func (b *Board) IsCardFrom(card *Card, cards []*Card) bool {
	for _, c := range cards {
		if c == card {
			return true
		}
	}
	return false
}

func (b *Board) AllMinionCards() []*Card {
	var cards []*Card
	for _, p := range b.Players {
		cards = append(cards, p.Minions...)
	}
	return cards
}

func (b *Board) EquipWeapon(weapon *Card) {
	if weapon.Player.Weapon != nil {
		b.DestroyWeapon(weapon.Player.Hero, weapon)
	}
	b.TriggerEventsFrom(b.AllActiveCards(), b.Context(weapon, weapon), EventBeforeWeaponEquip)
	weapon.Player.Weapon = weapon
	b.TriggerEventsFrom(b.AllActiveCards(), b.Context(weapon, weapon), EventAfterWeaponEquip)
}

func (b *Board) LoseDurability(source, target *Card, q int) {
	b.TriggerEventsFrom(b.AllActiveCards(), b.Context(source, target), EventBeforeLoseDurability)
	target.Health -= q
	b.TriggerEventsFrom(b.AllActiveCards(), b.Context(source, target), EventAfterLoseDurability)
	if target.Health <= 0 {
		b.DestroyWeapon(source, target)
	}
}

func (b *Board) DestroyWeapon(source, target *Card) {
	b.TriggerEventsFrom(b.AllActiveCards(), b.Context(source, target), EventBeforeWeaponDestroy)
	target.Player.Graveyard = append(target.Player.Graveyard, target.Player.Weapon)
	target.Player.Weapon = nil
	b.TriggerCardEvent(target, b.Context(source, target), EventAfterWeaponDestroy)
}

func (b *Board) AllActiveCards() []*Card {
	var cards []*Card
	for _, p := range b.Players {
		cards = append(cards, p.Hero)
		if p.Weapon != nil {
			cards = append(cards, p.Weapon)
		}
		if p.HeroPower != nil {
			cards = append(cards, p.HeroPower)
		}
		cards = append(cards, p.Minions...)
	}
	return cards
}
func (b *Board) AllCharacters() []*Card {
	var cards []*Card
	for _, p := range b.Players {
		cards = append(cards, p.Hero)
		cards = append(cards, p.Minions...)
	}
	return cards
}

func (b *Board) getOpponent(p *Player) *Player {
	if p == b.Players[0] {
		return b.Players[1]
	}
	return b.Players[0]
}

func (b *Board) cardIsActiveCard(card *Card) bool {
	for _, c := range b.AllActiveCards() {
		if c == card {
			return true
		}
	}
	return false
}

func (b *Board) insertCard(pos int, arr []*Card, card *Card) []*Card {
	if pos < 0 {
		pos = 0
	} else if (pos) > len(arr) {
		pos = len(arr)
	}
	return append(arr[:pos], append([]*Card{card}, arr[pos:]...)...)
}

func (b *Board) removeCard(pos int, arr []*Card) []*Card {
	return append(arr[:pos], arr[pos+1:]...)
}

func (b *Board) removeCardById(id string, arr []*Card) []*Card {
	for i, c := range arr {
		if c.Id == id {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

func (b *Board) Context(source, target *Card) *EventContext {
	ctx := &EventContext{
		Board: b,
	}
	if source != nil {
		ctx.Source = source
	}
	if target != nil {
		ctx.Target = target
	}
	return ctx
}

func (b *Board) getCardByIdFrom(id string, cards []*Card) *Card {
	for _, c := range cards {
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (b *Board) targetIsValid(target string, targets []string) bool {
	for _, t := range targets {
		if t == target {
			return true
		}
	}
	return false
}

func (b *Board) playedThisTurn(card *Card) bool {
	for i := len(b.History) - 1; i >= 0; i-- {
		h := b.History[i]
		if h.Turn != b.TurnCount {
			break
		}
		if h.Type == Summon && h.Target == card {
			return true
		}
	}
	return false
}
