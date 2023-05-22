<script setup lang="ts">
  import { createEventDispatcher } from "svelte";
  import { type Card, getCardAttack, getMaxHealth } from "../interfaces/cards";
  export let card: Card;
  export let width: number;
  const dispatch = createEventDispatcher();

  function attackColor() {
    let enchAttack = getCardAttack(card);
    if (card.attack === enchAttack) return "#E9EAE0";
    if (card.attack > enchAttack) return "#DB1F48";
    if (card.attack < enchAttack) return "#01949A";
  }

  function healthColor() {
    let maxHealth = card.maxHealth;
    const enchmaxh = getMaxHealth(card);
    if (card.maxHealth !== enchmaxh) {
      maxHealth = enchmaxh;
    }
    if (card.health === maxHealth) return "#E9EAE0";
    if (card.health < maxHealth) return "#DB1F48";
    if (card.health > maxHealth) return "#01949A";
  }
</script>

<div
  class="minion"
  style:width={width + "px"}
  class:targeted={card.targeted}
  class:playable={card.playable}
  class:selected={card.selected}
  class:CanAttack={card.canAttack}
  class:Sleeping={card.sleeping}
  on:click={() => dispatch("click")}
  on:keypress={() => dispatch("keypress")}
  on:mouseenter={() => dispatch("mouseenter")}
  on:mouseleave={() => dispatch("mouseleave")}
  on:focus={() => dispatch("focus")}
>
  <img
    class="image"
    style:width={width + "px"}
    src={card.image}
    alt="card portrait"
  />
  <div
    class="attack"
    style:color={attackColor()}
    style:width={width * 0.2 + "px"}
    style:height={width * 0.2 + "px"}
    style:font-size={width * 0.16 + "px"}
  >
    {getCardAttack(card)}
  </div>
  <div
    class="health"
    style:color={healthColor()}
    style:width={width * 0.2 + "px"}
    style:height={width * 0.2 + "px"}
    style:font-size={width * 0.16 + "px"}
  >
    {card.health}
  </div>
  {#if card.sleeping}
    <div class="sleeping" style:font-size={width * 0.7 + "px"}>💤</div>
  {/if}
</div>

<style>
  .CanAttack {
    cursor: pointer;
    border-color: green;
    box-shadow: 0 0 20px green;
  }
  .Sleeping {
    filter: grayscale(80%);
  }
  .selected {
    border-color: white;
    box-shadow: 0 0 30px white;
  }
  .targeted {
    cursor: pointer;
    border-color: orangered;
    box-shadow: 0 0 24px orangered;
  }
  .playable {
    box-shadow: 0 0 8px laightgreen;
  }
  .minion {
    z-index: -1;
    font-family: "Righteous";
    position: relative;
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    aspect-ratio: 1/1.3;
    margin-left: 0.2%;
    margin-right: 0.2%;
    border-style: solid;
    border-width: 2px;
    border-radius: 50%;
    background-size: cover;
    text-shadow: -1px -1px 1px black, 1px -1px 1px black, -1px 1px 1px black,
      1px 1px 1px black;
  }

  .attack {
    z-index: 1;
    background-color: yellow;
    text-align: center;
    font-weight: bold;
    border-style: solid;
    border-width: 1px;
    border-radius: 50%;
    background-image: linear-gradient(#ffed8a, #ffe55c);
  }

  .health {
    z-index: 1;
    background-color: red;
    width: 15%;
    height: 15%;
    text-align: center;
    font-weight: bolder;
    border-style: solid;
    border-width: 1px;
    border-radius: 50%;
    background-image: linear-gradient(#e7625f, #c85250);
  }
  .image {
    position: absolute;
    object-fit: cover;
    top: 50%;
    left: 50%;
    translate: -50% -50%;
    aspect-ratio: 1/1.3;
    border-style: solid;
    border-width: 0;
    border-radius: 50%;
  }
  .sleeping {
    position: absolute;
    left: 50%;
    top: 50%;
    translate: -50% -50%;
    animation: sleeping 0.7s infinite;
    animation-direction: alternate;
  }

  @keyframes sleeping {
    100% {
      scale: 1.1;
    }
  }
</style>
