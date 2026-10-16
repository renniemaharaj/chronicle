import type { Entity } from '../types/knowledge';
import type { Project } from '../types/project';
import { knownYear } from '../types/temporal';
import { safeSourceURL } from '../api/links';

export function searchText(value: string) {
  return value
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase();
}

export function portraitFor(entity: Entity, project: Project) {
  const asset = project.media?.find(
    (item) => item.entityId === entity.id && item.kind === 'image' && item.useAsAvatar && item.retrieved,
  );
  const address = asset?.url ?? String(entity.attributes.portrait ?? '');
  return {
    url: safeSourceURL(address),
    credit: asset
      ? `${asset.creator || 'Creator not recorded'} · ${asset.license || 'Rights not established'} · ${asset.title} · Cropped display`
      : 'Research portrait',
  };
}

export function contemporary(entity: Entity, subject: Entity): 'overlap' | 'unknown' | 'outside' {
  if (entity.kind !== 'person') {
    return 'outside';
  }
  const a = knownYear(entity.time.start);
  const b = knownYear(subject.time.start);
  const endA = knownYear(entity.time.end);
  const endB = knownYear(subject.time.end);
  if (a !== undefined && endB !== undefined && a > endB + (subject.time.end?.uncertaintyYears ?? 0)) {
    return 'outside';
  }
  if (b !== undefined && endA !== undefined && b > endA + (entity.time.end?.uncertaintyYears ?? 0)) {
    return 'outside';
  }
  if (a === undefined || b === undefined) {
    return 'unknown';
  }
  if (endA === undefined && endB === undefined) {
    return 'unknown';
  }
  return 'overlap';
}

export function entityDescription(entity: Entity, project: Project) {
  const relationships = project.relationships.filter(
    (item) =>
      (item.fromId === project.subjectId && item.toId === entity.id) ||
      (item.toId === project.subjectId && item.fromId === entity.id),
  );
  return [...new Set(relationships.map((item) => item.kind.replaceAll('_', ' ')))].join(' · ');
}
