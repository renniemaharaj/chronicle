<script lang="ts">
  import type { Campaign, Part } from '../../campaign/types';
  import { dateLabel } from '../../types/temporal';

  let { campaign, part, index }: { campaign: Campaign; part: Part; index: number } = $props();
  const passage = $derived(part.passages[index]);
</script>

<article
  class="campaign-stage"
  aria-label="Vertical reel capture region"
>
  <header>
    <span>CHRONICLE</span>
    <span>Part {part.number} / {campaign.parts.length}</span>
  </header>
  <div class="campaign-stage-title">
    <span class="eyebrow">{passage?.contextId ? 'Historical context' : 'The story'}</span>
    <h2>{campaign.title}</h2>
  </div>
  {#if passage}
    <div class="campaign-stage-passage">
      <p class="campaign-stage-date">
        {dateLabel(passage.time.start)}{passage.time.end ? ` to ${dateLabel(passage.time.end)}` : ''}
      </p>
      <small>
        {passage.time.start?.calendar ?? 'Calendar not recorded'}{passage.time.openEnd
          ? ' · end unknown'
          : ''}
      </small>
      {#if passage.time.start?.uncertaintyYears}
        <small>Date uncertainty: ±{passage.time.start.uncertaintyYears} years</small>
      {/if}
      <h3>{passage.title}</h3>
      <p class="campaign-caption">{passage.text}</p>
    </div>
    <footer>
      <span>Evidence and competing accounts preserved in Chronicle</span>
      <span>{index + 1} / {part.passages.length} passages · {passage.claimIds.length} linked claims</span>
    </footer>
  {/if}
</article>
