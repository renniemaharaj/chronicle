<script lang="ts">
  import { onDestroy } from 'svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  let { side }: { side: 'left' | 'right' | 'story' } = $props();
  let stopResize = () => {};

  function adjust(delta: number, height: number) {
    if (side === 'left') {
      ws.leftWidth = Math.max(180, Math.min(350, ws.leftWidth + delta));
    } else if (side === 'right') {
      ws.rightWidth = Math.max(280, Math.min(480, ws.rightWidth - delta));
    } else {
      const maximum = Math.max(170, height - 270);
      const pixels = Math.max(170, Math.min(maximum, height * ws.storyHeight / 100 - delta));
      ws.storyHeight = pixels / height * 100;
    }
  }

  function resize(event: PointerEvent) {
    if (event.button !== 0) {
      return;
    }

    event.preventDefault();
    stopResize();
    const height = (event.currentTarget as HTMLElement).parentElement?.clientHeight ?? 700;
    let previous = side === 'story' ? event.clientY : event.clientX;
    const move = (pointer: PointerEvent) => {
      const current = side === 'story' ? pointer.clientY : pointer.clientX;
      adjust(current - previous, height);
      previous = current;
    };

    document.body.classList.add('resizing-panels');
    stopResize = () => {
      window.removeEventListener('pointermove', move);
      window.removeEventListener('pointerup', stopResize);
      window.removeEventListener('pointercancel', stopResize);
      document.body.classList.remove('resizing-panels');
    };
    window.addEventListener('pointermove', move);
    window.addEventListener('pointerup', stopResize);
    window.addEventListener('pointercancel', stopResize);
  }

  onDestroy(() => stopResize());
</script>

<button
  class="resize-handle {side}"
  aria-label={`Resize ${side} panel`}
  title={side === 'story' ? 'Drag to resize Story. Arrow keys also resize; double-click to reset.' : `Resize ${side} panel`}
  onpointerdown={resize}
  ondblclick={() => {
    if (side === 'story') ws.storyHeight = 34;
  }}
  onkeydown={(event) => {
    const forward = side === 'story' ? 'ArrowDown' : 'ArrowRight';
    const backward = side === 'story' ? 'ArrowUp' : 'ArrowLeft';
    if (event.key !== forward && event.key !== backward) {
      return;
    }
    event.preventDefault();
    event.stopPropagation();
    const height = event.currentTarget.parentElement?.clientHeight ?? 700;
    adjust(event.key === forward ? 20 : -20, height);
  }}
></button>
