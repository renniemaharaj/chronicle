<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Input, Label } from 'flowbite-svelte';
  import { onMount } from 'svelte';
  import { Save } from '@lucide/svelte';
  import ModelSettings from './ModelSettings.svelte';
  import ResearchSources from './ResearchSources.svelte';
  import { api } from '../../api/client';
  import type { ProviderSettings } from '../../types/project';
  import { workspace as ws } from '../../state/workspace.svelte';

  let settings = $state<ProviderSettings>({
    provider: 'ollama',
    baseUrl: 'http://localhost:11434',
    model: '',
    mcpUrl: '',
    searchTool: '',
    retrieveTool: '',
    maxRequests: 40,
    concurrency: 2,
  });
  let tools = $state<{ name: string; description: string }[]>([]);
  let message = $state('');
  let busy = $state(false);

  onMount(async () => {
    try {
      settings = await api.settings();
    } catch {
      message = 'Start the local Go API to configure research providers.';
    }
  });

  async function save() {
    busy = true;
    try {
      settings = await api.saveSettings(settings);
      message = 'Settings saved on this workstation.';

      return true;
    } catch (error) {
      message = String(error);

      return false;
    } finally {
      busy = false;
    }
  }

  async function discover() {
    if (!(await save())) {
      return;
    }

    busy = true;
    try {
      tools = await api.tools();
      message = `${tools.length} MCP tools discovered. Select search and retrieval tools below.`;
    } catch (error) {
      message = String(error);
    } finally {
      busy = false;
    }
  }
</script>

<p class="dialog-intro">
  Choose a model and connect your research sources. Settings are saved on this workstation.
</p>
<form
  class="settings-form"
  onsubmit={(event) => {
    event.preventDefault();
    save();
  }}
>
  <ModelSettings bind:settings />
  <ResearchSources
    bind:settings
    {tools}
    {busy}
    {discover}
  />
  <fieldset class="settings-section">
    <legend>
      <span>3</span>
      Research limits
    </legend>
    <p>Control how much work each investigation can run.</p>
    <div class="form-row">
      <Label>
        Request budget<Input
          type="number"
          min="4"
          max="100"
          bind:value={settings.maxRequests}
          required
        />
      </Label>
      <Label>
        Concurrent questions<Input
          type="number"
          min="1"
          max="4"
          bind:value={settings.concurrency}
          required
        />
      </Label>
    </div>
  </fieldset>
  <footer class="settings-footer">
    {#if message}
      <p
        class="form-message"
        role="status"
      >
        {message}
      </p>
    {/if}
    <div class="dialog-actions">
      <Button
        type="button"
        variant="secondary"
        onclick={() => {
          ws.dialog = '';
        }}
      >
        Close
      </Button>
      <Button
        type="submit"
        variant="primary"
        disabled={busy}
      >
        <Save size={15} />{busy ? 'Working…' : 'Save configuration'}
      </Button>
    </div>
  </footer>
</form>
