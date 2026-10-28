<script lang="ts">
  import { Input, Label, Select } from 'flowbite-svelte';
  import { speech } from '../../speech/preferences.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import ElevenLabsSettings from './ElevenLabsSettings.svelte';
  import GoogleVoiceSettings from './GoogleVoiceSettings.svelte';
  import BrowserVoiceSettings from './BrowserVoiceSettings.svelte';
  import NarrationLibrary from './NarrationLibrary.svelte';
</script>

<div class="speech-settings">
  <p>Read the evidence-linked story shown in the narrative panel. Changes apply to the current passage.</p>
  <Label>
    Speech engine
    <Select
      placeholder=""
      aria-label="Speech engine"
      bind:value={speech.provider}
      onchange={() => speech.chooseProvider()}
    >
      <option value="browser">Browser voices</option>
      <option value="google">Google Cloud Text-to-Speech</option>
      <option value="elevenlabs">ElevenLabs</option>
      <option value="screen">Read on screen · no audio</option>
    </Select>
  </Label>
  <Label>
    Language
    <Input
      bind:value={speech.language}
      onchange={() => speech.save()}
      placeholder="en-US"
    />
  </Label>
  {#if speech.provider === 'browser'}
    <BrowserVoiceSettings />
  {:else if speech.provider === 'google'}
    <GoogleVoiceSettings />
  {:else if speech.provider === 'elevenlabs'}
    <ElevenLabsSettings />
  {/if}
  {#if speech.provider === 'google' || speech.provider === 'elevenlabs'}
    <NarrationLibrary provider={speech.provider} />
  {/if}
  <Label>
    Reading speed
    <Select
      placeholder=""
      aria-label="Reading speed"
      bind:value={ws.speed}
    >
      <option value={0.75}>0.75×</option>
      <option value={1}>1×</option>
      <option value={1.25}>1.25×</option>
      <option value={1.5}>1.5×</option>
      <option value={2}>2×</option>
    </Select>
  </Label>
  {#if speech.error}
    <p
      role="alert"
      class="speech-error"
    >
      {speech.error}
    </p>
  {/if}
  <p role="status">{speech.status}</p>
</div>
