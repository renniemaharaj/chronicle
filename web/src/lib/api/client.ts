import type { Project, ProviderSettings, ResearchRun } from '../types/project';

export async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`/api${path}`, {
    ...options,
    headers: { 'Content-Type': 'application/json', ...options?.headers },
  });
  const body = await response.json().catch(() => ({ error: 'The local Chronicle API is unavailable.' }));
  if (!response.ok) {
    throw new Error(body.error ?? `Request failed (${response.status})`);
  }

  return body as T;
}

export const api = {
  projects: () => request<Project[]>('/projects'),
  save: (project: Project) =>
    request<Project>(`/projects/${project.id}`, { method: 'PUT', body: JSON.stringify(project) }),
  settings: () => request<ProviderSettings>('/settings'),
  saveSettings: (settings: ProviderSettings) =>
    request<ProviderSettings>('/settings', { method: 'PUT', body: JSON.stringify(settings) }),
  tools: () => request<{ name: string; description: string }[]>('/tools'),
  start: (subject: string) => request<ResearchRun>('/research', { method: 'POST', body: JSON.stringify({ subject }) }),
  runs: () => request<ResearchRun[]>('/research'),
  run: (id: string) => request<ResearchRun>(`/research/${id}`),
  control: (id: string, action: string) => request<ResearchRun>(`/research/${id}/${action}`, { method: 'POST' }),
};
