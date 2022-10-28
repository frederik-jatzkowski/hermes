<script lang="ts">
  export let value: string;
  export let name: string;
</script>

<label for={name}><slot /></label>
<input type="password" {name} bind:value />

<style>
  label {
    display: block;
    text-align: left;
  }
  input {
    display: block;
    background-color: #555;
    color: #eee;
    border: none;
    max-width: fit-content;
  }
  input:active,
  input:focus,
  input:hover {
    outline: none;
    border: none;
  }
  input:active,
  input:focus {
    background-color: #444;
  }
</style>
