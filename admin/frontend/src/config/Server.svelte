<script lang="ts">
  import Button from "../util/Button.svelte";
  import { ERR } from "../util/colors";
  import Group from "../util/Group.svelte";
  import HeaderInput from "../util/HeaderInput.svelte";
  import type { ServerType, ServiceType } from "./types";

  export let server: ServerType;
  export let parent: ServiceType;

  function deleteServer() {
    parent.balancer.servers.splice(parent.balancer.servers.indexOf(server), 1);
    parent = parent;
  }
</script>

<service>
  <HeaderInput
    name="remoteAddress"
    bind:value={server.address}
    placeholder="Enter the server address here"
  />
  <Button scheme={ERR} on:click={deleteServer}>Delete Server</Button>
</service>

<style>
  service {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
    padding-right: 0rem;
    border: solid 0.2rem #222;
    border-right: none;
    border-bottom: none;
  }
</style>
