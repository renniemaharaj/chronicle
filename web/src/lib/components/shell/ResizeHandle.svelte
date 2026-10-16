<script lang="ts">
  import { workspace as ws } from '../../state/workspace.svelte';
  let { side }: { side: 'left' | 'right' } = $props();

  function resize(event: PointerEvent) {
    const start = event.clientX;
    const width = side === 'left' ? ws.leftWidth : ws.rightWidth;
    const move = (pointer: PointerEvent) => {
      const delta = (pointer.clientX - start) * (side === 'left' ? 1 : -1);
      if (side === 'left') ws.leftWidth = Math.max(180, Math.min(350, width + delta));
      else ws.rightWidth = Math.max(280, Math.min(480, width + delta));
    };
    const stop = () => {
      window.removeEventListener('pointermove', move);
      window.removeEventListener('pointerup', stop);
    };
    window.addEventListener('pointermove', move);
    window.addEventListener('pointerup', stop);
  }
</script>

<button
  class="resize-handle {side}"
  aria-label={`Resize ${side} panel`}
  onpointerdown={resize}
  onkeydown={(event) => {
    const delta = event.key === 'ArrowRight' ? 10 : event.key === 'ArrowLeft' ? -10 : 0;
    if (side === 'left') ws.leftWidth = Math.max(180, Math.min(350, ws.leftWidth + delta));
    else ws.rightWidth = Math.max(280, Math.min(480, ws.rightWidth - delta));
  }}
></button>
