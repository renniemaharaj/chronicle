<script lang="ts">
  import { Handle, Position, type NodeProps } from '@xyflow/svelte';
  import { BookOpen, Landmark, MapPin, Sparkles, Users } from '@lucide/svelte';

  type TemporalNodeData = {
    title: string;
    detail: string;
    year: string;
    kind: string;
    color: string;
    portrait?: string;
    count?: number;
    evidenceCount?: number;
  };

  let { data, selected } = $props<NodeProps & { data: TemporalNodeData; selected: boolean }>();

  const iconByKind: Record<string, typeof BookOpen> = {
    Person: Users,
    Place: MapPin,
    Organization: Landmark,
    'Scientific work': BookOpen,
    Publication: BookOpen,
    Idea: Sparkles,
    Disputed: Sparkles,
    'Historical phase': Landmark,
  };

  const NodeIcon = $derived(iconByKind[data.kind] ?? BookOpen);
</script>

<div
  class="temporal-node"
  class:chosen={selected}
  class:gold={data.color === 'gold'}
  class:violet={data.color === 'violet'}
  class:teal={data.color === 'teal'}
  class:conflicted={data.color === 'red'}
>
  <Handle type="target" position={Position.Left} />
  <div class="node-topline">
    <span class="node-icon"
      >{#if data.portrait}<img src={data.portrait} alt="" />{:else}<NodeIcon size={14} />{/if}</span
    >
    <span class="node-kind">{data.kind}</span>
  </div>
  <strong>{data.title}</strong>
  <small
    >{data.year}<span
      >{data.count && data.count > 1
        ? `${data.count} events`
        : data.evidenceCount
          ? `${data.evidenceCount} source`
          : ''}</span
    ></small
  >
  <Handle type="source" position={Position.Right} />
</div>
