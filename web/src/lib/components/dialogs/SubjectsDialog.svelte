<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Select } from 'flowbite-svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { contemporary, entityDescription, searchText } from '../../utils/subjects';
  import SubjectCard from '../subjects/SubjectCard.svelte';
  import CollectionSearch from '../subjects/CollectionSearch.svelte';
  let query = $state('');
  let kind = $state('All kinds');
  let limit = $state(30);
  const contemporaries = $derived(ws.dialog === 'Contemporaries');
  const candidates = $derived(
    ws.project.entities.filter(
      (entity) => entity.id !== ws.subject.id && (!contemporaries || contemporary(entity, ws.subject) !== 'outside'),
    ),
  );
  const kinds = $derived([...new Set(candidates.map((entity) => entity.kind))].sort());
  const results = $derived(
    candidates.filter(
      (entity) =>
        (kind === 'All kinds' || entity.kind === kind) &&
        searchText(`${entity.name} ${entity.summary} ${entityDescription(entity, ws.project)}`).includes(
          searchText(query),
        ),
    ),
  );
</script>

<CollectionSearch
  label="Search related subjects"
  placeholder="Search names, places, relationships…"
  bind:value={query}
  summary={`${results.length} of ${candidates.length} subjects`}
  oninput={() => { limit = 30; }}
>
  {#snippet controls()}
    <Select placeholder="" aria-label="Subject kind" bind:value={kind}
      ><option>All kinds</option>{#each kinds as item}<option>{item}</option>{/each}</Select
    >
  {/snippet}
</CollectionSearch>
{#if contemporaries}<p class="subject-search-note">
    People whose documented dates overlap {ws.subject.name}. Incomplete dates are labeled separately.
  </p>{/if}
<div class="subject-results">
  {#each results.slice(0, limit) as entity (entity.id)}<SubjectCard {entity} showOverlap={contemporaries} />{/each}
</div>
{#if !results.length}<p class="empty-state">No subjects match your search.</p>{/if}
{#if results.length > limit}<Button
    variant="secondary"
    onclick={() => {
      limit += 30;
    }}>Show more subjects</Button
  >{/if}
