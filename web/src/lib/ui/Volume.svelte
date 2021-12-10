<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  import Mute from './Mute.svelte';

  export let id = '';
  export let value: number;
  export let muted: boolean;

  const dispatch = createEventDispatcher();

  let hovered = 0;

  $: displayedValue = hovered === 0 ? value : hovered;
</script>

<div {id} class="max-w-min flex flex-col gap-2 align-center">
  <div class="flex">
    {#each new Array(100) as _, index}
      <button
        class="value"
        class:active={index + 1 <= value}
        class:hovering={hovered !== 0}
        class:hovered={index + 1 <= hovered}
        on:mouseover={() => (hovered = index + 1)}
        on:focus={() => (hovered = index + 1)}
        on:mouseout={() => (hovered = 0)}
        on:blur={() => (hovered = 0)}
        on:click={() => dispatch('value', { value: index + 1 })}
      />
    {/each}
  </div>
  <div class="flex items-center justify-between">
    <div class="value-display">{displayedValue}%</div>
    <Mute {muted} on:click={() => dispatch(muted ? 'unmute' : 'mute')} />
  </div>
</div>

<style lang="postcss">
  .value {
    @apply inline-block rounded-none border-t border-b;
    width: 50px;
    height: 50px;
    border-color: #007fff;
  }

  .value:nth-child(1) {
    @apply rounded-l-full border-l;
  }

  .value:nth-child(100) {
    @apply rounded-r-full border-r;
  }

  .value.active {
    background-color: #007fff;
  }

  /* .value.active:not(:nth-child(100)),
  .value.hovered:not(:nth-child(100)) {
    @apply border-r;
    border-right-color: white;
  } */

  .value.hovering:not(.hovered) {
    background-color: transparent;
  }

  .value.hovered {
    background-color: #007fff;
  }

  .value-display {
    @apply rounded-full max-w-min px-4 py-2 pointer-events-none text-2xl text-white text-center font-bold mr-2;
    background-color: #007fff;
    min-width: 4rem;
  }
</style>
