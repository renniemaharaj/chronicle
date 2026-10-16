import { readFile } from 'node:fs/promises';

const [name = 'list', filename] = process.argv.slice(2);
const endpoint = `${process.env.CHRONICLE_API ?? 'http://127.0.0.1:8080'}/api/mcp`;

try {
  const args = filename ? JSON.parse(await readFile(filename, 'utf8')) : {};
  const response = await fetch(endpoint, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json, text/event-stream',
      'MCP-Protocol-Version': '2025-11-25',
    },
    body: JSON.stringify({
      jsonrpc: '2.0', id: 1,
      method: name === 'list' ? 'tools/list' : 'tools/call',
      params: name === 'list' ? {} : { name, arguments: args },
    }),
    signal: AbortSignal.timeout(60000),
  });
  if (!response.ok) throw new Error(`Chronicle returned HTTP ${response.status}: ${await response.text()}`);
  const result = await response.json();
  if (result.error) throw new Error(JSON.stringify(result.error));
  process.stdout.write(`${JSON.stringify(result.result, null, 2)}\n`);
  if (result.result?.isError) process.exitCode = 1;
} catch (error) {
  process.stderr.write(`${error.message}\nStart the local backend with npm start.\n`);
  process.exitCode = 1;
}
