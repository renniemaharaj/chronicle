<script lang="ts">
  import Button from '../ui/Button.svelte';
  import { Label, Select } from 'flowbite-svelte';
  import { speech } from '../../speech/preferences.svelte';
  import { speechRequest } from '../../speech/audio-api';

  let voices = $state<{ name: string; ssmlGender: string }[]>([]);
  let loading = $state(false);
  let error = $state('');

  async function load() {
    loading = true;
    error = '';
    try {
      voices = await speechRequest(`voices?language=${encodeURIComponent(speech.language)}`);
    } catch (cause) {
      error = (cause as Error).message;
    } finally {
      loading = false;
    }
  }
</script>

<p>
  Google receives narration during generation. The backend uses Application Default Credentials; see
  docs/speech.md.
</p>
<Button
  disabled={loading}
  onclick={load}
>
  {loading ? 'Loading voices…' : 'Load Google voices'}
</Button>
<Label>
  Voice
  <Select
    placeholder=""
    aria-label="Voice"
    bind:value={speech.voice}
    onchange={() => speech.save()}
  >
    <option value="">Choose a Google voice</option>
    {#if speech.voice && !voices.some((voice) => voice.name === speech.voice)}
      <option value={speech.voice}>{speech.voice}</option>
    {/if}
    {#each voices as voice}
      <option value={voice.name}>{voice.name} · {voice.ssmlGender}</option>
    {/each}
  </Select>
</Label>
{#if error}
  <p
    role="alert"
    class="speech-error"
  >
    {error}
  </p>
{/if}
