<script lang="ts">
  import { Check, GitBranch } from '@lucide/svelte';
  import { orderingImplications } from '../../graph/interpretations';
  import { workspace as ws } from '../../state/workspace.svelte';
</script>

<div class="conflict-list">
  {#each ws.project.conflicts as conflict}
    <article class="conflict-card">
      <div class="conflict-card-title">
        <GitBranch size={14} /><strong>{conflict.title}</strong><span>{conflict.severity}</span>
      </div>
      <p>{conflict.explanation}</p>
      <div class="conflict-mode-buttons">
        {#each ['Shadows', 'Parallel', 'Isolation', 'Consequences'] as mode}<button
            class:active={ws.conflictMode === mode}
            onclick={() => {
              ws.conflictMode = mode;
              ws.showConflicts = true;
              ws.fitToken++;
            }}>{mode}</button
          >{/each}
      </div>
      {#each conflict.claimIds as claimId, index}
        {@const claim = ws.project.claims.find((claim) => claim.id === claimId)}
        {#if claim}
          <div class="hypothesis" class:alt={index % 2 === 1} class:adopted={ws.activeClaims[conflict.id] === claimId}>
            <small>Hypothesis {index + 1} · {claim.status.replaceAll('_', ' ')}</small>
            <strong>{claim.text}</strong>
            <span>{claim.evidenceIds.length} evidence record{claim.evidenceIds.length === 1 ? '' : 's'}</span>
            <button
              onclick={() => {
                ws.activeClaims[conflict.id] = claimId;
                ws.fitToken++;
              }}
            >
              {#if ws.activeClaims[conflict.id] === claimId}<Check size={10} /> Active interpretation{:else}Explore this
                interpretation{/if}
            </button>
            {#each ws.project.evidence.filter((item) => item.claimId === claimId) as evidence}
              <p class="comparison-excerpt">
                {evidence.excerpt}<small>{evidence.relation} · {evidence.location}</small>
              </p>
            {/each}
          </div>
        {/if}
      {/each}
      {#if ws.conflictMode === 'Consequences' || Object.hasOwn(ws.activeClaims, conflict.id)}
        <div class="consequences">
          <strong>Narrative implications</strong
          >{#each [...orderingImplications(ws.project, ws.activeClaims[conflict.id]), ...conflict.consequences] as consequence}<p
            >
              {consequence}
            </p>{/each}<small>Interpretation changes this view. Canonical claims and evidence are preserved.</small>
        </div>
      {/if}
    </article>
  {/each}
  {#if !ws.project.conflicts.length}<p class="empty-state">
      No disagreements have been recorded. This does not establish that none exist.
    </p>{/if}
</div>
