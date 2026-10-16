<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Input, Label, Select } from 'flowbite-svelte';
  import { onMount } from 'svelte';
  import { Plug, Save } from '@lucide/svelte';
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

<p class="dialog-intro">Choose a model and connect your research sources. Settings are saved on this workstation.</p>
<form class="settings-form" onsubmit={(event) => { event.preventDefault(); save(); }}>
  <fieldset class="settings-section">
    <legend><span>1</span> Model connection</legend>
    <p>Select the model that will investigate and organize your research.</p>
    <div class="form-row">
      <Label>Provider
        <Select placeholder="" aria-label="Provider" bind:value={settings.provider}>
          <option value="ollama">Ollama</option>
          <option value="openai">OpenAI compatible</option>
        </Select>
      </Label>
      <Label>Model name<Input bind:value={settings.model} placeholder="e.g. qwen3:8b" required /></Label>
    </div>
    <Label>Endpoint URL
      <Input type="url" bind:value={settings.baseUrl} placeholder="http://localhost:11434" required />
    </Label>
    {#if settings.provider === 'openai'}
      <Label>API key
        <Input type="password" bind:value={settings.apiKey} autocomplete="off"
          placeholder={settings.hasApiKey ? 'Leave blank to keep the saved key' : 'Optional for local endpoints'} />
      </Label>
    {/if}
  </fieldset>
  <fieldset class="settings-section">
    <legend><span>2</span> Research sources</legend>
    <p>Connect an MCP server, then choose its search and page retrieval tools.</p>
    <Label>Server URL
      <Input type="url" bind:value={settings.mcpUrl} placeholder="http://localhost:3000/mcp" />
    </Label>
    <Button type="button" variant="secondary" class="settings-discover" disabled={busy || !settings.mcpUrl} onclick={discover}>
      <Plug size={15} /> {busy ? 'Connecting…' : 'Discover available tools'}
    </Button>
    <div class="form-row">
      <Label>Search tool<Input bind:value={settings.searchTool} list="mcp-tools" placeholder="Select or enter tool name" /></Label>
      <Label>Retrieval tool<Input bind:value={settings.retrieveTool} list="mcp-tools" placeholder="Select or enter tool name" /></Label>
    </div>
    <datalist id="mcp-tools">
      {#each tools as tool}<option value={tool.name}>{tool.description}</option>{/each}
    </datalist>
  </fieldset>
  <fieldset class="settings-section">
    <legend><span>3</span> Research limits</legend>
    <p>Control how much work each investigation can run.</p>
    <div class="form-row">
      <Label>Request budget<Input type="number" min="4" max="100" bind:value={settings.maxRequests} required /></Label>
      <Label>Concurrent questions<Input type="number" min="1" max="4" bind:value={settings.concurrency} required /></Label>
    </div>
  </fieldset>
  <footer class="settings-footer">
    {#if message}<p class="form-message" role="status">{message}</p>{/if}
    <div class="dialog-actions">
      <Button type="button" variant="secondary" onclick={() => { ws.dialog = ''; }}>Close</Button>
      <Button type="submit" variant="primary" disabled={busy}><Save size={15} />{busy ? 'Working…' : 'Save configuration'}</Button>
    </div>
  </footer>
</form>
