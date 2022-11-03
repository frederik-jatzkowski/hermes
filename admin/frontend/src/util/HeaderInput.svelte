<script lang="ts">
  export let value: string;
  export let name: string;
  export let placeholder: string;
</script>

<input type="text" {name} bind:value {placeholder} />

<style>
  input::placeholder {
    color: #bbb;
  }
  input {
    display: block;
    background-color: unset;
    color: #eee;
    border: none;
    width: 50%;
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
