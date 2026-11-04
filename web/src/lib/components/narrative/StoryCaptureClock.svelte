<script lang="ts">
  import { onMount } from 'svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { captureClock } from '../../state/story-presentation.svelte';

  onMount(() => {
    if (new URLSearchParams(location.search).get('capture') !== 'story') {
      return;
    }

    const update = (event: Event) => {
      const input = (event as CustomEvent).detail;
      if (
        input?.eventId !== ws.event?.id ||
        !Number.isFinite(input.seconds) ||
        input.seconds < 0 ||
        input.seconds > 3600 ||
        !Number.isFinite(input.duration) ||
        input.duration <= 0 ||
        input.duration > 3600
      ) {
        return;
      }

      Object.assign(captureClock, {
        ...input,
        projectId: ws.project.id,
        active: true,
      });
    };

    window.addEventListener('chronicle:story-capture-clock', update);

    return () => {
      window.removeEventListener('chronicle:story-capture-clock', update);
      captureClock.active = false;
    };
  });
</script>
