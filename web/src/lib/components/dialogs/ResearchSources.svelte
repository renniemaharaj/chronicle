<script lang="ts">
  import Button from '../ui/Button.svelte';
  import { Input, Label } from 'flowbite-svelte';
  import { Plug } from '@lucide/svelte';
  import type { ProviderSettings } from '../../types/project';

  let {
    settings = $bindable(),
    tools,
    busy,
    discover,
  }: {
    settings: ProviderSettings;
    tools: { name: string; description: string }[];
    busy: boolean;
    discover: () => Promise<void>;
  } = $props();
</script>

<fieldset class="settings-section">
  <legend>
    <span>2</span>
    Research sources
  </legend>
  <p>Connect an MCP server, then choose its search and page retrieval tools.</p>
  <Label>
    Server URL
    <Input
      type="url"
      bind:value={settings.mcpUrl}
      placeholder="http://localhost:3000/mcp"
    />
  </Label>
  <Button
    type="button"
    variant="secondary"
    class="settings-discover"
    disabled={busy || !settings.mcpUrl}
    onclick={discover}
  >
    <Plug size={15} />
    {busy ? 'Connecting…' : 'Discover available tools'}
  </Button>
  <div class="form-row">
    <Label>
      Search tool<Input
        bind:value={settings.searchTool}
        list="mcp-tools"
        placeholder="Select or enter tool name"
      />
    </Label>
    <Label>
      Retrieval tool<Input
        bind:value={settings.retrieveTool}
        list="mcp-tools"
        placeholder="Select or enter tool name"
      />
    </Label>
  </div>
  <datalist id="mcp-tools">
    {#each tools as tool}
      <option value={tool.name}>{tool.description}</option>
    {/each}
  </datalist>
</fieldset>
