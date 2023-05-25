<script lang="ts">
  import Card from "./Card.svelte";
  import Minion from "./Minion.svelte";
  import {
    type Action,
    type Board,
    type Card as ICard,
    getTags,
    getCost,
  } from "../interfaces/cards";
  import Hero from "./Hero.svelte";
  import Heropower from "./Heropower.svelte";
  import Weapon from "./Weapon.svelte";

  export let board: Board;
  export let socket: WebSocket;
  export let log;

  const turnLabel = board.myTurn ? "Pass turn" : "Enemy turn";

  let fullCardView: ICard | undefined = undefined;
  let width = window.innerWidth * 0.09;
  let nextAction: Action | null = null;

  board.myHand.forEach((card: ICard) => {
    card.playable = playable(card);
  });
  board.myHeroPower.playable = playable(board.myHeroPower);

  function update() {
    board = board;
  }

  function playable(card: ICard): boolean {
    if (hasBloodPayment(card))
      if (board.myHero.health > getCost(card)) return true;
      else return false;
    return board.myMana >= getCost(card);
  }

  function isAlly(card: ICard): boolean {
    if (card.id === board.myHero.id) return true;
    if (card.id === board.myWeapon?.id) return true;
    if (board.myMinions.some((c) => c.id === card.id)) return true;
    return false;
  }

  function hasBloodPayment(card: ICard): boolean {
    return getTags(card).some((t) => t === "bloodPayment")!;
  }

  function useHeroPower(): void {
    if (!board.myTurn || !playable(board.myHeroPower)) return;
    send({
      Type: "heropower",
      SourceId: "",
      TargetId: "",
      Position: 0,
    });
  }

  function prepareAttack(card: ICard): void {
    if (selectedAllyToAttack(card)) {
      card.selected = true;
      nextAction = {
        Type: "attack",
        SourceId: card.id,
        TargetId: "",
        Position: 0,
      };
      targetForAttack(card);
      return;
    } else if (selectedEnemyToAttack(card) || selectedTargetNeededCard(card)) {
      nextAction.TargetId = card.id;
      send(nextAction);
    }
    nextAction = null;
    untargetAll();
  }

  function selectedAllyToAttack(card: ICard) {
    return nextAction === null && isAlly(card) && card.canAttack;
  }

  function selectedEnemyToAttack(card: ICard) {
    return nextAction !== null && !isAlly(card) && card.targeted;
  }

  function selectedTargetNeededCard(card: ICard) {
    return (
      nextAction !== null && nextAction.SourceId !== card.id && card.targeted
    );
  }

  function playFromHand(pos: number, card: ICard) {
    if (!board.myTurn || !playable(card)) return;
    if (hasTargets(card)) {
      card.targeted = true;
      log("select a target");
      targetIds(card.targets);
      log(card.targets);
      nextAction = {
        Type: "play",
        SourceId: card.id,
        TargetId: "",
        Position: pos,
      };
      return;
    }
    send({
      Type: "play",
      SourceId: card.id,
      TargetId: "",
      Position: pos,
    });
  }

  function endTurn(): void {
    if (!board.myTurn) return;
    send({
      Type: "endturn",
      SourceId: "",
      TargetId: "",
      Position: 0,
    });
  }

  function send(data: Action) {
    socket.send(JSON.stringify(data));
  }

  function targetForAttack(card: ICard): void {
    let targets = [board.enemyHero, ...board.enemyMinions];
    if (getTags(card).some((t) => "rush") && card.canAttack && card.sleeping) {
      targets.shift();
    }
    let taunt = taunts(targets);
    if (taunt.length > 0) target(taunt);
    else target(targets);
    update();
  }

  function taunts(cards: ICard[]) {
    return cards.filter((c) => getTags(c).some((t) => t === "taunt"));
  }

  function target(cards: ICard[]) {
    cards.forEach((c) => {
      c.targeted = true;
    });
    update();
  }

  function targetIds(ids: string[]) {
    [
      board.myHero,
      ...board.myMinions,
      board.enemyHero,
      ...board.enemyMinions,
    ].forEach((c) => {
      if (ids.some((id) => id === c.id)) c.targeted = true;
    });
    update();
  }

  function untargetAll(): void {
    [
      board.myHero,
      ...board.myMinions,
      board.enemyHero,
      ...board.enemyMinions,
    ].forEach((c) => {
      c.targeted = false;
      c.selected = false;
    });
    update();
  }

  function hasTargets(card: ICard): boolean {
    return card.targets !== null && card.targets.length > 0;
  }

  function manaString(
    mana: number,
    maxMana: number,
    maxMaxMana: number
  ): string {
    let d = 0;
    let m = mana;
    let mm = maxMana - mana;
    if (mm < 0) {
      d = maxMana - mana;
      mm = 0;
    }
    let mmm = maxMaxMana - mm + d;
    return "🟦".repeat(m) + "◼️".repeat(mm) + "◾".repeat(mmm);
  }
