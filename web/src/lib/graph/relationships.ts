import type { Edge, Node } from '@xyflow/svelte';
import type { Project } from '../types/project';
import type { CardData } from './types';

export function projectEdges(project: Project, nodes: Node<CardData>[], labels: boolean, relations: boolean): Edge[] {
  const edges: Edge[] = [];
  const eventNodes = nodes.filter((node) => node.data.recordType === 'event');
  for (let index = 1; index < eventNodes.length; index++) {
    const previous = eventNodes[index - 1];
    const node = eventNodes[index];
    edges.push({
      id: `sequence:${index}`,
      source: previous.id,
      target: node.id,
      type: 'smoothstep',
      class: 'sequence-edge',
      label: labels && index === 1 ? 'chronology, not causation' : undefined,
    });
  }
  if (!relations) return edges;

  for (const node of eventNodes) {
    const participants = new Set(
      project.events.filter((event) => node.data.eventIds.includes(event.id)).flatMap((event) => event.participants),
    );
    for (const entityId of [...participants].slice(0, 3)) {
      const entityNode = nodes.find((item) => item.id === `entity:${entityId}`);
      if (!entityNode) continue;
      edges.push({
        id: `participant:${node.id}:${entityId}`,
        source: entityNode.id,
        target: node.id,
        type: 'default',
        animated: node.selected,
        class: 'participant-edge',
      });
    }
  }
  for (const relation of project.relationships) {
    const source = `entity:${relation.fromId}`;
    const target = `entity:${relation.toId}`;
    if (!nodes.some((node) => node.id === source) || !nodes.some((node) => node.id === target)) continue;
    edges.push({
      id: relation.id,
      source,
      target,
      type: 'default',
      label: labels ? relation.kind : undefined,
      class: 'relationship-edge',
    });
  }

  return edges;
}
