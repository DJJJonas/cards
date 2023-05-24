<script lang="ts">
  import Board from "./lib/Board.svelte";
  import { type Board as B, type Event } from "./interfaces/cards";

  const host = import.meta.env.DEV ? "localhost:8123" : window.location.host;
  const protocol = host.startsWith("localhost") ? "ws://" : "wss://";
  const gameWSUrl = protocol + host + "/ws/connect";
  const socket = new WebSocket(gameWSUrl);

  let board: B;
  let toggleConsole: boolean = true;
  let debugmsg: string = "";

  socket.addEventListener("message", (event) => {
    const msg = JSON.parse(event.data);
    switch (msg.type) {
      case "board":
        board = msg.data;
        board.lastEvents.forEach((e) => handleHistoricEvent(board, e));
        break;
      case "error":
        alert(msg.data);
        break;
      case "result":
        alert(msg.data);
        break;
    }
  });

  function handleHistoricEvent(b: B, e: Event) {
    if (!e) return;
    switch (e.type) {
      case "draw":
        if (e.heroId === b.myHero.id && e.source.id === b.myHero.id) {
          log(`You drew ${e.target.name}`);
        } else if (e.heroId === b.enemyHero.id) {
          log(`Enemy drew a card`);
        }
        log(`${e.source.name} drew ${e.target.name}`);
        break;
      case "endturn":
        if (e.heroId === b.myHero.id) {
          log(`Your turn is done`);
        } else if (e.heroId === b.enemyHero.id) {
          log(`Enemy's turn is done`);
        }
        break;
      case "heropower":
        if (e.heroId === b.myHero.id) {
          let msg = `You used your hero power ${e.source.name}`;
          msg += e.target ? " on " + e.target.name : "";
          log(msg);
        } else if (e.heroId === b.enemyHero.id) {
          let msg = `Enemy used their hero power ${e.source.name}`;
          msg += e.target ? " on " + e.target.name : "";
          log(msg);
        }
        break;
      case "attack":
        log(`${e.source.name} attacked ${e.target.name}`);
        break;
      case "play":
        if (e.heroId === b.myHero.id) {
          log(`You played ${e.source.name}`);
        }
        break;
      case "summon":
        log(`${e.source.name} summoned ${e.target.name}`);
      case "secret":
        // TODO: rewrite this after secret rework to show the secret's name
        log(`${e.source.name} revealed a secret`);
    }
  }

  function log(msg: string) {
    debugmsg += msg + "\n";
  }
</script>

<main>
  <button
    class="console-toggle"
    on:click={() => (toggleConsole = !toggleConsole)}>Toggle Logs</button
  >
  {#if toggleConsole}
    <button class="console-clear" on:click={() => (debugmsg = "")}>Clear</button
    >
    <textarea class="console" bind:value={debugmsg} />
  {/if}
  {#if board}
    {#key board}
      <Board {board} {socket} {log} />
    {/key}
  {:else}
    Loading...
  {/if}
</main>

<style>
  button {
    cursor: pointer;
    border-width: 0;
    background-color: #101010;
    color: ghostwhite;
  }
  .console-toggle {
    position: absolute;
    top: 16%;
    left: 1%;
  }
  .console-clear {
    position: absolute;
    top: 16%;
    left: 6%;
  }

  .console {
    color: ghostwhite;
    position: absolute;
    top: 20%;
    left: 1%;
    width: 20%;
    height: 40%;
    background-color: #10101060;
    border-width: 0;
  }
  .console::-webkit-scrollbar {
    width: 8px;
  }

  .console::-webkit-scrollbar-track {
    background-color: #f1f1f1;
  }

  .console::-webkit-scrollbar-thumb {
    background-color: #888;
  }

  .console::-webkit-scrollbar-thumb:hover {
    background-color: #555;
  }
</style>
