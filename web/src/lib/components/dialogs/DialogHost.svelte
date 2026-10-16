<script lang="ts">
  import { Modal } from 'flowbite-svelte';
  import Button from '../ui/Button.svelte';
  import { X } from '@lucide/svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import WorkspaceDialog from './WorkspaceDialog.svelte';
  import SettingsDialog from './SettingsDialog.svelte';
  import ResearchDialog from './ResearchDialog.svelte';
  import HistoricalContext from '../narrative/HistoricalContext.svelte';
  import SpeechDialog from './SpeechDialog.svelte';
  import SubjectsDialog from './SubjectsDialog.svelte';
  import MediaDialog from './MediaDialog.svelte';
  import SubjectLinksDialog from './SubjectLinksDialog.svelte';
  import RecordDialog from './RecordDialog.svelte';
</script>

{#if ws.dialog}
  <Modal
    bind:open={() => Boolean(ws.dialog), (open) => { if (!open) ws.dialog = ''; }}
    dismissable={false}
    outsideclose
    size={ws.dialog === 'Settings' ? 'lg' : 'md'}
    aria-label={ws.dialog}
    class="chronicle-modal"
    classes={{ header: 'dialog-header', body: 'dialog-body' }}
  >
      {#snippet header()}
        <div>
          <span class="eyebrow">Chronicle workspace</span>
          <h2>{ws.dialog}</h2>
        </div>
        <Button icon aria-label="Close dialog" onclick={() => (ws.dialog = '')}><X size={18} /></Button>
      {/snippet}
        {#if ['Library', 'Exports'].includes(ws.dialog)}<WorkspaceDialog />
        {:else if ['Related Subjects', 'Contemporaries'].includes(ws.dialog)}<SubjectsDialog />
        {:else if ['Images', 'Videos'].includes(ws.dialog)}<MediaDialog />
        {:else if ws.dialog === 'Social Networks & Website'}<SubjectLinksDialog />
        {:else if ws.dialog === 'Historical context'}<HistoricalContext />
        {:else if ws.dialog === 'Voice & playback'}<SpeechDialog />
        {:else if ws.dialog === 'Settings'}<SettingsDialog />
        {:else if ws.dialog === 'Research'}<ResearchDialog />
        {:else}<RecordDialog />{/if}
  </Modal>
{/if}
