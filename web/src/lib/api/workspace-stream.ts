import { workspace as ws } from '../state/workspace.svelte';
import type { Project } from '../types/project';

// Backend changes, including assistant-created projects, appear without a page reload.
export function watchWorkspace() {
  const stream = new EventSource('/api/workspace/events');
  const versions = new Map<string, string>();
  stream.onmessage = (event) => {
    const projects = (JSON.parse(event.data) as Project[]).filter((project) => !project.fixture);
    ws.connected = true;
    for (const project of projects) {
      const fingerprint = JSON.stringify(project);
      if (versions.get(project.id) === fingerprint) continue;
      versions.set(project.id, fingerprint);
      const index = ws.projects.findIndex((item) => item.id === project.id);
      if (index < 0) {
        ws.projects.push(project);
        if (!ws.hasProject) ws.open(project);
        else ws.notify(`Research available in Library: ${project.title}`);
      } else {
        ws.projects[index] = project;
        if (ws.project.id === project.id) ws.project = project;
      }
    }
  };
  stream.onerror = () => {
    ws.connected = false;
  };
  return () => stream.close();
}
