<script lang="ts">
  import { getMaxHealth, type Card } from "../interfaces/cards";
  export let card: Card;
  export let width = 300;
  function healthColor() {
    let maxhealth = getMaxHealth(card);
    if (card.health === maxhealth) return "#E9EAE0";
    if (card.health < maxhealth) return "#DB1F48";
    if (card.health > maxhealth) return "#01949A";
  }
</script>

<div
  class="heropower"
  style:background-image="url({card.image})"
  style:width={width + "px"}
  style:height={width + "px"}
  style={card.playable
    ? "box-shadow: 0 0 8px lightgreen;color:lightgreen;"
    : ""}
>
  <div class="mana" style:font-size={width * 0.23 + "px"}>
    {card.mana}
  </div>
</div>

<style>
  .mana {
    color: white;
    position: absolute;
    top: -10%;
    width: 30%;
    height: 30%;
    text-align: center;
    background-image: linear-gradient(to top, #3434d7, #3b3bd7);
    border-radius: 40%;
  }
  .heropower {
    color: rgb(22, 2, 2);
    position: relative;
    border-style: hidden;
    border-radius: 50%;
    border-style: solid;
    border-width: 2px;
    display: flex;
    justify-content: center;
    background-size: cover;
  }
</style>
