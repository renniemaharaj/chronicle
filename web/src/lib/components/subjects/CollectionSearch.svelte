<script lang="ts">
  import { Input, Label } from 'flowbite-svelte';
  import { onMount, type Snippet } from 'svelte';
  import { Search } from '@lucide/svelte';

  let {
    value = $bindable(''), label, placeholder, summary = '', controls, oninput,
  }: {
    value?: string;
    label: string;
    placeholder: string;
    summary?: string;
    controls?: Snippet;
    oninput?: () => void;
  } = $props();
  let input = $state<HTMLInputElement>();
  onMount(() => input?.focus());
</script>

<div class="subject-search-toolbar">
  <Label class="subject-search">
    <Search size={18} />
    <Input bind:elementRef={input} aria-label={label} {placeholder} bind:value {oninput} />
  </Label>
  {#if summary || controls}
    <div class="subject-search-meta">
      <span>{summary}</span>
      {@render controls?.()}
    </div>
  {/if}
</div>
