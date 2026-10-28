export function colorFor(kind: string): string {
  if (['Publication', 'Idea'].includes(kind)) {
    return 'gold';
  }

  if (['Education', 'Collaboration'].includes(kind)) {
    return 'teal';
  }

  if (kind === 'Appointment') {
    return 'violet';
  }

  if (kind === 'Dispute') {
    return 'red';
  }

  return 'blue';
}
