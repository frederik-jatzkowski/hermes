<script lang="ts">
  import Tile from "../util/Tile.svelte";
  import { authenticate } from "./auth";

  let username: string;
  let password: string;
</script>

<Tile heading="Login">
  <p>For access to the admin panel you have to log in first.</p>
  <label for="username">Username:</label>
  <input type="text" name="username" bind:value={username} />
  <label for="password">Password:</label>
  <input type="password" name="password" bind:value={password} />
  <button on:click={() => authenticate(username, password)}>Authorize</button>
</Tile>

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
  button {
    display: block;
    cursor: pointer;
    background-color: #555;
    color: #eee;
    outline: none;
    border: none;
  }
  input:active,
  input:focus,
  button:hover {
    outline: none;
    border: none;
    background-color: #444;
  }
</style>
