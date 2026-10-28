export function svelteLayout(compiler, source) {
  const ast = compiler.parse(source, { modern: true });
  const edits = new Set();
  const structural = /^(RegularElement|Component|.*Block|Svelte.*|SlotElement)$/;

  function before(position) {
    if (position > 0 && !/\n\s*$/.test(source.slice(0, position))) {
      edits.add(position);
    }
  }

  function after(position) {
    if (position < source.length && !/^\s*\n/.test(source.slice(position))) {
      edits.add(position);
    }
  }

  function fragment(node) {
    let previous;
    for (const child of node.nodes) {
      if (child.type === 'Text' && !child.data.trim()) {
        continue;
      }

      if (previous && structural.test(previous.type) && structural.test(child.type)) {
        before(child.start);
      }

      previous = child;
    }
  }

  function block(node) {
    const prefix = source.slice(node.start, node.start + 12);
    if (!/^\{(?:#|:else if)/.test(prefix)) {
      return;
    }

    const expressions = [node.test, node.expression, node.context, node.key, ...(node.parameters ?? [])];
    const expressionEnd = Math.max(node.start, ...expressions.filter(Boolean).map((item) => item.end));
    after(source.indexOf('}', expressionEnd) + 1);

    const closing = source.lastIndexOf('{/', node.end);
    if (closing >= node.start && source.indexOf('}', closing) + 1 === node.end) {
      before(closing);
    }

    for (const branch of ['alternate', 'fallback', 'then', 'catch']) {
      const first = node[branch]?.nodes?.[0];
      if (!first || first.elseif) {
        continue;
      }

      const start = source.lastIndexOf('{:', first.start);
      if (start > node.start) {
        before(start);
        after(source.indexOf('}', start) + 1);
      }
    }
  }

  function visit(node) {
    if (!node || typeof node !== 'object') {
      return;
    }

    if (node.type === 'Fragment') {
      fragment(node);
    } else if (node.type?.endsWith('Block')) {
      block(node);
    }

    for (const [key, value] of Object.entries(node)) {
      if (['instance', 'module', 'expression', 'test', 'attributes', 'loc', 'metadata'].includes(key)) {
        continue;
      }

      if (Array.isArray(value)) {
        value.forEach(visit);
      } else if (value && typeof value === 'object') {
        visit(value);
      }
    }
  }

  visit(ast.fragment);

  let result = source;
  for (const position of [...edits].sort((left, right) => right - left)) {
    result = result.slice(0, position) + '\n' + result.slice(position);
  }

  return result;
}
