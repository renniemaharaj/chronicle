export const backend = process.env.CHRONICLE_API ?? 'http://127.0.0.1:8080';

export async function reelRequest(path, body) {
  const response = await fetch(`${backend}/api${path}`, {
    method: body === undefined ? 'GET' : 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  const result = await response.json();
  if (!response.ok) {
    throw new Error(result.error ?? `Chronicle returned ${response.status}`);
  }

  return result;
}

export async function reelAction(campaign, number, action) {
  return reelRequest(`/campaigns/${campaign.id}/actions`, {
    id: campaign.id,
    revision: campaign.revision,
    part: number,
    ...action,
  });
}
