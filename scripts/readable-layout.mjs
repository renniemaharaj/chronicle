import { expandObject } from './readable-objects.mjs';

export function javascriptLayout(ts, source) {
  const file = ts.createSourceFile('source.ts', source, ts.ScriptTarget.Latest, true);
  const edits = new Map();

  function insert(position, text) {
    edits.set(position, (edits.get(position) ?? '') + text);
  }

  function gap(node, lines) {
    const position = node.getStart(file);
    const comments = ts.getLeadingCommentRanges(source, node.getFullStart()) ?? [];
    const start = comments[0]?.pos ?? position;
    const whitespace = source.slice(0, start).match(/\s*$/)?.[0] ?? '';
    const existing = (whitespace.match(/\n/g) ?? []).length;

    if (existing < lines && start > 0) {
      const addition = '\n'.repeat(lines - existing);
      if (!(edits.get(start) ?? '').endsWith(addition)) {
        insert(start, addition);
      }
    }
  }

  function isControl(node) {
    return (
      ts.isIfStatement(node) ||
      ts.isIterationStatement(node, false) ||
      ts.isTryStatement(node) ||
      ts.isSwitchStatement(node)
    );
  }

  function isFunction(node) {
    if (ts.isVariableStatement(node)) {
      return node.declarationList.declarations.some((declaration) => {
        const value = declaration.initializer;

        return value && (ts.isArrowFunction(value) || ts.isFunctionExpression(value));
      });
    }

    return (
      ts.isFunctionDeclaration(node) ||
      ts.isMethodDeclaration(node) ||
      ts.isGetAccessor(node) ||
      ts.isSetAccessor(node)
    );
  }

  function statements(list) {
    list.forEach((node, index) => {
      const previous = list[index - 1];
      if (!previous) {
        return;
      }

      const nextPhase = ts.isExpressionStatement(previous) && ts.isVariableStatement(node);
      if (
        isControl(previous) ||
        isFunction(previous) ||
        isFunction(node) ||
        ts.isReturnStatement(node) ||
        nextPhase
      ) {
        gap(node, 2);
      }
    });
  }

  function expand(statement, elseIf = false) {
    if (!statement || ts.isBlock(statement) || (elseIf && ts.isIfStatement(statement))) {
      return;
    }

    insert(statement.getStart(file), '{\n');
    insert(statement.end, '\n}');
  }

  function visit(node) {
    expandObject(ts, node, source, insert);

    if (ts.isIfStatement(node)) {
      expand(node.thenStatement);
      expand(node.elseStatement, true);
    } else if (ts.isIterationStatement(node, false)) {
      expand(node.statement);
    }

    if (ts.isBlock(node) || ts.isSourceFile(node) || ts.isCaseClause(node) || ts.isDefaultClause(node)) {
      statements(node.statements);
    } else if (ts.isClassDeclaration(node)) {
      statements(node.members);
    }

    ts.forEachChild(node, visit);
  }

  visit(file);

  let result = source;
  const positions = [...edits.keys()].sort((left, right) => right - left);
  for (const position of positions) {
    result = result.slice(0, position) + edits.get(position) + result.slice(position);
  }

  return result;
}
