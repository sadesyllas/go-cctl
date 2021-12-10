<script lang="ts" context="module">
  import { getDevices, setDefault, setProfile, setVolume, toggleMute } from '$lib/audio';

  export async function load(): Promise<unknown> {
    await getDevices();

    return {};
  }
</script>

<script lang="ts">
  import { devices } from '$lib/audio';
  import { AudioDeviceBus, BluetoothAudioDeviceProfile, bluetoothAudioDeviceProfileToString } from '$lib/audio/types';

  import Volume from '$lib/ui/Volume.svelte';

  $: bluetoothCard = $devices?.cards.filter((card) => card.bus === AudioDeviceBus.Bluetooth)[0];
  $: bluetoothCardProfiles = bluetoothCard?.profiles || [];
  $: bluetoothCardActiveProfile = bluetoothCard?.activeProfile;
  $: sources = $devices?.sources || [];
  $: defaultSource = sources.find((source) => source.isDefault);
  $: defaultSourceIndex = defaultSource?.index;
  $: sinks = $devices?.sinks || [];
  $: defaultSink = sinks.find((sink) => sink.isDefault);
  $: defaultSinkIndex = defaultSink?.index;

  async function onBluetoothCardProfileChange(event: Event) {
    const profile = <BluetoothAudioDeviceProfile>Number((<HTMLSelectElement>event.target).value);

    await setProfile(bluetoothCard.index, profile)
      .then(() => getDevices())
      .then(async () => {
        if (profile === BluetoothAudioDeviceProfile.HeadsetHeadUnit) {
          const bluetoothSource = sources.find((source) => source.bluetoothProtocol);

          await setDefault('source', bluetoothSource.index, bluetoothSource.name);

          const bluetoothSink = sinks.find((sink) => sink.bluetoothProtocol);

          await setDefault('sink', bluetoothSink.index, bluetoothSink.name);
        }
      })
      .then(() => getDevices());
  }

  async function onDefaultAudioDeviceChange(type: 'source' | 'sink', event: Event) {
    const index = Number((<HTMLSelectElement>event.target).value);
    const name = (type === 'source' ? sources : sinks).find((target) => target.index === index).name;

    await setDefault(type, index, name).then(() => getDevices());
  }

  async function onVolumeChange(type: 'source' | 'sink', volume: number) {
    const target = type === 'source' ? defaultSource : defaultSink;

    await setVolume(type, target.index, volume).then(() => getDevices());
  }

  async function onMuteToggle(type: 'source' | 'sink', mute: boolean) {
    const target = type === 'source' ? defaultSource : defaultSink;

    await toggleMute(type, target.index, mute).then(() => getDevices());
  }
</script>

<main class="p-10">
  {#if bluetoothCard}
    <section>
      <label for="bluetooth-card">Bluetooth Audio Card Profiles</label>

      <select
        id="bluetooth-card"
        on:change={async (e) => await onBluetoothCardProfileChange(e)}
        bind:value={bluetoothCardActiveProfile}
      >
        {#each bluetoothCardProfiles as bluetoothCardProfile}
          <option value={bluetoothCardProfile}>
            {bluetoothCard.description} [{bluetoothAudioDeviceProfileToString(bluetoothCardProfile)}]
          </option>
        {/each}
      </select>
    </section>
  {/if}

  <section>
    <label for="sources">Audio Sources</label>

    <select
      id="sources"
      on:change={async (e) => await onDefaultAudioDeviceChange('source', e)}
      bind:value={defaultSourceIndex}
    >
      {#each sources as source}
        <option value={source.index}>{source.description}</option>
      {/each}
    </select>
  </section>

  <section>
    <label for="sinks">Audio Sinks</label>

    <select
      id="sinks"
      on:change={async (e) => await onDefaultAudioDeviceChange('sink', e)}
      bind:value={defaultSinkIndex}
    >
      {#each sinks as sink}
        <option value={sink.index}>{sink.description}</option>
      {/each}
    </select>
  </section>

  <section>
    <label for="input-volume">Input Volume</label>

    <Volume
      id="input-volume"
      value={defaultSource?.volume}
      muted={defaultSource?.isMuted}
      on:value={async ({ detail: { value } }) => await onVolumeChange('source', value)}
      on:mute={async () => await onMuteToggle('source', true)}
      on:unmute={async () => await onMuteToggle('source', false)}
    />
  </section>

  <section>
    <label for="output-volume">Output Volume</label>

    <Volume
      id="output-volume"
      value={defaultSink?.volume}
      muted={defaultSink?.isMuted}
      on:value={async ({ detail: { value } }) => await onVolumeChange('sink', value)}
      on:mute={async () => await onMuteToggle('sink', true)}
      on:unmute={async () => await onMuteToggle('sink', false)}
    />
  </section>
</main>

<style lang="postcss">
  section {
    @apply mb-4;
  }

  select {
    @apply px-4 py-2 text-2xl font-bold rounded-full outline-none block appearance-none text-white w-full;
    background-color: #007fff;
  }
</style>
