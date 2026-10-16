<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { safeSourceURL } from "../../api/links";
  import { ArrowRight, Check, ShieldCheck } from '@lucide/svelte';
  import type { Claim, Evidence } from '../../types/knowledge';
  import { workspace as ws } from '../../state/workspace.svelte';
  let { claims, evidence }: { claims: Claim[]; evidence: Evidence[] } = $props();
</script>

<section class="evidence-section">
  <h3><ShieldCheck size={13} /> Evidence & claims <span>{claims.length}</span></h3>
  {#each claims as claim}
    <div class="claim-card">
      <div>
        <span class="claim-mark" class:supported={claim.status !== 'unverified'}><ShieldCheck size={11} /></span><strong
          >{claim.text}</strong
        >
      </div>
      <span class="claim-status" class:disputed={claim.status === 'disputed'}>{claim.status.replaceAll('_', ' ')}</span>
      {#each evidence.filter((item) => item.claimId === claim.id) as item}
        {@const source = ws.project.sources.find((source) => source.id === item.sourceId)}
        <p class="evidence-excerpt">{item.mode === 'quotation' ? `“${item.excerpt}”` : item.excerpt}</p>
        <small>{item.mode} · {item.relation} · {item.location}</small>
        <a href={safeSourceURL(source?.url)} target="_blank" rel="noreferrer"
          >{source?.title ?? 'Source missing'} <ArrowRight size={10} /></a
        >
        {#if !item.verified}<small class="evidence-limit"
            >{ws.project.fixture
              ? 'Curated source summary; no machine-verified quotation.'
              : 'Passage verification incomplete.'}</small
          >{/if}
      {/each}
      {#each claim.assessment as reason}<small class="assessment-reason">{reason}</small>{/each}
      {#if !claim.evidenceIds.length}<p class="source-note">No supporting passage has been retrieved.</p>{/if}
    </div>
  {/each}
  {#if !claims.length}<p class="empty-state">No factual claims are attached to this record.</p>{/if}
  <Button class="show-sources" onclick={() => (ws.inspectorTab = 'Sources')}
    >Inspect source provenance <ArrowRight size={12} /></Button
  >
</section>
