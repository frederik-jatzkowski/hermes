<script lang="ts">
  import { ColorScheme, INTERACTIVE } from "./colors";

  export let scheme: ColorScheme = INTERACTIVE;
</script>

<button on:click style={`color:${scheme.color}; background-color:${scheme.bg};`}><slot /></button>

<style>
  button {
    display: inline;
    max-width: max-content;
    cursor: pointer;
    transition: all 0.05s ease;
    outline: none;
    border: none;
    padding: 0.5rem 1rem;
  }
  button:active {
    outline: none;
    border: none;
    filter: brightness(80%);
  }
</style>
