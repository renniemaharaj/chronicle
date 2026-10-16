import type { Project } from '../types/project';

export function mediaMarkdown(project: Project): string[] {
  const lines = ['## Images and external videos', ''];
  for (const asset of project.media ?? []) {
    lines.push(`- ${asset.kind}: ${asset.title} — ${asset.url}`,
      `  - Subject: ${asset.entityId}; source: ${asset.sourceId}; retrieved link: ${asset.retrieved}.`,
      `  - ${asset.caption} Credit: ${asset.creator || 'Unknown'}; license: ${asset.license || 'Not established'} ${asset.licenseUrl || ''}.`,
      `  - Identifying claims: ${asset.claimIds.join(', ')}.`);
  }
  lines.push('', '## Websites and social references', '');
  for (const link of project.subjectLinks ?? []) {
    lines.push(`- ${link.label}: ${link.url}; subject: ${link.entityId}; verified reference: ${link.verified}.`,
      `  - Method: ${link.verificationMethod}; checked: ${link.verifiedAt || 'Not recorded'}; source: ${link.sourceId}; claims: ${link.claimIds.join(', ')}.`);
  }
  return lines;
}

export function mediaHTML(project: Project, escape: (value: string) => string) {
  return `<h2>Media and public references</h2><pre style="white-space:pre-wrap;overflow-wrap:anywhere">${escape(mediaMarkdown(project).join('\n'))}</pre>`;
}
