<script lang="ts">
  import { Eye, Heart, MessageCircle, Share2 } from '@lucide/svelte';
  import { metricLabel, totalMetric, type MetricName } from '../../campaign/observations';
  import type { Campaign } from '../../campaign/types';

  let { campaign }: { campaign: Campaign } = $props();
  const fields: { key: MetricName; label: string; icon: typeof Eye }[] = [
    {
      key: 'views',
      label: 'Views',
      icon: Eye,
    },
    {
      key: 'reactions',
      label: 'Reactions',
      icon: Heart,
    },
    {
      key: 'comments',
      label: 'Comments',
      icon: MessageCircle,
    },
    {
      key: 'shares',
      label: 'Shares',
      icon: Share2,
    },
  ];
</script>

<div
  class="reel-series-metrics"
  aria-label="Series performance"
>
  <div class="reel-metrics-label">
    Facebook performance
    <span>Manual observations</span>
  </div>
  <div class="reel-metric-grid">
    {#each fields as field}
      {@const total = totalMetric(campaign, field.key)}
      <div
        title={`${total.observed} of ${campaign.parts.length} episodes have observed ${field.label.toLowerCase()}`}
      >
        <strong><field.icon size={15} />{metricLabel(total.value)}</strong>
        <span>{field.label}</span>
      </div>
    {/each}
  </div>
</div>
