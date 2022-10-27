<script lang="ts">
  export let heading: string;
</script>

<tile>
  <h2>{heading}</h2>
  <slot />
</tile>

<style>
  tile {
    display: flex;
    gap: 0.5em;
    flex-direction: column;
    background-color: #333;
    padding: 1em;
  }
  h2 {
    text-align: left;
    font-weight: 400;
    border-bottom: solid 0.1em #555;
    padding-bottom: 0.2em;
    margin: 0;
  }
  p {
    margin: 0;
  }
</style>
