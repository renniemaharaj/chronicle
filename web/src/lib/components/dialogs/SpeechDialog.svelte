<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Input, Label, Select } from 'flowbite-svelte';
  import { onMount } from 'svelte';
  import { speech } from '../../speech/preferences.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  let browserVoices = $state<SpeechSynthesisVoice[]>([]);
  let googleVoices = $state<{ name: string; languageCodes: string[]; ssmlGender: string }[]>([]);
  let loading = $state(false);
  let error = $state('');
  onMount(() => {
    if (!('speechSynthesis' in window)) {
      return;
    }
    const refresh = () => {
      browserVoices = speechSynthesis.getVoices();
    };
    refresh();
    speechSynthesis.addEventListener('voiceschanged', refresh);
    return () => speechSynthesis.removeEventListener('voiceschanged', refresh);
  });
  async function loadGoogle() {
    loading = true;
    error = '';
    try {
      const response = await fetch(`/api/speech/voices?language=${encodeURIComponent(speech.language)}`);
      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.error);
      }
      googleVoices = data;
      if (!data.some((item: { name: string }) => item.name === speech.voice)) {
        speech.voice = data[0]?.name ?? '';
      }
      speech.save();
    } catch (cause) {
      error = (cause as Error).message;
    } finally {
      loading = false;
    }
  }
</script>

<div class="speech-settings">
  <p>Read the same evidence-linked story shown in the narrative panel. Changes apply to the current event.</p>
  <Label
    >Speech engine
    <Select placeholder=""
      aria-label="Speech engine"
      bind:value={speech.provider}
      onchange={() => {
        speech.voice = '';
        speech.save();
      }}
    >
      <option value="browser">Browser voices</option><option value="google">Google Cloud Text-to-Speech</option>
      <option value="screen">Read on screen · no audio</option>
    </Select>
  </Label>
  <Label>Language <Input bind:value={speech.language} onchange={() => speech.save()} placeholder="en-US" /></Label>
  {#if speech.provider === 'browser'}
    <Label
      >Voice <Select placeholder="" aria-label="Voice" bind:value={speech.voice} onchange={() => speech.save()}>
        <option value="">System default</option>
        {#each browserVoices as voice}<option value={voice.voiceURI}>{voice.name} · {voice.lang}</option>{/each}
      </Select></Label
    >
    <p>Google voices appear here when supplied by your browser. Availability depends on your system.</p>
    {#if !browserVoices.length}<p>
        No browser voices are listed. Choose Google Cloud or read on screen if system speech is unavailable.
      </p>{/if}
  {:else if speech.provider === 'google'}
    <p>
      Narration text is sent to Google. The local backend uses Application Default Credentials; see <code
        >docs/speech.md</code
      >.
    </p>
    <Button variant="primary" disabled={loading} onclick={loadGoogle}
      >{loading ? 'Loading voices…' : 'Load Google voices'}</Button
    >
    <Label
      >Voice <Select placeholder="" aria-label="Voice" bind:value={speech.voice} onchange={() => speech.save()}>
        <option value="">Choose a Google voice</option>
        {#each googleVoices as voice}<option value={voice.name}>{voice.name} · {voice.ssmlGender}</option>{/each}
      </Select></Label
    >
  {/if}
  <Label
    >Reading speed <Select placeholder="" aria-label="Reading speed" bind:value={ws.speed}
      ><option value={0.75}>0.75×</option><option value={1}>1×</option>
      <option value={1.25}>1.25×</option><option value={1.5}>1.5×</option><option value={2}>2×</option></Select
    ></Label
  >
  {#if error || speech.error}<p role="alert" class="speech-error">{error || speech.error}</p>{/if}
  <p role="status">{speech.status}</p>
</div>