</script>

<main style="color: white;font-family:'Righteous';">
  <input type="range" min="100" max="300" bind:value={width} />
  {#if board}
    {#if fullCardView}
      <div class="full-card-view">
        <Card card={fullCardView} width={width * 2} />
      </div>
    {/if}

    <div class="my-hero">
      <Hero
        card={board.myHero}
        width={width * 1}
        on:click={() => prepareAttack(board.myHero)}
      />
    </div>

    <div class="enemy-hero">
      <Hero
        card={board.enemyHero}
        width={width * 1}
        on:click={() => prepareAttack(board.enemyHero)}
      />
    </div>

    {#if board.myWeapon}
      <div
        class="my-weapon"
        on:mouseenter={() => (fullCardView = board.myWeapon)}
        on:mouseleave={() => (fullCardView = null)}
      >
        <Weapon card={board.myWeapon} width={width * 0.8} />
      </div>
    {/if}

    {#if board.enemyWeapon}
      <div
        class="enemy-weapon"
        on:mouseenter={() => (fullCardView = board.enemyWeapon)}
        on:mouseleave={() => (fullCardView = null)}
      >
        <Weapon card={board.enemyWeapon} width={width * 0.8} />
      </div>
    {/if}

    <div
      class="my-heropower"
      on:click={() => useHeroPower()}
      on:keypress={() => useHeroPower()}
      on:mouseenter={() => (fullCardView = board.myHeroPower)}
      on:mouseleave={() => (fullCardView = null)}
    >
      <Heropower card={board.myHeroPower} width={width * 0.65} />
    </div>

    <div
      class="enemy-heropower"
      on:mouseenter={() => (fullCardView = board.enemyHeroPower)}
      on:mouseleave={() => (fullCardView = null)}
    >
      <Heropower card={board.enemyHeroPower} width={width * 0.65} />
    </div>

    {#if board.mySecrets}
      <div class="my-secrets">
        {#each board.mySecrets as secret}
          <div
            class="secret"
            on:mouseenter={() => (fullCardView = secret)}
            on:mouseleave={() => (fullCardView = null)}
          >
            ?
          </div>
        {/each}
      </div>
    {/if}

    {#if board.enemySecretCount > 0}
      <div class="enemy-secrets">
        {#each Array(board.enemySecretCount) as _}
          <div class="secret">?</div>
        {/each}
      </div>
    {/if}

    <div class="my-deck">
      {board.myDeckSize}
    </div>
    <div class="enemy-deck">
      {board.enemyDeckSize}
    </div>

    <div class="my-mana">
      {manaString(board.myMana, board.myMaxMana, board.myMaxMaxMana)}
    </div>

    <div class="enemy-mana">
      <!-- "◼️" -->
      {manaString(board.enemyMana, board.enemyMaxMana, board.enemyMaxMaxMana)}
    </div>

    <div class="minions my-minions">
      {#each board.myMinions as minion (minion.id)}
        <Minion
          bind:card={minion}
          width={width * 0.7}
          on:click={() => prepareAttack(minion)}
          on:mouseenter={() => (fullCardView = minion)}
          on:mouseleave={() => (fullCardView = null)}
        />
      {/each}
    </div>

    <div class="minions enemy-minions">
      {#each board.enemyMinions as minion (minion.id)}
        <Minion
          bind:card={minion}
          width={width * 0.7}
          on:click={() => prepareAttack(minion)}
          on:mouseenter={() => (fullCardView = minion)}
          on:mouseleave={() => (fullCardView = null)}
        />
      {/each}
    </div>

    <div class="my-hand">
      {#each board.myHand as card (card.id)}
        <Card
          {card}
          {width}
          on:click={() => playFromHand(0, card)}
          on:mouseenter={() => (card.selected = true)}
          on:mouseleave={() => (card.selected = false)}
        />
      {/each}
    </div>

    <button
      class="turn-button"
      on:click={() => endTurn()}
      style:background-image={board.myTurn
        ? "linear-gradient(#ffffcc, #f3e5ab)"
        : "linear-gradient(#616151, #464231)"}
      style:width={width * 1.2 + "px"}
      style:font-size={width * 0.2 + "px"}
      >{turnLabel}
    </button>
  {/if}
</main>

<style>
  .my-secrets {
    position: absolute;
    display: flex;
    justify-content: center;
    width: 10%;
    left: 50%;
    bottom: 7%;
    transform: translate(-50%, 50%);
    font-size: 10pt;
  }
  .enemy-secrets {
    position: absolute;
    display: flex;
    justify-content: center;
    width: 10%;
    left: 50%;
    top: 0.5%;
    transform: translate(-50%);
    font-size: 10pt;
  }
  .secret {
    font-size: 18pt;
    text-align: center;
    background-color: yellow;
    text-shadow: 1px 1px 2px black, -1px 1px 2px black, 1px -1px 2px black,
      -1px -1px 2px black;
    width: 16%;
    aspect-ratio: 1;
    border-style: solid;
    border-color: #222222;
    border-radius: 50%;
  }
  .enemy-mana {
    position: absolute;
    top: 10%;
    left: 1%;
  }
  .enemy-minions {
    position: absolute;
    top: 31.5%;
    left: 50%;
    translate: -50%;
    transition: all 0.2s ease;
  }
  .enemy-hero {
    position: absolute;
    top: 4%;
    left: 50%;
    translate: -50%;
    transition: all 0.2s ease;
  }
  .enemy-heropower {
    position: absolute;
    top: 8%;
    left: 57%;
    cursor: pointer;
  }
  .turn-button {
    position: absolute;
    padding: 0;
    margin: 0;
    top: 50%;
    translate: 0 -50%;
    right: 1%;
    cursor: pointer;
  }
  .minions {
    display: flex;
    justify-content: space-around;
    align-items: center;
  }
  .my-minions {
    position: absolute;
    display: flex;
    bottom: 31.5%;
    left: 50%;
    translate: -50%;
    transition: all 0.2s ease;
  }
  .my-heropower {
    position: absolute;
    bottom: 13%;
    left: 57%;
    cursor: pointer;
  }

  .my-weapon {
    position: absolute;
    bottom: 10%;
    left: 36%;
  }

  .enemy-weapon {
    position: absolute;
    top: 7%;
    left: 36%;
  }
  .my-hero {
    position: absolute;
    bottom: 10%;
    left: 50%;
    translate: -50%;
    transition: all 0.2s ease;
  }
  .my-hand {
    position: absolute;
    display: flex;
    flex-direction: row;
    transition: bottom 0.2s ease;
    bottom: -10%;
    left: 50%;
    translate: -50% 50%;
    justify-content: space-around;
    align-items: center;
    transition: all 0.2s ease;
  }
  .my-hand:hover {
    bottom: 20%;
    height: 40%;
  }
  .my-mana {
    position: absolute;
    bottom: 10%;
    left: 1%;
  }

  .my-deck {
    position: absolute;
    font-size: 40px;
    text-align: center;
    background-color: #282828;
    border-style: solid;
    border-width: 4px;
    border-radius: 50%;
    bottom: 20%;
    right: 10%;
    width: 50px;
    height: 50px;
  }
  .enemy-deck {
    position: absolute;
    font-size: 40px;
    text-align: center;
    background-color: #282828;
    border-style: solid;
    border-width: 4px;
    border-radius: 50%;
    top: 20%;
    right: 10%;
    width: 50px;
    height: 50px;
  }

  .full-card-view {
    position: absolute;
    top: 50%;
    left: 14%;
    transform: translate(-50%, -50%);
  }
</style>
