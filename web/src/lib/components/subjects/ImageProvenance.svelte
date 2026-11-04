<script lang="ts">
  import type { MediaAsset } from '../../types/media';
  import type { Project } from '../../types/project';
  import { intervalLabel } from '../../types/temporal';
  import { safeSourceURL } from '../../api/links';

  let { asset, project }: { asset: MediaAsset; project: Project } = $props();
</script>

{#if asset.generated}
  <small class="story-visual-kind">AI illustration · Interpretive reconstruction</small>
  {#if asset.depicts}
    <small>Depicts: {intervalLabel(asset.depicts)}</small>
  {/if}
  <details class="image-generation-details">
    <summary>Generation & historical guidance</summary>
    <p>Generated with {asset.generated.tool}. This illustration is separate from historical evidence.</p>
    <p>{asset.generated.prompt}</p>
    {#each asset.generated.referenceSourceIds as id}
      {@const source = project.sources.find((item) => item.id === id)}
      {#if source && safeSourceURL(source.url)}
        <a
          href={safeSourceURL(source.url)}
          target="_blank"
          rel="noopener noreferrer"
        >
          {source.title}
        </a>
      {/if}
    {/each}
  </details>
{/if}
