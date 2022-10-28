<script lang="ts">
  import { ColorScheme, ERR, INFO } from "./colors";

  export let scheme: ColorScheme = INFO;
</script>

<p style={`color:${scheme.color}; background-color:${scheme.bg};`}><slot /></p>

<style>
  p {
    display: block;
    text-align: justify;
    padding: 1rem;
    margin: 0;
  }
</style>
