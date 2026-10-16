import type { Project } from '../types/project';
import { intervalLabel } from '../types/temporal';
import { reportHTML, reportMarkdown } from './reports';

export const exportFormats = [
  'Project JSON',
  'Markdown report',
  'HTML report',
  'PDF report',
  'Event CSV',
  'Graph JSON',
  'Bibliography',
  'Conflict report',
];

export function exportProject(project: Project, format: string) {
  let text = '';
  let extension = 'json';
  let type = 'application/json';
  if (format === 'Project JSON') text = JSON.stringify(project, null, 2);
  if (format === 'Markdown report') {
    text = reportMarkdown(project);
    extension = 'md';
    type = 'text/markdown';
  }
  if (format === 'HTML report' || format === 'PDF report') {
    text = reportHTML(project);
    extension = 'html';
    type = 'text/html';
  }
  if (format === 'PDF report') {
    const print = window.open('', '_blank');
    if (!print) throw new Error('Allow a new window to print this report to PDF.');
    print.document.write(text);
    print.document.close();
    print.focus();
    print.print();
    return;
  }
  if (format === 'Event CSV') {
    const rows = [
      [
        'id',
        'title',
        'date',
        'precision',
        'calendar',
        'approximate',
        'uncertaintyYears',
        'summary',
        'claimIds',
        'temporalValue',
        'sourceUrls',
      ],
    ];
    for (const event of project.events) {
      const date = event.time.start;
      rows.push([
        event.id,
        event.title,
        intervalLabel(event.time),
        date?.precision ?? 'unknown',
        date?.calendar ?? '',
        String(date?.approximate ?? false),
        String(date?.uncertaintyYears ?? 0),
        event.summary,
        event.claimIds.join(';'),
        JSON.stringify(event.time),
        project.sources
          .filter((source) =>
            project.evidence.some((item) => item.sourceId === source.id && event.claimIds.includes(item.claimId)),
          )
          .map((source) => source.url)
          .join(';'),
      ]);
    }
    text = rows.map((row) => row.map((cell) => `"${cell.replaceAll('"', '""')}"`).join(',')).join('\n');
    extension = 'csv';
    type = 'text/csv';
  }
  if (format === 'Graph JSON')
    text = JSON.stringify(
      {
        formatVersion: 1,
        entities: project.entities,
        events: project.events,
        relationships: project.relationships,
        conflicts: project.conflicts,
        claims: project.claims,
        evidence: project.evidence,
        sources: project.sources,
      },
      null,
      2,
    );
  if (format === 'Bibliography') {
    text = project.sources
      .map(
        (source) =>
          `${source.author || 'Author unknown'}. ${source.title}. ${source.publisher || 'Publisher unknown'}. ${source.url}\nClassification: ${source.classification}. Retrieved: ${source.retrievedAt || 'unknown'}. Independence group: ${source.independenceGroup}.`,
      )
      .join('\n\n');
    extension = 'txt';
    type = 'text/plain';
  }
  if (format === 'Conflict report') {
    const claims = project.claims.filter((claim) =>
      project.conflicts.some((conflict) => conflict.claimIds.includes(claim.id)),
    );
    text = JSON.stringify(
      {
        formatVersion: 1,
        project: project.title,
        conflicts: project.conflicts,
        claims,
        evidence: project.evidence.filter((evidence) => claims.some((claim) => claim.id === evidence.claimId)),
        sources: project.sources,
      },
      null,
      2,
    );
  }

  const url = URL.createObjectURL(new Blob([text], { type }));
  const link = document.createElement('a');
  link.href = url;
  link.download = `chronicle-${project.id}-${format.toLowerCase().replaceAll(' ', '-')}.${extension}`;
  link.click();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}
