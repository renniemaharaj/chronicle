import type { Project, ResearchRun } from '../types/project';

export function completion(done: number, total: number) {
  return total ? Math.floor((done / total) * 100) : 0;
}

export function researchCompletion(run: Pick<ResearchRun, 'plan' | 'status'>) {
  const done = run.plan.filter((question) => question.status === 'completed').length;
  const percent = completion(done, run.plan.length);

  return {
    done,
    total: run.plan.length,
    percent: run.status === 'completed' ? percent : Math.min(99, percent),
  };
}

export function storyCompletion(project: Project) {
  const supported = new Set(
    project.claims
      .filter(
        (claim) =>
          ['supported_limited', 'corroborated', 'disputed'].includes(claim.status) &&
          project.evidence.some(
            (item) =>
              item.claimId === claim.id &&
              item.verified &&
              ['supporting', 'partially_supporting'].includes(item.relation),
          ),
      )
      .map((claim) => claim.id),
  );
  const evidence = project.events.filter((event) => event.claimIds.some((id) => supported.has(id))).length;
  const narration = project.events.filter(
    (event) =>
      event.narration?.text.trim() &&
      event.narration.claimIds.length &&
      event.narration.claimIds.every((id) => supported.has(id)),
  ).length;
  const reviewedAt = project.researchRecords
    ?.flatMap((record) => [record.createdAt, ...record.activities.map((activity) => activity.time)])
    .sort()
    .at(-1);
  const latest = project.researchRecords?.at(-1);

  return {
    total: project.events.length,
    evidence,
    narration,
    evidencePercent: completion(evidence, project.events.length),
    narrationPercent: completion(narration, project.events.length),
    reviewedAt,
    research: latest?.plan.length ? researchCompletion(latest) : null,
  };
}
