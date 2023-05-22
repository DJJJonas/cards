export interface Action {
  Type: string;
  SourceId: string;
  TargetId: string;
  Position: number;
}

export interface Event {
  type: string;
  turn: number;
  heroId: string;
  source: Card;
  target: Card;
}

export interface Enchantment {
  id: string;
  mana: number;
  name: string;
  attack: number;
  maxHealth: number;
  spellDamage: number;
  tags: string[];
  text: string;
}

export interface Card {
  id: string;
  mana: number;
  name: string;
  attack: number;
  health: number;
  maxHealth: number;
  shield: number;
  rarity: string;
  text: string;
  image: string;
  type: string;
  tribe: string;
  sleeping: boolean;
  canAttack: boolean;
  targets: string[] | null;
  tags: string[] | null;
  enchantments: Enchantment[] | null;
  // Frontend only fields
  selected: boolean;
  targeted: boolean;
  playable: boolean;
}

export interface Board {
  myTurn: boolean; // 0, 1
  turnCount: number;
  lastEvents: Event[];

  myHero: Card;
  myHeroPower: Card;
  myWeapon: Card | null;
  myHand: Card[];
  myMinions: Card[];
  myMana: number;
  myMaxMana: number;
  myMaxMaxMana: number;
  myDeckSize: number;

  enemyHero: Card;
  enemyHeroPower: Card;
  enemyWeapon: Card | null;
  enemyHandSize: number;
  enemyMinions: Card[];
  enemyMana: number;
  enemyMaxMana: number;
  enemyMaxMaxMana: number;
  enemyDeckSize: number;
}

export function queryCard(b: Board, id: string): Card | undefined {
  const toSearch = [
    ...b.myMinions,
    ...b.myHand,
    b.myHero,
    b.myHeroPower,
    b.myWeapon,
    ...b.enemyMinions,
    b.enemyHero,
    b.enemyHeroPower,
    b.enemyWeapon,
  ];

  return toSearch.find((c) => c?.id === id);
}

export function getCardAttack(card: Card): number {
  let enchattack = 0;
  if (card.enchantments) {
    for (let i = 0; i < card.enchantments.length; i++) {
      enchattack += card.enchantments[i].attack;
    }
  }
  return card.attack + enchattack;
}

export function getMaxHealth(card: Card): number {
  let enchmaxhealth = 0;
  if (card.enchantments) {
    for (let i = 0; i < card.enchantments.length; i++) {
      enchmaxhealth += card.enchantments[i].maxHealth;
    }
  }
  return card.maxHealth + enchmaxhealth;
}

export function getCost(card: Card): number {
  let enchcost = 0;
  if (card.enchantments) {
    for (let i = 0; i < card.enchantments.length; i++) {
      enchcost += card.enchantments[i].mana;
    }
  }
  return card.mana + enchcost;
}

export function getTags(card: Card): string[] {
  let tags: string[] = [];
  card.tags?.forEach((tag) => {
    tags.push(tag);
  });
  card.enchantments?.forEach((ench) => {
    ench.tags?.forEach((tag) => {
      tags.push(tag);
    });
  });
  return tags;
}
