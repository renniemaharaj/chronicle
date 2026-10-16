<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
</script>

<details class="narrative-themes">
  <summary>Key themes</summary>
  {#each [...new Set(ws.project.events.map((event) => event.kind))] as kind}
    <Button
      onclick={() => {
        const event = ws.project.events.find((item) => item.kind === kind);
        if (event) ws.select(event.id);
      }}>{kind} · {ws.project.events.filter((event) => event.kind === kind).length}</Button
    >
  {/each}
</details>
