<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import {
      getCardAttack,
      getCost,
      getTags
  } from "../interfaces/cards";

  export let card: import("../interfaces/cards").Card;
  export let width: number;

  let dispatch = createEventDispatcher();

  let raritymap = new Map();
  raritymap.set("common", "#EFEBE0");
  raritymap.set("rare", "#055C9D");
  raritymap.set("epic", "#8155BA");
  raritymap.set("legendary", "#FFA500");

  let isBleeding = getTags(card).some((t) => t === "bloodPayment");

  function cardShadow(card: import("../interfaces/cards").Card): string {
    let shadows = [];
    if (isBleeding) {
      shadows.push("0 0 12px red");
    }
    if (card.selected) {
      shadows.push("0 0 10px white");
    }
    if (card.targeted) {
      shadows.push("0 0 12px orangered");
    }
    if (card.playable) {
      shadows.push("0 0 8px lightgreen");
    }
    return shadows.join(",");
  }

  function cardBGColor() {
    switch (card.tags[1]) {
      case "neutral":
        return "linear-gradient(to bottom, #6C5E51, #897767)";
      case "druid":
        return "linear-gradient(to bottom, #5C3623, #814C31)";
      case "hunter":
        return "linear-gradient(to bottom, #415E25, #5A8333)";
      case "mage":
        return "linear-gradient(to bottom, #59648F, #737EA8)";
      case "paladin":
        return "linear-gradient(to bottom, #9C7D4F, #7A623E)";
      case "priest":
        return "linear-gradient(to bottom, #A5A7AD, #8A8D95)";
      case "rogue":
        return "linear-gradient(to bottom, #39393E, #515159)";
      case "shaman":
        return "linear-gradient(to bottom, #343F62, #465483)";
      case "warlock":
        return "linear-gradient(to bottom, #4E3058, #6B4279)";
      case "warrior":
        return "linear-gradient(to bottom, #832B20, #AC382A)";
    }
  }

  function attackColor() {
    let enchAttack = getCardAttack(card);
    if (card.attack === enchAttack) return "#E9EAE0";
    if (card.attack > enchAttack) return "#DB1F48";
    if (card.attack < enchAttack) return "#01949A";
  }

  function healthColor() {
    if (card.health === card.maxHealth) return "#E9EAE0";
    if (card.health < card.maxHealth) return "#DB1F48";
    if (card.health > card.maxHealth) return "#01949A";
  }
</script>

<div
  class="card"
  style:width={width + "px"}
  style:height={width * 1.5 + "px"}
  style:box-shadow={cardShadow(card)}
  style:background-image={cardBGColor()}
  on:click={() => dispatch("click")}
  on:keypress={() => dispatch("keypress")}
  on:mouseenter={() => dispatch("mouseenter")}
  on:mouseleave={() => dispatch("mouseleave")}
>
<img class={card.type} src={card.image} alt="card portrait" />
<div class="mana" style:font-size={width * 0.174 + "px"}>{getCost(card)}</div>
  <div class="name" style:font-size={width * 0.092 + "px"}>{card.name}</div>
  <div class="text" style:font-size={width * 0.08 + "px"}>
    {@html card.text}
  </div>
  {#if card.type !== "spell" && card.type !== "heropower"}
    <div
      class="attack"
      style:color={attackColor()}
      style:font-size={width * 0.16 + "px"}
    >
      {getCardAttack(card)}
    </div>
    <div
      class="health"
      style:color={healthColor()}
      style:font-size={width * 0.16 + "px"}
      style:background-image={card.type === "weapon"
        ? "linear-gradient(#283336, #362b28)"
        : "linear-gradient(#e7625f, #c85250)"}
    >
      {card.health}
    </div>
  {/if}
  <div
    style:display={card.rarity === "basic" ? "none" : "block"}
    class="rarity"
    style:background-color={raritymap.get(card.rarity)}
  />
  <div
    style:display={card.tribe ? "block" : "none"}
    style:font-size={width * 0.09 + "px"}
    class="tribe"
  >
    {card.tribe}
  </div>
</div>

<style>
  .card {
    position: relative;
    cursor: pointer;
    margin: 5px;
    color: #f4f2eb;
    text-align: center;
    text-shadow: -1px -1px 1px black, 1px -1px 1px black, -1px 1px 1px black,
      1px 1px 1px black;
    border-style: solid;
    border-width: 1px;
    border-color: #282120;
    border-radius: 4%;
  }
  .mana {
    position: absolute;
    top: -4%;
    left: -4%;
    width: 20%;
    border-style: solid;
    border-width: 1px;
    border-color: #282120;
    border-radius: 40%;
    background-image: linear-gradient(to top, #3434d7, #3b3bd7);
  }
  .image {
    position: absolute;
    object-fit: cover;
    left: 50%;
    translate: -50%;
    border-style: solid;
    border-color: #a49393;
    border-width: 2px;
    border-radius: 50%;
  }

  .minion {
    position: absolute;
    object-fit: cover;
    left: 50%;
    translate: -50%;
    border-style: solid;
    border-color: #a49393;
    border-width: 2px;
    border-radius: 50%;
    top: -6%;
    width: 76%;
    height: 65.2%;
  }

  .weapon {
    position: absolute;
    object-fit: cover;
    left: 50%;
    translate: -50%;
    border-style: solid;
    border-color: #a49393;
    border-width: 2px;
    border-radius: 50%;
    top: -2%;
    width: 84%;
    aspect-ratio: 1;
  }
  .spell {
    position: absolute;
    object-fit: cover;
    left: 50%;
    translate: -50%;
    border-style: solid;
    border-color: #a49393;
    border-width: 2px;
    border-radius: 10%;
    top: 2%;
    width: 90%;
    height: 50%;
  }
  .name {
    position: absolute;
    margin: 0;
    padding: 0;
    font-weight: bolder;
    top: 47%;
    left: 50%;
    translate: -50%;
    background-color: #eed6d3;
    width: 95%;
    height: 10%;
    border-style: solid;
    border-color: black;
    border-width: 1px;
  }
  .text {
    color: #282120;
    position: absolute;
    font-family: Arial;
    text-shadow: none;
    background-color: #eed6d3;
    top: 60%;
    left: 50%;
    translate: -50%;
    width: 95%;
    height: 38%;
  }
  .attack {
    position: absolute;
    top: 88%;
    left: -3%;
    background-image: linear-gradient(#ffed8a, #ffe55c);
    width: 20%;
    height: 13%;
    border-style: solid;
    border-color: #282120;
    border-radius: 50%;
  }
  .health {
    position: absolute;
    top: 88%;
    left: 79%;
    width: 20%;
    height: 13%;
    border-style: solid;
    border-color: #282120;
    border-radius: 50%;
  }
  .rarity {
    position: absolute;
    width: 6%;
    height: 6%;
    top: 54%;
    left: 50%;
    translate: -50%;
    border-style: solid;
    border-color: #282120;
    border-width: 1px;
    border-radius: 50%;
  }
  .tribe {
    position: absolute;
    background-color: #bc9b82;
    width: 56%;
    height: 8%;
    top: 91.5%;
    left: 50%;
    translate: -50%;
    border-style: solid;
    border-width: 1px;
    border-color: #a1846e;
    border-radius: 10% 10% 0 0;
  }
</style>
