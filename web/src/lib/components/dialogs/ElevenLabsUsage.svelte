<script lang="ts">
  import type { ElevenUsage } from '../../speech/elevenlabs.svelte';

  let { usage, error }: { usage: ElevenUsage | null; error: string } = $props();
</script>

{#if usage}
  <div role="status">
    <p>Plan: {usage.tier} · {usage.status}</p>
    <p>
      Account usage: {usage.character_count.toLocaleString()} / {usage.character_limit.toLocaleString()}
      characters/credits reported by ElevenLabs.
    </p>
    <p>
      Remaining: {Math.max(0, usage.character_limit - usage.character_count).toLocaleString()}.
      {#if usage.next_character_count_reset_unix}
        Resets {new Date(usage.next_character_count_reset_unix * 1000).toLocaleDateString()}.
      {/if}
    </p>
    <p>Usage covers your account. Credit cost varies by model; refresh after generating audio.</p>
  </div>
{:else if error}
  <p role="status">{error} Voice playback may still be available with this key's permissions.</p>
{/if}
