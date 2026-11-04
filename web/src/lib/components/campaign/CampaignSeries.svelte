<script lang="ts">
  import { TabItem, Tabs } from 'flowbite-svelte';
  import type { CampaignWorkspace } from '../../campaign/workspace.svelte';
  import SeriesSummary from './SeriesSummary.svelte';
  import EpisodesTable from './EpisodesTable.svelte';
  import SeriesOverview from './SeriesOverview.svelte';
  import SeriesDistribution from './SeriesDistribution.svelte';
  import ClipsLibrary from './ClipsLibrary.svelte';
  import CampaignAnalytics from './CampaignAnalytics.svelte';

  let { model }: { model: CampaignWorkspace } = $props();
</script>

<Tabs
  class="reel-main-tabs"
  tabStyle="underline"
  divider={false}
  bind:selected={model.tab}
  classes={{ content: 'reel-tab-content' }}
>
  <TabItem
    key="Series"
    title="Series"
  >
    <SeriesSummary {model} />
    <Tabs
      class="reel-section-tabs"
      tabStyle="underline"
      divider={false}
      bind:selected={model.section}
      classes={{ content: 'reel-section-content' }}
    >
      <TabItem
        key="Episodes"
        title={`Episodes (${model.studio.campaign?.parts.length ?? 0})`}
      >
        <EpisodesTable {model} />
      </TabItem>
      <TabItem
        key="Overview"
        title="Overview"
      >
        <SeriesOverview campaign={model.studio.campaign!} />
      </TabItem>
      <TabItem
        key="Distribution"
        title="Distribution"
      >
        <SeriesDistribution {model} />
      </TabItem>
    </Tabs>
  </TabItem>
  <TabItem
    key="Clips Library"
    title="Clips Library"
  >
    <ClipsLibrary {model} />
  </TabItem>
  <TabItem
    key="Analytics"
    title="Analytics"
  >
    <CampaignAnalytics {model} />
  </TabItem>
</Tabs>
