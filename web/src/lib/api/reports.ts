import { mediaMarkdown } from './media-report';
import type { Project } from '../types/project';
import { intervalLabel } from '../types/temporal';

export function reportMarkdown(project: Project): string {
  const lines = [
    `# ${project.title}`,
    '',
    project.fixture
      ? 'Curated fixture. Coverage is incomplete; synthetic records are explicitly labeled.'
      : 'Research record. Coverage is limited to retrieved evidence.',
    '',
  ];
  for (const chapter of project.chapters) {
    lines.push(`## ${chapter.title}`, '', `${chapter.interpretation ? 'Interpretation: ' : ''}${chapter.summary}`, '');
  }
  lines.push('## Historical context', '');
  for (const period of project.contextPeriods ?? []) {
    lines.push(`### ${period.title}`, '', `${period.scope}; ${intervalLabel(period.time)}`, '', period.summary, '', `Claims: ${period.claimIds.join(', ')}`, '');
  }
  lines.push('## Chronology', '');
  for (const event of project.events) {
    lines.push(
      `### ${intervalLabel(event.time)}: ${event.title}`,
      '',
      event.narration?.text || event.summary,
      '',
      `Calendar: ${event.time.start?.calendar ?? 'unknown'}; precision: ${event.time.start?.precision ?? 'unknown'}.`,
      '',
    );
    for (const claim of project.claims.filter((claim) => event.claimIds.includes(claim.id))) {
      lines.push(`- Claim ${claim.id}: ${claim.text} [${claim.status}]`);
      for (const evidence of project.evidence.filter((item) => item.claimId === claim.id)) {
        lines.push(
          `  - ${evidence.mode}: ${evidence.excerpt} (${evidence.relation}; ${evidence.sourceId}; ${evidence.location}; verified: ${evidence.verified})`,
        );
      }
    }
    lines.push('');
  }
  lines.push('## Competing claims', '');
  for (const conflict of project.conflicts) {
    lines.push(
      `### ${conflict.title}`,
      '',
      conflict.explanation,
      '',
      ...conflict.claimIds.map(
        (id) => `- ${id}: ${project.claims.find((claim) => claim.id === id)?.text ?? 'Missing claim'}`,
      ),
      '',
      ...conflict.consequences.map((item) => `- Implication: ${item}`),
      '',
    );
  }
  lines.push(...mediaMarkdown(project), '', '## Sources', '');
  for (const source of project.sources)
    lines.push(
      `- [${source.title}](${source.url}) — ${source.author}; ${source.publisher}. ${source.classification}; retrieved ${source.retrievedAt}; independence group ${source.independenceGroup}.`,
    );
  lines.push(
    '',
    '## Coverage and limitations',
    '',
    ...project.coverage.map(
      (period) => `- ${period.startYear}–${period.endYear}: ${period.status}. ${period.question}. ${period.notes}`,
    ),
  );
  lines.push('', '## Research annotations', '', ...project.notes.map((note) => `- ${note.targetId}: ${note.text}`));

  return lines.join('\n');
}

export { reportHTML } from './html-report';
