export function episodeDescription(campaign, part, project) {
  const text = part.passages.map((passage) => passage.text).join('');
  const claims = new Set(part.passages.flatMap((passage) => passage.claimIds));
  for (const period of project.contextPeriods ?? []) {
    if (text.includes(period.summary)) {
      period.claimIds.forEach((id) => claims.add(id));
    }
  }

  const evidence = new Set(
    project.claims.filter((claim) => claims.has(claim.id)).flatMap((claim) => claim.evidenceIds),
  );
  const sources = new Set(
    project.evidence.filter((item) => evidence.has(item.id)).map((item) => item.sourceId),
  );
  const references = campaign.sources.filter((source) => sources.has(source.id));
  const title = `Rome | Part ${String(part.number).padStart(2, '0')} | ${part.passages[0].title}`.slice(
    0,
    100,
  );
  const events = [...new Set(part.passages.map((passage) => passage.title))];
  const visuals = [
    ...new Map(
      part.passages.filter((item) => item.visual).map((item) => [item.visual.image.id, item.visual.image]),
    ).values(),
  ];
  const description = [
    title,
    '',
    `The Rome story, in sequence. This episode: ${events.join('; ')}.`,
    'Dates, historical traditions and competing accounts retain their qualifications in the narration.',
    '',
    'Sources:',
    ...references.map((source) => `${source.title}\n${source.url}`),
    '',
    'Map: OpenStreetMap contributors. Modern geographic references; historical borders are not shown.',
    'https://www.openstreetmap.org/copyright',
    '',
    'AI-generated narration: Bill (ElevenLabs), reused from saved Chronicle audio.',
    ...visuals.map(
      (image) =>
        `${image.generated ? 'AI-generated interpretive illustration' : 'Source image'}: ${image.title}. ${image.caption}`,
    ),
    '#Rome #RomanHistory #History #Shorts',
  ].join('\n');

  return {
    title,
    description,
    eventIds: [...new Set(part.passages.map((passage) => passage.eventId))],
    projectHash: campaign.projectHash,
    captionTiming: 'Transitions follow detected narration pauses where available',
    visual: 'Recorded Chronicle story map with Follow location enabled',
  };
}
