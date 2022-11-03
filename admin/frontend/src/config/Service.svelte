<script lang="ts">
  import Button from "../util/Button.svelte";
  import { ERR, SUCCESS } from "../util/colors";
  import Group from "../util/Group.svelte";
  import HeaderInput from "../util/HeaderInput.svelte";
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
      address: "",
    });
    service = service;
  }

  // toggle dropdown
  let open: boolean = false;
  function toggleDropdown() {
    open = !open;
  }
</script>

<service class:open>
  <opener on:click={toggleDropdown} />
  <HeaderInput
    name="hostname"
    bind:value={service.hostName}
    placeholder="Enter the service host name here"
  />
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
    position: relative;
    flex-direction: column;
    padding: 1rem;
    padding-right: 0rem;
    padding-bottom: 0rem;
    border: solid 0.2rem #222;
    border-right: none;
    border-bottom: none;
    height: 2.8rem;
    overflow: hidden;
  }
  service.open {
    height: fit-content;
  }
  service opener {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    height: 2.8rem;
    width: 2.8rem;
    cursor: pointer;
  }
  service opener::before {
    content: "";
    position: absolute;
    inset: 0;
    transition: all 0.2s ease;
    background: #ddd;
    clip-path: polygon(50% 65%, 25% 40%, 75% 40%);
  }
  service opener:hover::before {
    background: #ccc;
    clip-path: polygon(50% 62%, 25% 43%, 75% 43%);
  }
  service.open opener::before {
    clip-path: polygon(50% 35%, 25% 60%, 75% 60%);
  }
  service.open opener:hover::before {
    background: #ccc;
    clip-path: polygon(50% 38%, 25% 57%, 75% 57%);
  }
</style>
