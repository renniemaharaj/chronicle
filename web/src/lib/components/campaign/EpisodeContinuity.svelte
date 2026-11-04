<script lang="ts">
  import { Route } from '@lucide/svelte';
  import { dateLabel } from '../../types/temporal';
  import type { Campaign, Part } from '../../campaign/types';

  let { campaign, part }: { campaign: Campaign; part: Part } = $props();
  const first = $derived(part.passages[0]);
  const last = $derived(part.passages.at(-1));
</script>

<details class="reel-detail-card">
  <summary><Route size={16} />Story continuity</summary>
  <p class="reel-detail-note">
    {part.passages.length} reading passages. {part.number > 1
      ? `Continues from episode ${part.number - 1}.`
      : 'Begins the series.'}
  </p>
  <div class="reel-continuity-stop">
    <small>Opens with · {dateLabel(first?.time.start)}</small>
    <strong>{first?.title}</strong>
    <p>{first?.text}</p>
  </div>
  <div class="reel-continuity-stop">
    <small>Closes with · {dateLabel(last?.time.start)}</small>
    <strong>{last?.title}</strong>
    <p>{last?.text}</p>
  </div>
  <p class="reel-detail-note">
    {part.number < campaign.parts.length
      ? `Next: episode ${part.number + 1}.`
      : 'Final episode in this saved series.'} Recorded dates and competing accounts stay with the story.
  </p>
</details>
