import { mediaHTML } from './media-report';
import type { Project } from '../types/project';
import { intervalLabel } from '../types/temporal';

function escape(value: string): string {
  return value.replaceAll('&', '&amp;').replaceAll('<', '&lt;').replaceAll('>', '&gt;').replaceAll('"', '&quot;');
}

export function reportHTML(project: Project): string {
  const chapters = project.chapters.map((chapter) => `<section><p class="label">${chapter.interpretation ? 'Interpretation' : 'Documented chronology'}</p><h2>${escape(chapter.title)}</h2><p>${escape(chapter.summary)}</p><p class="references">Claims: ${escape(chapter.claimIds.join(', '))}</p></section>`).join('');
  const events = project.events.map((event) => {
    const claims = project.claims.filter((claim) => event.claimIds.includes(claim.id));
    const assertions = claims.map((claim) => {
      const evidence = project.evidence.filter((item) => item.claimId === claim.id).map((item) => `<blockquote>${escape(item.excerpt)}<footer>${escape(item.mode)}; ${escape(item.relation)}; ${escape(item.location)}; source ${escape(item.sourceId)}; verified ${escape(String(item.verified))}</footer></blockquote>`).join('');
      return `<article><h4>${escape(claim.text)}</h4><p class="label">${escape(claim.status.replaceAll('_', ' '))} · ${escape(claim.id)}</p>${evidence}<p class="references">${escape(claim.assessment.join(' '))}</p></article>`;
    }).join('');
    return `<section><h3>${escape(intervalLabel(event.time))} · ${escape(event.title)}</h3><p>${escape(event.narration?.text || event.summary)}</p><p class="references">Precision: ${escape(event.time.start?.precision ?? 'unknown')}; calendar: ${escape(event.time.start?.calendar ?? 'unknown')}; temporal record: ${escape(JSON.stringify(event.time))}</p>${assertions}</section>`;
  }).join('');
  const contexts = (project.contextPeriods ?? []).map((period) => `<section><h3>${escape(period.title)}</h3><p>${escape(period.scope)}; ${escape(intervalLabel(period.time))}</p><p>${escape(period.summary)}</p><p>Claims: ${escape(period.claimIds.join(", "))}</p></section>`).join('');
  const conflicts = project.conflicts.map((conflict) => `<section><h3>${escape(conflict.title)}</h3><p>${escape(conflict.explanation)}</p><ul>${conflict.claimIds.map((id) => `<li>${escape(project.claims.find((claim) => claim.id === id)?.text ?? id)}</li>`).join('')}</ul><p>${escape(conflict.consequences.join(' '))}</p></section>`).join('');
  const sources = project.sources.map((source) => `<li id="${escape(source.id)}">${escape(source.author || 'Author unknown')}. <em>${escape(source.title)}</em>. ${escape(source.publisher || 'Publisher unknown')}. ${escape(source.url)}. Retrieved ${escape(source.retrievedAt || 'not recorded')}. ${escape(source.classification)}; independence group ${escape(source.independenceGroup)}.</li>`).join('');
  const coverage = project.coverage.map((period) => `<li>${escape(String(period.startYear))}–${escape(String(period.endYear))}: ${escape(period.status)}. ${escape(period.question)}. ${escape(period.notes)}</li>`).join('');
  return `<!doctype html><html lang="en"><head><meta charset="utf-8"><title>${escape(project.title)} | Chronicle</title><style>
body{max-width:820px;margin:3rem auto;padding:0 2rem;font:15px/1.7 Georgia,serif;color:#243044}
h1{font-size:27px;line-height:1.1}h2{margin-top:2rem}h3{margin-top:1.5rem}h4{margin:1rem 0 .2rem}
.label{font:13px/1.6 system-ui,sans-serif;letter-spacing:1px;text-transform:uppercase;color:#64748b}
.references,blockquote footer{font:13px/1.6 system-ui,sans-serif;color:#66758a;overflow-wrap:anywhere}
blockquote{border-left:2px solid #b0bdce;padding:0 1rem;margin:1rem 0}section{break-inside:avoid}li{margin-bottom:.7rem;overflow-wrap:anywhere}
@media print{body{margin:0;max-width:none;font-size:11pt}h1{font-size:20pt}a{color:inherit}}
</style></head><body><p class="label">Chronicle · Research record · Format v1</p><h1>${escape(project.title)}</h1><p>${project.fixture ? 'Curated fixture; historical coverage is incomplete. Synthetic subjects are labeled in their titles.' : 'Research findings are limited to the evidence retrieved. Alternative claims and coverage gaps remain available.'}</p>${chapters}<h2>Historical context</h2>${contexts}<h2>Chronology and evidence</h2>${events}<h2>Competing interpretations</h2>${conflicts || '<p>No recorded disagreements. This does not establish consensus.</p>'}${mediaHTML(project, escape)}<h2>Bibliography</h2><ol>${sources}</ol><h2>Coverage and limitations</h2><ul>${coverage}</ul><h2>Research annotations</h2>${project.notes.map((note) => `<p>${escape(note.targetId)}: ${escape(note.text)}</p>`).join('')}</body></html>`;
}
