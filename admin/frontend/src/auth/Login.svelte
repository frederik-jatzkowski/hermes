<script lang="ts">
  import Button from "../util/Button.svelte";
  import Textfield from "../util/Textfield.svelte";
  import Password from "../util/Password.svelte";
  import Tile from "../util/Tile.svelte";
  import Message from "../util/Message.svelte";
  import { authenticate } from "./auth";

  let username: string;
  let password: string;
</script>

<Tile heading="Login" maxWidth={30}>
  <Message>For access to the admin panel you have to log in first.</Message>
  <Textfield name="username" bind:value={username}>Username:</Textfield>
  <Password name="password" bind:value={password}>Password:</Password>
  <Button on:click={() => authenticate(username, password)}>Authorize</Button>
</Tile>
