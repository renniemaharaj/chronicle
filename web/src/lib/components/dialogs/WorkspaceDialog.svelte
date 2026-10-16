<script lang="ts">
  import { BookOpen, Download, FolderOpen, Plus, Upload } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { importProject, saveProject } from '../../state/persistence';
  import { exportFormats, exportProject } from '../../api/exports';

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
    {#each exportFormats as format}<button
        onclick={() => {
          try {
            exportProject(ws.project, format);
          } catch (error) {
            ws.notify(String(error));
          }
        }}
        ><Download size={17} /><span>{format}</span><small
          >{format === 'Project JSON' ? 'Portable format v1' : 'Preserves source references'}</small
        ></button
      >{/each}
  </div>
{:else}
  <p class="dialog-intro">Your saved research projects. Import a portable project or begin a new investigation.</p>
  {#if !ws.projects.length}<p class="dialog-intro">No research projects yet.</p>{/if}
  <div class="library-projects">
    {#each ws.projects as project}<button class="library-project" onclick={() => ws.open(project)}
        ><FolderOpen size={23} /><span
          ><strong>{project.title}</strong><small
            >{project.events.length} events · {project.claims.length} claims · {project.conflicts.length} conflicts</small
          ><small>{project.fixture ? 'Fixture · incomplete coverage' : 'Research project'}</small></span
        ><BookOpen size={15} /></button
      >{/each}
  </div>
  <div class="dialog-actions">
    <label class="button secondary"
      ><Upload size={14} /> Import project<input
        type="file"
        accept="application/json,.json"
        onchange={upload}
        hidden
      /></label
    ><button class="button" onclick={() => (ws.dialog = 'Research')}><Plus size={14} /> New research</button>
  </div>
{/if}
