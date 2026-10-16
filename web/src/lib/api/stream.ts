import { workspace as ws } from '../state/workspace.svelte';
import type { ResearchRun } from '../types/project';

export function watchRun(id: string) {
  const stream = new EventSource(`/api/research/${id}/events`);
  stream.onmessage = (event) => {
    ws.run = JSON.parse(event.data) as ResearchRun;
    const index = ws.runs.findIndex((run) => run.id === id);
    if (index >= 0) ws.runs[index] = ws.run;
    else ws.runs.unshift(ws.run);
    if (ws.run.project) {
      const projectIndex = ws.projects.findIndex((project) => project.id === ws.run?.project?.id);
      if (projectIndex < 0) ws.projects.push(ws.run.project);
      else ws.projects[projectIndex] = ws.run.project;
    }
    if (['completed', 'failed', 'paused', 'cancelled'].includes(ws.run.status)) stream.close();
  };
  stream.onerror = () => {
    stream.close();
    ws.notify('Research updates disconnected. Reopen the session to reconnect.');
  };
  return () => stream.close();
}
