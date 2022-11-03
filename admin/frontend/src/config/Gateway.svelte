<script lang="ts">
  import Button from "../util/Button.svelte";
  import { ERR, SUCCESS } from "../util/colors";
  import Group from "../util/Group.svelte";
  import HeaderInput from "../util/HeaderInput.svelte";
  import Textfield from "../util/Textfield.svelte";
  import Service from "./Service.svelte";
  import type { ConfigType, GatewayType } from "./types";

  export let gateway: GatewayType;
  export let parent: ConfigType;

  function deleteGateway() {
    parent.gateways.splice(parent.gateways.indexOf(gateway), 1);
    parent = parent;
  }
  function addService() {
    gateway.services.unshift({
      hostName: "",
      balancer: {
        algorithm: "RoundRobin",
        servers: [],
      },
    });
    gateway = gateway;
  }
</script>

<gateway>
  <HeaderInput
    name="localAddress"
    bind:value={gateway.address}
    placeholder="Enter the gateway address here"
  />
  <Group>
    <Button scheme={ERR} on:click={deleteGateway}>Delete Gateway</Button>
    <Button scheme={SUCCESS} on:click={addService}>Add Service</Button>
  </Group>
  {#each gateway.services as service}
    <Service {service} bind:parent={gateway} />
  {/each}
</gateway>

<style>
  gateway {
    display: flex;
    flex-direction: column;
    padding: 1rem;
    padding-right: 0rem;
    padding-bottom: 0rem;
    border: solid 0.2rem #222;
  }
</style>
