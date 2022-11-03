<script lang="ts">
  import Group from "../util/Group.svelte";
  import type { ConfigType } from "./types";
  import Gateway from "./Gateway.svelte";
  import Button from "../util/Button.svelte";
  import Message from "../util/Message.svelte";
  import { SUCCESS } from "../util/colors";

  export let config: ConfigType;

  function addGateway() {
    config.gateways.unshift({
      address: "",
      services: [],
    });
    config = config;
  }
</script>

<Message>
  This configuration was originally applied on
  {new Date(config.unix * 1000).toUTCString()}
</Message>

<Group>
  <Button scheme={SUCCESS} on:click={addGateway}>Add Gateway</Button>
</Group>

{#each config.gateways as gateway}
  <Gateway {gateway} bind:parent={config} />
{/each}
