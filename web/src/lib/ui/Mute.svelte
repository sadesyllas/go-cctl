<script lang="ts">
  import { mdiVolumeHigh, mdiVolumeOff } from '@mdi/js';

  export let muted: boolean;
</script>

<div class="block" on:click>
  <svg viewBox="0 0 24 24" width="50" height="50" class:muted>
    <path d={muted ? mdiVolumeOff : mdiVolumeHigh} />
  </svg>
</div>

<style lang="postcss">
  svg {
    @apply cursor-pointer;
    fill: #007fff;
  }

  .muted {
    fill: #f5c6cb;
  }
</style>
