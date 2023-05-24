package cards

const (
	EventEndOfAction          = "end_of_action"
	EventBeforeDrawCard       = "before_draw_card"
	EventAfterDrawCard        = "after_draw_card"
	EventBeforeAddToHand      = "before_add_to_hand"
	EventAfterAddToHand       = "after_add_to_hand"
	EventBeforeShuffleDeck    = "before_shuffle_deck"
	EventAfterShuffleDeck     = "after_shuffle_deck"
	EventStartOfGame          = "start_of_game"
	EventStartOfTurn          = "start_of_turn"
	EventEndOfTurn            = "end_of_turn"
	EventAfterCardPlay        = "after_play"
	EventBeforeCardPlay       = "before_play"
	EventBeforeHeroPower      = "before_hero_power"
	EventHeroPower            = "hero_power"
	EventAfterHeroPower       = "after_hero_power"
	EventBattlecry            = "battlecry"
	EventSpellCast            = "spell_cast"
	EventBeforeSummon         = "before_summon"
	EventSummon               = "summoned"
	EventAfterSummon          = "after_summon"
	EventBeforeDamage         = "before_damage"
	EventAfterDamage          = "after_damage"
	EventBeforeHeal           = "before_heal"
	EventAfterHeal            = "after_heal"
	EventBeforeDestroyMinion  = "before_destroy_minion"
	EventDestroyMinion        = "destroy_minion"
	EventDeathrattle          = "deathrattle"
	EventAfterDestroyMinion   = "after_destroy_minion"
	EventBeforeAttack         = "before_attack"
	EventAfterAttack          = "after_attack"
	EventBeforeLoseDurability = "before_lose_durability"
	EventAfterLoseDurability  = "after_lose_durability"
	EventBeforeWeaponDestroy  = "before_weapon_destroy"
	EventAfterWeaponDestroy   = "after_weapon_destroy"
	EventBeforeWeaponEquip    = "before_weapon_equip"
	EventAfterWeaponEquip     = "after_weapon_equip"

	// Historic
	Draw    = "draw"
	Attack  = "attack"
	Heal    = "heal"
	Play    = "play"
	EndTurn = "endturn"
	Summon  = "summon"
)

type Event func(ctx *EventContext) error

type EventContext struct {
	Board        *Board
	This         *Card
	Source       *Card
	Target       *Card
	HealAmount   int
	DamageAmount int
}
