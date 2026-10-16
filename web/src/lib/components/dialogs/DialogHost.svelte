<script lang="ts">
  import { X } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import WorkspaceDialog from './WorkspaceDialog.svelte';
  import SettingsDialog from './SettingsDialog.svelte';
  import ResearchDialog from './ResearchDialog.svelte';
  import RecordDialog from './RecordDialog.svelte';
</script>

{#if ws.dialog}
  <div
    class="modal-backdrop"
    role="presentation"
    onclick={(event) => {
      if (event.target === event.currentTarget) ws.dialog = '';
    }}
  >
    <div class="dialog" role="dialog" aria-modal="true" aria-label={ws.dialog} tabindex="-1">
      <div class="dialog-header">
        <div>
          <span class="eyebrow">Chronicle workspace</span>
          <h2>{ws.dialog}</h2>
        </div>
        <button aria-label="Close dialog" onclick={() => (ws.dialog = '')}><X size={18} /></button>
      </div>
      <div class="dialog-body">
        {#if ['Library', 'Exports'].includes(ws.dialog)}<WorkspaceDialog />
        {:else if ws.dialog === 'Settings'}<SettingsDialog />
        {:else if ws.dialog === 'Research'}<ResearchDialog />
        {:else}<RecordDialog />{/if}
      </div>
    </div>
  </div>
{/if}
