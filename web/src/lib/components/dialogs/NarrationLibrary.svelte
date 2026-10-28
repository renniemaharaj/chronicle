<script lang="ts">
  import { workspace as ws } from '../../state/workspace.svelte';
  import { speech } from '../../speech/preferences.svelte';
  import { storyPassages } from '../../speech/passages';
  import { speechChunks } from '../../speech/chunks';
  import type { CloudProvider } from '../../speech/audio-api';
  import NarrationPreparation from './NarrationPreparation.svelte';

  let { provider }: { provider: CloudProvider } = $props();
  const parts = $derived(
    storyPassages(ws.project, ws.orderedEvents, ws.activeClaims).flatMap((passage) =>
      speechChunks(passage.text).map((text) => ({
        text,
        voice: speech.voice,
        language: speech.language,
        model: speech.model,
        passageId: passage.id,
        title: ws.project.events.find((event) => event.id === passage.id)?.title,
      })),
    ),
  );
  const signature = $derived(JSON.stringify([provider, parts]));
</script>

{#if speech.voice && parts.length}
  {#key signature}
    <NarrationPreparation
      {provider}
      {parts}
    />
  {/key}
{:else}
  <p>Choose a voice and open a story to prepare saved narration.</p>
{/if}
