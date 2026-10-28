<script lang="ts">
  import type { Project } from '../../types/project';
  import { storyCompletion } from '../../utils/completion';

  let { project }: { project: Project } = $props();
  const status = $derived(storyCompletion(project));
</script>

<div class="story-readiness">
  <p>Evidence linked: {status.evidencePercent}% of saved events ({status.evidence}/{status.total}).</p>
  <p>Narration text: {status.narrationPercent}% complete ({status.narration}/{status.total} passages).</p>
  <p>Audio readiness is checked in Voice &amp; playback for your selected voice and model.</p>
  {#if status.reviewedAt}
    <p>Last research review: {new Date(status.reviewedAt).toLocaleDateString()}.</p>
  {/if}
  <p>These percentages cover saved events and tasks. Historical coverage can still have gaps.</p>
</div>
