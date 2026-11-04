import { request } from '../api/client';
import type { Action, Campaign } from './types';

export const campaignAPI = {
  list: () => request<Campaign[]>('/campaigns'),
  get: (id: string) => request<Campaign>(`/campaigns/${id}`),
  create: (projectId: string, seconds: number) =>
    request<Campaign>('/campaigns', {
      method: 'POST',
      body: JSON.stringify({ projectId, seconds }),
    }),
  update: (campaign: Campaign, part: number, action: Action) =>
    request<Campaign>(`/campaigns/${campaign.id}/actions`, {
      method: 'POST',
      body: JSON.stringify({
        id: campaign.id,
        revision: campaign.revision,
        part,
        ...action,
      }),
    }),
};
