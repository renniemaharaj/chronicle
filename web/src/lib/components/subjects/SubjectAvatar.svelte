<script lang="ts">
  import type { Entity } from '../../types/knowledge';
  import type { Project } from '../../types/project';
  import { portraitFor } from '../../utils/subjects';
  let { entity, project, size = 52 }: { entity: Entity; project: Project; size?: number } = $props();
  let failed = $state('');
  const portrait = $derived(portraitFor(entity, project));
  const initials = $derived(
    entity.name
      .split(/\s+/)
      .slice(0, 2)
      .map((part) => part[0])
      .join(''),
  );
</script>

<span
  class="subject-avatar"
  style={`--avatar-size: ${size}px`}
  title={portrait.url ? portrait.credit : `No verified image: ${entity.name}`}
>
  {#if portrait.url && portrait.url !== failed}
    <img
      src={portrait.url}
      alt={entity.name}
      loading="lazy"
      referrerpolicy="no-referrer"
      onerror={() => {
        failed = portrait.url ?? '';
      }}
    />
  {:else}<span aria-label={`${entity.name}, image unavailable`}>{initials}</span>{/if}
</span>
