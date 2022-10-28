<script lang="ts">
  export let checked: boolean;
  export let name: string;
</script>

<label for={name}><slot /></label>
<input type="checkbox" {name} bind:checked />

<style>
  input {
    margin: 0;
    width: 1.1rem;
    height: 1.1rem;
    cursor: pointer;
  }
</style>
