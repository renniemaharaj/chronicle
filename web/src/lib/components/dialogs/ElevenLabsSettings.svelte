<script lang="ts">
  import { onMount } from 'svelte';
  import { Input, Label, Select } from 'flowbite-svelte';
  import Button from '../ui/Button.svelte';
  import { speech } from '../../speech/preferences.svelte';
  import { ElevenLabsAccount } from '../../speech/elevenlabs.svelte';
  import ElevenLabsUsage from './ElevenLabsUsage.svelte';

  const account = new ElevenLabsAccount();
  let apiKey = $state('');
  const managed = $derived(account.settings.keySource === 'environment/config');
  const selected = $derived(account.voices.find((voice) => voice.voice_id === speech.voice));
  onMount(() => {
    void account.load();
  });

  async function saveKey(remove = false) {
    const value = remove ? '' : apiKey;
    apiKey = '';
    await account.save(value);
  }
</script>

<p>ElevenLabs receives narration when generating missing audio. Saved passages replay locally.</p>
<p role="status">
  {#if !account.checked}
    Checking server key…
  {:else if account.settings.hasApiKey}
    API key detected: {managed ? 'environment or server configuration' : 'saved in this workspace'}.
  {:else}
    No API key configured. Set ELEVENLABS_API_KEY on the server or save a key below.
  {/if}
</p>
{#if !managed}
  <Label>
    ElevenLabs API key
    <Input
      type="password"
      autocomplete="new-password"
      bind:value={apiKey}
    />
  </Label>
  <p>The key is saved on the server. It is never returned to the browser or stored in browser preferences.</p>
  <div class="dialog-actions">
    <Button
      disabled={account.loading || !apiKey.trim()}
      onclick={() => saveKey()}
    >
      Save API key
    </Button>
    {#if account.settings.hasApiKey}
      <Button
        disabled={account.loading}
        onclick={() => saveKey(true)}
      >
        Remove saved key
      </Button>
    {/if}
  </div>
{/if}
<Button
  disabled={account.loading}
  onclick={() => account.load()}
>
  {account.loading ? 'Loading ElevenLabs…' : 'Refresh voices, models & usage'}
</Button>
<Label>
  Voice
  <Select
    placeholder=""
    aria-label="ElevenLabs voice"
    bind:value={speech.voice}
    onchange={() => speech.save()}
  >
    <option value="">Choose an ElevenLabs voice</option>
    {#if speech.voice && !selected}
      <option value={speech.voice}>Previously selected voice · {speech.voice}</option>
    {/if}
    {#each account.voices as voice}
      <option value={voice.voice_id}>
        {voice.name} · {voice.labels?.accent ?? voice.category}
      </option>
    {/each}
  </Select>
</Label>
{#if selected?.description}
  <p>{selected.description}</p>
{/if}
{#if account.checked && account.settings.hasApiKey && !account.loading && !account.voices.length && !account.error}
  <p>No voices were returned for this key.</p>
{/if}
<Label>
  Model
  <Select
    placeholder=""
    aria-label="ElevenLabs model"
    bind:value={speech.model}
    onchange={() => speech.save()}
  >
    {#if !account.models.some((model) => model.model_id === speech.model)}
      <option value={speech.model}>{speech.model}</option>
    {/if}
    {#each account.models as model}
      <option value={model.model_id}>{model.name}</option>
    {/each}
  </Select>
</Label>
<ElevenLabsUsage
  usage={account.usage}
  error={account.usageError}
/>
{#if account.error}
  <p
    role="alert"
    class="speech-error"
  >
    {account.error}
  </p>
{/if}
