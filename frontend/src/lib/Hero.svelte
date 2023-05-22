<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { getMaxHealth, type Card } from "../interfaces/cards";
  export let card: Card;
  export let width = 300;
  let dispatcher = createEventDispatcher();
  function healthColor() {
    let maxhealth = getMaxHealth(card);
    if (card.health === maxhealth) return "#E9EAE0";
    if (card.health < maxhealth) return "#DB1F48";
    if (card.health > maxhealth) return "#01949A";
  }
</script>

<div
  class="hero"
  style:width={width + "px"}
  style:height={width + "px"}
  class:targeted={card.targeted}
  class:playable={card.playable}
  class:selected={card.selected}
  class:CanAttack={card.canAttack}
  class:Sleeping={card.sleeping}
  on:click={() => dispatcher("click")}
  on:keypress={() => dispatcher("keypress")}
>
  <img src={card.image} alt="hero" />
  <div
    class="health"
    style:color={healthColor()}
    style:font-size={width * 0.14 + "px"}
  >
    {card.health}
  </div>
  {#if card.shield > 0}
    <div class="shield" style:font-size={width * 0.15 + "px"}>
      {card.shield}
    </div>
  {/if}
</div>

<style>
  .CanAttack {
    cursor: pointer;
    box-shadow: 0 0 20px green;
    border-color: "green";
  }
  .Sleeping {
    filter: grayscale(80%);
  }
  .selected {
    box-shadow: 0 0 30px white;
    border-color: "white";
  }
  .targeted {
    cursor: pointer;
    box-shadow: 0 0 12px orangered;
  }
  .playable {
    box-shadow: 0 0 8px laightgreen;
  }
  .hero {
    position: relative;
    text-shadow: 1px 1px 2px black, -1px 1px 2px black, 1px -1px 2px black,
      -1px -1px 2px black;
    text-align: center;
    border-style: none;
    border-style: solid;
    border-color: black;
    border-radius: 50% 50% 0 0;
  }
  .hero > img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 50% 50% 0 0;
  }
  .health {
    position: absolute;
    bottom: -10%;
    left: -4%;
    width: 18%;
    height: 18%;
    border-style: solid;
    border-radius: 50%;
    border-color: black;
    background-image: linear-gradient(to top, #c85250, #e7625f);
  }
  .shield {
    position: absolute;
    bottom: -10%;
    left: 86%;
    width: 18%;
    height: 18%;
    border-style: solid;
    border-radius: 50%;
    border-color: black;
    background-image: linear-gradient(to top, #494b49, #5f635f);
  }
</style>
