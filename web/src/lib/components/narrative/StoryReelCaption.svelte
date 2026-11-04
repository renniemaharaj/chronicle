<script lang="ts">
  import { workspace as ws } from '../../state/workspace.svelte';
  import { presentationTime } from '../../state/story-presentation.svelte';
  import { activeStoryCaption } from '../../utils/story-captions';
  import { storyPassages } from '../../speech/passages';

  const passage = $derived(
    storyPassages(ws.project, ws.orderedEvents, ws.activeClaims).find((item) => item.id === ws.event?.id),
  );
  const time = $derived(presentationTime(ws.project.id, ws.event?.id));
  const caption = $derived(activeStoryCaption(passage?.text ?? '', time.seconds, time.duration));
</script>

{#if caption}
  <div
    class="story-reel-caption"
    aria-label="Narration captions"
  >
    {caption}
  </div>
{/if}
