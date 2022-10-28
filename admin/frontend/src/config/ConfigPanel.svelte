<script lang="ts">
  import Button from "../util/Button.svelte";
  import Group from "../util/Group.svelte";
  import { ERR, SUCCESS } from "../util/colors";
  import Spinner from "../util/Spinner.svelte";
  import Tile from "../util/Tile.svelte";
  import { ApplyConfig, ConfigHistory } from "./fetch";
  import { ConfigState } from "./states";
  import type { ConfigHistoryType } from "./types";
  import Config from "./Config.svelte";
  import Message from "../util/Message.svelte";

  let state: ConfigState = ConfigState.Loading;
  let history: ConfigHistoryType;
  let current: number;
  let errors: string[] = [];

  async function loadConfig() {
    history = await ConfigHistory();
    current = history.length - 1;
    state = ConfigState.Edit;
  }

  async function applyConfig() {
    state = ConfigState.Applying;
    errors = [];
    const res = await ApplyConfig(history[current]);
    if (res.ok) {
      state = ConfigState.Success;

      return;
    }

    state = ConfigState.Edit;
    errors = res.exceptions;
  }

  loadConfig();
</script>

<Tile heading="Configuration">
  {#each errors as error}
    <Message scheme={ERR}>{error}</Message>
  {/each}
  {#if state == ConfigState.Loading}
    <Spinner>Loading configuration history...</Spinner>
  {:else if state == ConfigState.Edit}
    <Group>
      {#if current > 0}
        <Button on:click={() => current--}>
          Previous ({new Date(history[current - 1].unix).toUTCString()})
        </Button>
      {/if}
      {#if current < history.length - 1}
        <Button on:click={() => current++}>
          Next ({new Date(history[current + 1].unix).toUTCString()})
        </Button>
      {/if}
      <Button scheme={ERR} on:click={applyConfig}>Apply This Configuration</Button>
    </Group>
    <Config config={history[current]} />
  {:else if state == ConfigState.Applying}
    <Spinner>Applying configuration...</Spinner>
  {:else if state == ConfigState.Success}
    <Message scheme={SUCCESS}>Successfully applied configuration!</Message>
    <Button on:click={loadConfig}>Back to the Editor.</Button>
  {/if}
</Tile>
