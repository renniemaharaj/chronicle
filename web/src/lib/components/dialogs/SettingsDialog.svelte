<script lang="ts">
  import { onMount } from 'svelte';
  import { Check, Plug, Save } from '@lucide/svelte';
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
    } catch (error) {
      message = String(error);
    } finally {
      busy = false;
    }
  }

  async function discover() {
    await save();
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
  Connect a local model or an OpenAI-compatible endpoint and an MCP research server. Credentials stay in your local
  configuration.
</p>
<form
  class="settings-form"
  onsubmit={(event) => {
    event.preventDefault();
    save();
  }}
>
  <div class="form-row">
    <label
      >Model provider<select bind:value={settings.provider}
        ><option value="ollama">Ollama</option><option value="openai">OpenAI compatible</option></select
      ></label
    ><label>Model name<input bind:value={settings.model} placeholder="e.g. qwen3:8b" required /></label>
  </div>
  <label
    >Model endpoint<input
      type="url"
      bind:value={settings.baseUrl}
      placeholder="http://localhost:11434"
      required
    /></label
  >
  {#if settings.provider === 'openai'}<label
      >API key<input
        type="password"
        bind:value={settings.apiKey}
        placeholder={settings.hasApiKey ? 'Saved key (leave blank to retain)' : 'Optional for local endpoints'}
        autocomplete="off"
      /></label
    >{/if}
  <label>MCP server URL<input type="url" bind:value={settings.mcpUrl} placeholder="http://localhost:3000/mcp" /></label>
  <button type="button" class="button secondary" disabled={busy || !settings.mcpUrl} onclick={discover}
    ><Plug size={14} /> Discover tools</button
  >
  <div class="form-row">
    <label>Search tool<input bind:value={settings.searchTool} list="mcp-tools" placeholder="Tool name" /></label><label
      >Retrieval tool<input bind:value={settings.retrieveTool} list="mcp-tools" placeholder="Tool name" /></label
    >
  </div>
  <datalist id="mcp-tools"
    >{#each tools as tool}<option value={tool.name}>{tool.description}</option>{/each}</datalist
  >
  <div class="form-row">
    <label>Request budget<input type="number" min="4" max="100" bind:value={settings.maxRequests} /></label><label
      >Concurrent questions<input type="number" min="1" max="4" bind:value={settings.concurrency} /></label
    >
  </div>
  <p class="settings-note">
    MCP tool schemas are inspected before use. Research runs record model settings and tool activity. Retrieved
    documents are treated as untrusted evidence.
  </p>
  {#if message}<p class="form-message" role="status">{message}</p>{/if}
  <div class="dialog-actions">
    <button class="button" disabled={busy}><Save size={14} /> Save configuration</button>
  </div>
</form>
