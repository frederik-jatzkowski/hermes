<script lang="ts">
  import Login from "./auth/Login.svelte";
  import { authenticated } from "./auth/auth";
  import Panel from "./panel/Panel.svelte";
</script>

<main>
  {#if !$authenticated}
    <Login />
  {:else}
    <Panel />
  {/if}
</main>

<style>
  main {
    text-align: center;
    padding: 1em;
    max-width: 240px;
    margin: 0 auto;
    display: flex;
    gap: 1em;
    flex-wrap: wrap;
  }

  h1 {
    color: #eee;
    text-transform: uppercase;
    text-align: left;
    font-size: 3em;
    font-weight: 100;
    margin: 0;
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>
