<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import {
    getTags,
    type Card,
    getMaxHealth,
    getCardAttack,
  } from "../interfaces/cards";

  export let card: import("../interfaces/cards").Card;
  export let width: number;

  let dispatch = createEventDispatcher();

  function cardShadow(card: import("../interfaces/cards").Card): string {
    let shadows = [];
    if (card.targeted) {
      shadows.push("0 0 12px orangered");
    }
    return shadows.join(",");
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
  style:height={width + "px"}
  style:box-shadow={cardShadow(card)}
  on:click={() => dispatch("click")}
  on:keypress={() => dispatch("keypress")}
  on:mouseenter={() => dispatch("mouseenter")}
  on:mouseleave={() => dispatch("mouseleave")}
>
  <img class="image" src={card.image} alt="card portrait" />
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
  >
    {card.health}
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
  .image {
    position: absolute;
    object-fit: cover;
    top: -6%;
    left: 50%;
    translate: -50%;
    width: 100%;
    height: 100%;
    border-style: solid;
    border-color: #a49393;
    border-width: 4px;
    border-radius: 50%;
  }
  .attack {
    position: absolute;
    top: 75%;
    left: -3%;
    background-image: linear-gradient(#ffed8a, #ffe55c);
    width: 20%;
    height: 20%;
    border-style: solid;
    border-color: #282120;
    border-radius: 50%;
  }
  .health {
    position: absolute;
    top: 75%;
    left: 79%;
    background-image: linear-gradient(#283336, #362b28);
    width: 20%;
    height: 20%;
    border-style: solid;
    border-color: #282120;
    border-radius: 50%;
  }
</style>
