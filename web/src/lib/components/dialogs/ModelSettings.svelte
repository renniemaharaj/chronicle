<script lang="ts">
  import { Input, Label, Select } from 'flowbite-svelte';
  import type { ProviderSettings } from '../../types/project';

  let { settings = $bindable() }: { settings: ProviderSettings } = $props();
</script>

<fieldset class="settings-section">
  <legend>
    <span>1</span>
    Model connection
  </legend>
  <p>Select the model that will investigate and organize your research.</p>
  <div class="form-row">
    <Label>
      Provider
      <Select
        placeholder=""
        aria-label="Provider"
        bind:value={settings.provider}
      >
        <option value="ollama">Ollama</option>
        <option value="openai">OpenAI compatible</option>
      </Select>
    </Label>
    <Label>
      Model name<Input
        bind:value={settings.model}
        placeholder="e.g. qwen3:8b"
        required
      />
    </Label>
  </div>
  <Label>
    Endpoint URL
    <Input
      type="url"
      bind:value={settings.baseUrl}
      placeholder="http://localhost:11434"
      required
    />
  </Label>
  {#if settings.provider === 'openai'}
    <Label>
      API key
      <Input
        type="password"
        bind:value={settings.apiKey}
        autocomplete="off"
        placeholder={settings.hasApiKey
          ? 'Leave blank to keep the saved key'
          : 'Optional for local endpoints'}
      />
    </Label>
  {/if}
</fieldset>
