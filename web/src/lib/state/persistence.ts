import { api } from '../api/client';
import { workspace as ws } from './workspace.svelte';
import type { Project } from '../types/project';

export async function initialize() {
  try {
    const cached = localStorage.getItem('chronicle-workspace');
    if (cached) {
      const saved = JSON.parse(cached);
      ws.dark = saved.dark ?? true;
      ws.projects = (saved.projects ?? []).filter((project: Project) => !project.fixture);
      const project = ws.projects.find((project) => project.id === saved.projectId) ?? ws.projects[0];
      if (project) ws.open(project);
      if (saved.selectedId) ws.selectedId = saved.selectedId;
    }
  } catch {
    ws.notify('The browser workspace could not be restored. Import a project or begin new research.');
  }

  try {
    const [projects, runs] = await Promise.all([api.projects(), api.runs()]);
    ws.connected = true;
    ws.runs = runs;
    ws.run = runs.find((run) => ['running', 'pending'].includes(run.status)) ?? ws.run;
    for (const project of projects.filter((item) => !item.fixture)) {
      const index = ws.projects.findIndex((item) => item.id === project.id);
      if (index < 0) ws.projects.push(project);
      else ws.projects[index] = project;
      if (ws.project.id === project.id || !ws.hasProject) ws.open(project);
    }
  } catch {
    ws.connected = false;
  }
}

let cacheWarningShown = false;

export function cacheWorkspace() {
  try {
    localStorage.setItem(
      'chronicle-workspace',
      JSON.stringify({
        projectId: ws.project.id,
        selectedId: ws.selectedId,
        dark: ws.dark,
        projects: ws.projects,
      }),
    );
  } catch {
    if (!cacheWarningShown) {
      cacheWarningShown = true;
      ws.notify('Browser storage is full. Use the local backend or export your project.');
    }
  }
}

export async function saveProject(project = ws.project) {
  if (!project.id) return;
  const index = ws.projects.findIndex((item) => item.id === project.id);
  if (index >= 0) ws.projects[index] = project;
  else ws.projects.push(project);
  cacheWorkspace();

  if (ws.connected) {
    try {
      const saved = await api.save(project);
      ws.project = saved;
      const savedIndex = ws.projects.findIndex((item) => item.id === saved.id);
      if (savedIndex >= 0) ws.projects[savedIndex] = saved;
      ws.notify('Project saved to your local workspace');
      return;
    } catch (error) {
      ws.notify(`Saved in browser. ${String(error)}`);
      return;
    }
  }

  ws.notify('Project saved in this browser');
}

export function importProject(value: unknown): Project {
  const project = value as Project;
  const lists = [
    'entities',
    'events',
    'claims',
    'sources',
    'evidence',
    'conflicts',
    'relationships',
    'chapters',
    'coverage',
    'notes',
  ];
  if (
    project?.formatVersion !== 1 ||
    !project.id ||
    !project.subjectId ||
    !lists.every((key) => Array.isArray((project as unknown as Record<string, unknown>)[key]))
  ) {
    throw new Error('This file is not a Chronicle v1 project.');
  }
  if (!project.entities.some((entity) => entity.id === project.subjectId)) {
    throw new Error('The project subject does not exist in its entity list.');
  }

  project.layout ??= {};
  return project;
}
