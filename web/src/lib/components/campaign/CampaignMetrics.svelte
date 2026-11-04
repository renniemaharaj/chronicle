<script lang="ts">
  import Button from '../ui/Button.svelte';
  import type { Studio } from '../../campaign/studio.svelte';
  import type { Metrics } from '../../campaign/types';
  import type { Platform } from '../../campaign/platforms';

  let { studio, platform = 'facebook' }: { studio: Studio; platform?: Platform } = $props();
  let views = $state<number>();
  let reactions = $state<number>();
  let comments = $state<number>();
  let shares = $state<number>();
  let revenue = $state<number>();
  let currency = $state('USD');

  async function save() {
    studio.busy = true;

    const metrics: Metrics = {
      observedAt: new Date().toISOString(),
      views,
      reactions,
      comments,
      shares,
      revenue,
      currency,
      method: 'manual',
    };
    try {
      await studio.act({
        action: 'metrics',
        platform,
        metrics,
      });
      views = reactions = comments = shares = revenue = undefined;
    } catch (error) {
      studio.error = (error as Error).message;
    } finally {
      studio.busy = false;
    }
  }
</script>

<h3>Performance observations</h3>
<p>
  Enter values you observed on {platform === 'youtube' ? 'YouTube' : 'Facebook'}. Leave unavailable values
  blank. Revenue is an observed amount, never a forecast.
</p>
<div class="campaign-metrics">
  <label>
    Views <input
      type="number"
      min="0"
      bind:value={views}
    />
  </label>
  <label>
    Reactions <input
      type="number"
      min="0"
      bind:value={reactions}
    />
  </label>
  <label>
    Comments <input
      type="number"
      min="0"
      bind:value={comments}
    />
  </label>
  <label>
    Shares <input
      type="number"
      min="0"
      bind:value={shares}
    />
  </label>
  <label>
    Revenue <input
      type="number"
      min="0"
      step="0.01"
      bind:value={revenue}
    />
  </label>
  <label>
    Currency <input
      maxlength="3"
      bind:value={currency}
    />
  </label>
</div>
<Button
  disabled={studio.busy}
  onclick={save}
>
  Save observed performance
</Button>
{#each (studio.part?.metrics ?? []).filter((metric) => (metric.platform ?? 'facebook') === platform) as metric}
  <div class="campaign-observation">
    <small>{new Date(metric.observedAt).toLocaleString()} · manual observation</small>
    <p>
      Views {metric.views ?? 'unknown'} · Reactions {metric.reactions ?? 'unknown'} · Comments {metric.comments ??
        'unknown'} · Shares {metric.shares ?? 'unknown'}
    </p>
    <p>Revenue {metric.revenue === undefined ? 'unknown' : `${metric.currency} ${metric.revenue}`}</p>
  </div>
{/each}
