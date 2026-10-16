<script lang="ts">
  import { workspace as ws } from '../../state/workspace.svelte';
  import { safeSourceURL } from '../../api/links';
  let { claimIds }: { claimIds: string[] } = $props();
</script>

<div class="context-evidence">
  {#each claimIds as id}
    {@const claim = ws.project.claims.find((item) => item.id === id)}
    {#if claim}
      <details>
        <summary>{claim.text}</summary>
        <p>{claim.status.replaceAll('_', ' ')}</p>
        {#each ws.project.evidence.filter((item) => item.claimId === id) as evidence}
          {@const source = ws.project.sources.find((item) => item.id === evidence.sourceId)}
          <blockquote>{evidence.excerpt}</blockquote>
          <small>{evidence.relation} · {evidence.mode} · {evidence.location}</small>
          {#if safeSourceURL(source?.url)}
            <a href={safeSourceURL(source?.url)} target="_blank" rel="noopener noreferrer">{source?.title}</a>
          {:else}<p>{source?.title ?? 'Source unavailable'}</p>{/if}
        {/each}
      </details>
    {/if}
  {/each}
</div>
