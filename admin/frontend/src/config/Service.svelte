<script lang="ts">
  import Button from "../util/Button.svelte";
  import { ERR, SUCCESS } from "../util/colors";
  import Group from "../util/Group.svelte";
  import Textfield from "../util/Textfield.svelte";
  import Server from "./Server.svelte";
  import type { GatewayType, ServiceType } from "./types";

  export let service: ServiceType;
  export let parent: GatewayType;

  function deleteService() {
    parent.services.splice(parent.services.indexOf(service), 1);
    parent = parent;
  }
  function addServer() {
    service.balancer.servers.unshift({
      address: "Enter the server address here.",
    });
    service = service;
  }
</script>

<service>
  <Textfield name="host" bind:value={service.hostName}>Service Hostname:</Textfield>
  <Textfield name="host" bind:value={service.balancer.algorithm}>Load Balancer Algorithm:</Textfield
  >
  <Group>
    <Button scheme={ERR} on:click={deleteService}>Delete Service</Button>
    <Button scheme={SUCCESS} on:click={addServer}>Add Server</Button>
  </Group>
  {#each service.balancer.servers as server}
    <Server {server} bind:parent={service} />
  {/each}
</service>

<style>
  service {
    display: flex;
    flex-direction: column;
    padding: 1rem;
    padding-right: 0rem;
    padding-bottom: 0rem;
    border: solid 0.2rem #222;
    border-right: none;
    border-bottom: none;
  }
</style>
