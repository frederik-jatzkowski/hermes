<script lang="ts">
  import Login from "./auth/Login.svelte";
  import { authenticated } from "./auth/auth";
  import ConfigPanel from "./config/ConfigPanel.svelte";
</script>

<main>
  {#if !$authenticated}
    <Login />
  {:else}
    <ConfigPanel />
  {/if}
</main>
<footer>
  © Frederik L. Jatzkowski | 2021 - {new Date().getFullYear()}
</footer>

<style>
  main {
    text-align: center;
    padding: 2rem;
    max-width: 240px;
    margin: 0 auto;
    display: flex;
    gap: 2rem;
    align-items: center;
    justify-content: center;
    flex-wrap: wrap;
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
  footer {
    position: fixed;
    right: 0rem;
    bottom: 0rem;
    left: 0rem;
    padding: 0.5rem;
    background-color: #222;
  }
</style>
