<script lang="ts">
  import { ChartNoAxesColumn, Eye, Heart, MessageCircle, Share2 } from '@lucide/svelte';
  import { latestObservation, type MetricName } from '../../campaign/observations';
  import type { Part } from '../../campaign/types';

  let { part }: { part: Part } = $props();
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

<section class="reel-detail-card">
  <h3><ChartNoAxesColumn size={16} />Facebook performance</h3>
  {#each fields as field}
    {@const observation = latestObservation(part, field.key)}
    <div class="reel-detail-metric">
      <field.icon size={15} />
      <strong>{observation?.[field.key]?.toLocaleString() ?? '—'}</strong>
      <span>{field.label}</span>
    </div>
  {/each}
  <p class="reel-detail-note">
    {part.metrics.length
      ? 'Manual observations. Missing values remain unknown.'
      : 'No performance observations yet.'}
  </p>
</section>
