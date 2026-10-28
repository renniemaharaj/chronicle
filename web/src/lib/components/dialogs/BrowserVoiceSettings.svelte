<script lang="ts">
  import { onMount } from 'svelte';
  import { Label, Select } from 'flowbite-svelte';
  import { speech } from '../../speech/preferences.svelte';

  let voices = $state<SpeechSynthesisVoice[]>([]);
  onMount(() => {
    if (!('speechSynthesis' in window)) {
      return;
    }

    const refresh = () => {
      voices = speechSynthesis.getVoices();
    };

    refresh();
    speechSynthesis.addEventListener('voiceschanged', refresh);

    return () => speechSynthesis.removeEventListener('voiceschanged', refresh);
  });
</script>

<Label>
  Voice
  <Select
    placeholder=""
    aria-label="Voice"
    bind:value={speech.voice}
    onchange={() => speech.save()}
  >
    <option value="">System default</option>
    {#each voices as voice}
      <option value={voice.voiceURI}>{voice.name} · {voice.lang}</option>
    {/each}
  </Select>
</Label>
<p>Browser voices depend on your system. They speak each time and cannot be saved as audio files here.</p>
{#if !voices.length}
  <p>No browser voices are listed. Choose a cloud voice or read on screen.</p>
{/if}
