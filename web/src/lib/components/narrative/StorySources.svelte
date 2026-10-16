<script lang="ts">
  import { ExternalLink, FileText } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { safeSourceURL } from '../../api/links';
  import ContextEvidence from './ContextEvidence.svelte';

  const linkedClaims = $derived(new Set([
    ...(ws.event?.claimIds ?? []), ...(ws.event?.narration?.claimIds ?? []),
  ]));
  const claims = $derived(ws.project.claims.filter((claim) =>
    claim.eventId === ws.selectedId || linkedClaims.has(claim.id),
  ));
  const sourceIds = $derived(new Set(ws.project.evidence.filter((item) =>
    claims.some((claim) => claim.id === item.claimId),
  ).map((item) => item.sourceId)));
  const sources = $derived(ws.project.sources.filter((source) => sourceIds.has(source.id)));
</script>

<div class="story-sidebar-section-heading">
  <h2>Passage sources</h2><span>{sources.length}</span>
</div>
<div class="story-sidebar-records">
  {#each sources as source}
    <article>
      <FileText size={17} />
      <div>
        {#if safeSourceURL(source.url)}
          <a href={safeSourceURL(source.url)} target="_blank" rel="noreferrer">
            {source.title}<ExternalLink size={12} />
          </a>
        {:else}<strong>{source.title}</strong>{/if}
        <small>{source.publisher || source.author || 'Attribution not recorded'}</small>
        <small>{source.classification} source</small>
      </div>
    </article>
  {:else}
    <p class="story-empty-note">No source passages are linked to this event yet.</p>
  {/each}
  <ContextEvidence claimIds={claims.map((claim) => claim.id)} />
</div>
