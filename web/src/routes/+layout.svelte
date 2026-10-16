<script lang="ts">
  import { ThemeProvider } from 'flowbite-svelte';
  import { chronicleTheme } from '$lib/components/ui/theme';
  import '$lib/styles/flowbite.css';
  import '@xyflow/svelte/dist/style.css';
  import './app.css';
  import type { Snippet } from 'svelte';
  let { children }: { children: Snippet } = $props();
</script>

<ThemeProvider theme={chronicleTheme}>
  {@render children()}
</ThemeProvider>
