<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { BookOpen, Download, FolderOpen, Plus, Upload } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import SubjectAvatar from '../subjects/SubjectAvatar.svelte';
  import { searchText } from '../../utils/subjects';
  import { importProject, saveProject } from '../../state/persistence';
  import { exportFormats, exportProject } from '../../api/exports';
  import CollectionSearch from '../subjects/CollectionSearch.svelte';

  let query = $state('');
  const results = $derived(ws.projects.filter((project) => searchText(project.title).includes(searchText(query))));

  async function upload(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;
    try {
      const project = importProject(JSON.parse(await file.text()));
      ws.open(project);
      await saveProject(project);
    } catch (error) {
      ws.notify(String(error));
    }
  }
</script>

{#if ws.dialog === 'Exports'}
  <p class="dialog-intro">
    Export this project with its evidence, temporal uncertainty, and competing claims. PDF opens a print-ready report.
  </p>
  <div class="export-grid">
    {#each exportFormats as format}<Button
        onclick={() => {
          try {
            exportProject(ws.project, format);
          } catch (error) {
            ws.notify(String(error));
          }
        }}
        ><Download size={17} /><span>{format}</span><small
          >{format === 'Project JSON' ? 'Portable format v1' : 'Preserves source references'}</small
        ></Button
      >{/each}
  </div>
{:else}
  <p class="dialog-intro">Your saved research projects. Import a portable project or begin a new investigation.</p>
  {#if !ws.projects.length}<p class="dialog-intro">No research projects yet.</p>{/if}
  <CollectionSearch
    label="Search workspace projects"
    placeholder="Search your research projects…"
    bind:value={query}
  />
  <div class="library-projects">
    {#each results as project}
      {@const subject = project.entities.find((entity) => entity.id === project.subjectId)}<Button
        class="library-project"
        onclick={() => ws.open(project)}
        >{#if subject}<SubjectAvatar entity={subject} {project} size={72} />{:else}<FolderOpen size={23} />{/if}<span class="library-project-copy"
          ><strong>{project.title}</strong><small
            >{project.events.length} events · {project.claims.length} claims · {project.conflicts.length} conflicts</small
          ><small>{project.fixture ? 'Fixture · incomplete coverage' : 'Research project'}</small></span
        ><BookOpen size={15} /></Button
      >{/each}
  </div>
  <div class="dialog-actions">
    <Button tag="label" variant="secondary"
      ><Upload size={14} /> Import project<input
        type="file"
        accept="application/json,.json"
        onchange={upload}
        hidden
      /></Button
    ><Button variant="primary" onclick={() => (ws.dialog = 'Research')}><Plus size={14} /> New research</Button>
  </div>
{/if}
