<script lang="ts">
  import { onMount } from 'svelte';
  import { connectAudioWS } from '$lib/audio';
  import '../app.postcss';

  // async function loopGetDevices() {
  //   await getDevices();

  //   setTimeout(loopGetDevices, 5_000);
  // }

  onMount(connectAudioWS);
</script>

<slot />
