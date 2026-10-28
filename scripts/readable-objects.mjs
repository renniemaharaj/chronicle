export function expandObject(ts, node, source, insert) {
  if (!ts.isObjectLiteralExpression(node) || node.properties.length < 2) {
    return;
  }

  if (!node.properties.some((property) => ts.isPropertyAssignment(property))) {
    return;
  }

  const last = node.properties.at(-1);
  const closing = node.end - 1;
  if (!source.slice(last.end, closing).trimStart().startsWith(',')) {
    insert(last.end, ',');
  }

  const positions = [...node.properties.map((property) => property.getStart()), closing];
  for (const position of positions) {
    const whitespace = source.slice(0, position).match(/\s*$/)?.[0] ?? '';
    if (!whitespace.includes('\n')) {
      insert(position, '\n');
    }
  }
}
