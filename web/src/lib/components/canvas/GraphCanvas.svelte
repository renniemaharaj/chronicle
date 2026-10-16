<script lang="ts">
  import { Background, Controls, MiniMap, SvelteFlow, useSvelteFlow, type Node } from '@xyflow/svelte';
  import TemporalNode from '../../TemporalNode.svelte';
  import { workspace as ws } from '../../state/workspace.svelte';
  import { projectGraph } from '../../graph/projector';
  import type { CardData } from '../../graph/types';

  const nodeTypes = { chronicle: TemporalNode };
  const projection = $derived(
    projectGraph(ws.project, {
      resolution: ws.resolution,
      from: ws.startYear,
      to: ws.endYear,
      relations: ws.showRelations,
      conflicts: ws.showConflicts,
      conflictMode: ws.conflictMode,
      labels: ws.showLabels,
      selectedId: ws.selectedId,
      activeClaims: ws.activeClaims,
      view: ws.view,
    }),
  );
  let nodes = $state<Node<CardData>[]>([]);
  let edges = $state<import('@xyflow/svelte').Edge[]>([]);
  const flow = useSvelteFlow();

  $effect(() => {
    nodes = projection.nodes;
    edges = projection.edges;
  });

  $effect(() => {
    const token = ws.fitToken;
    const resolution = ws.resolution;
    const mode = ws.conflictMode;
    const timer = setTimeout(() => flow?.fitView({ padding: 0.12, duration: 450 }), 90);
    return () => clearTimeout(timer);
  });

  function select(node: Node<CardData>) {
    if (node.data.recordType === 'conflict') {
      ws.inspectorTab = 'Conflicts';
      ws.rightOpen = true;
      return;
    }
    ws.select(node.data.recordId, node.data.recordType);
    if (node.data.recordType === 'evidence') ws.inspectorTab = 'Sources';
  }
</script>

<div class="graph-canvas" aria-label="Temporal knowledge graph">
  <div class="graph-legend">
    <div><i class="legend-dot blue-dot"></i>Life events</div>
    <div><i class="legend-dot orange-dot"></i>People</div>
    <div><i class="legend-dot green-dot"></i>Places</div>
    <div><i class="legend-dot gold-dot"></i>Works & ideas</div>
    <div><i class="legend-dot violet-dot"></i>Institutions</div>
  </div>
  <SvelteFlow
    bind:nodes
    bind:edges
    {nodeTypes}
    fitView
    fitViewOptions={{ padding: 0.12 }}
    minZoom={0.15}
    maxZoom={2}
    nodesConnectable={false}
    onlyRenderVisibleElements={true}
    proOptions={{ hideAttribution: true }}
    oninit={() => flow.fitView({ padding: 0.12 })}
    onnodeclick={({ node }) => select(node as Node<CardData>)}
    onnodedragstop={({ nodes: moved }) => {
      for (const node of moved) ws.project.layout[node.id] = node.position;
    }}
  >
    <Background patternColor={ws.dark ? '#273347' : '#c6d1de'} gap={24} size={0.7} />
    <Controls
      position="bottom-right"
      buttonBgColor="var(--surface-raised)"
      buttonColor="var(--text)"
      buttonBorderColor="var(--line)"
      showLock={false}
    />
    <MiniMap
      position="bottom-right"
      width={112}
      height={72}
      bgColor="var(--surface-raised)"
      maskColor={ws.dark ? '#07101b99' : '#eaf0f699'}
      maskStrokeColor="var(--line-strong)"
      pannable
      zoomable
      nodeColor={(node) =>
        node.data.color === 'gold'
          ? 'var(--gold)'
          : node.data.color === 'red'
            ? 'var(--red)'
            : node.data.color === 'teal'
              ? 'var(--teal)'
              : 'var(--blue)'}
    />
  </SvelteFlow>
  {#if ws.showContext}
    <div class="context-ribbon">
      <span>Historical context · {ws.subject.kind}</span><span
        >{ws.project.fixture ? 'Curated example · coverage incomplete' : 'Evidence before narrative'}</span
      >
    </div>
  {/if}
  {#if projection.hidden}<div class="graph-limit">
      {projection.hidden} nodes outside display budget. Narrow the period.
    </div>{/if}
</div>
