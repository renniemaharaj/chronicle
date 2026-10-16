<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';
  import { Pause, Play, Square, Check, Activity } from '@lucide/svelte';
  import { api } from '../../api/client';
  import { workspace as ws } from '../../state/workspace.svelte';

  async function control(action: string) {
    if (!ws.run) return;
    try {
      ws.run = await api.control(ws.run.id, action);
    } catch (error) {
      ws.notify(String(error));
    }
  }
</script>

{#if ws.run}
  <div class="run-status">
    <span class="live-dot" class:inactive={ws.run.status !== 'running'}></span><strong>{ws.run.phase}</strong><span
      >{ws.run.status}</span
    >
  </div>
  <p class="run-question">{ws.run.question || 'Preparing a research plan'}</p>
  <div class="run-metrics">
    <div><strong>{ws.run.requests}</strong><span>requests used</span></div>
    <div><strong>{ws.run.sources}</strong><span>sources retrieved</span></div>
    <div><strong>{ws.run.claims}</strong><span>claims extracted</span></div>
    <div><strong>{ws.run.conflicts}</strong><span>conflicts retained</span></div>
  </div>
  {#if ws.run.error}<p class="run-error">{ws.run.error}</p>{/if}
  <div class="research-plan">
    {#each ws.run.plan as question}<div>
        <span class:done={question.status === 'completed'}><Check size={12} /></span><strong
          >{question.question}<small>{question.interval} · {question.status}</small></strong
        >
      </div>{/each}
  </div>
  <div class="tool-log">
    {#each ws.run.activities.slice(-8).reverse() as activity}<div>
        <time>{activity.time.slice(11, 19)}</time><span>{activity.message}</span>
      </div>{/each}
  </div>
  <div class="dialog-actions">
    {#if ws.run.status === 'running'}<Button variant="secondary" onclick={() => control('pause')}
        ><Pause size={13} /> Pause</Button
      ><Button variant="danger" onclick={() => control('cancel')}><Square size={12} /> Cancel</Button>{/if}
    {#if ['paused', 'failed', 'cancelled'].includes(ws.run.status)}<Button
        variant="primary"
        onclick={() => control('resume')}><Play size={13} /> Resume</Button
      >{/if}
    {#if ws.run.project}<Button variant="primary" onclick={() => ws.open(ws.run!.project!)}>Explore research graph</Button
      >{/if}
  </div>
{/if}
