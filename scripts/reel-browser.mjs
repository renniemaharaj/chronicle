import { execFile } from 'node:child_process';
import { promisify } from 'node:util';
import { randomUUID } from 'node:crypto';
import { appendFile, unlink, writeFile } from 'node:fs/promises';

const execute = promisify(execFile);

// Keep tab operations sequential; separate argument files also isolate independent readers.
export async function browserTool(name, args = {}) {
  const argumentsFile = `.chronicle/social/driver-call-${randomUUID()}.json`;
  await writeFile(argumentsFile, JSON.stringify(args));

  let output;
  try {
    output = await execute(process.execPath, ['scripts/playwright-tool.mjs', name, argumentsFile], {
      maxBuffer: 8 * 1024 * 1024,
      timeout: 130000,
    });
  } catch (error) {
    await appendFile('.chronicle/social/browser-driver.log', error.stdout ?? String(error));
    throw new Error(`${name} failed; inspect .chronicle/social/browser-driver.log`, {
      cause: error.stdout ?? String(error),
    });
  } finally {
    await unlink(argumentsFile);
  }

  const result = JSON.parse(output.stdout);
  await appendFile('.chronicle/social/browser-driver.log', `${output.stdout}\n`);
  if (result.isError) {
    throw new Error(`${name} failed; inspect .chronicle/social/browser-driver.log`);
  }

  if (name === 'browser_tabs' && ['select', 'new'].includes(args.action)) {
    // Bring the real tab forward so background throttling does not stall form updates.
    await browserTool('browser_run_code_unsafe', {
      code: 'async (page) => { await page.bringToFront(); }',
    });
  }

  return result.content
    .filter((item) => item.type === 'text')
    .map((item) => item.text)
    .join('\n');
}

export async function browserValue(fn) {
  const text = await browserTool('browser_evaluate', { function: String(fn) });
  const match = text.match(/### Result\n([\s\S]*?)(?:\n###|$)/);
  if (!match) {
    throw new Error('The browser returned no evaluation result.');
  }

  return JSON.parse(match[1]);
}

export async function browserWait(fn, condition, seconds = 120) {
  const deadline = Date.now() + seconds * 1000;
  let value;
  do {
    value = await browserValue(fn);
    if (condition(value)) {
      return value;
    }

    await new Promise((accept) => setTimeout(accept, 2000));
  } while (Date.now() < deadline);

  throw new Error(`The expected browser state was not reached: ${JSON.stringify(value)}`);
}

export async function browserClick(target, element = target) {
  const deadline = Date.now() + 30000;
  do {
    try {
      return await browserTool('browser_click', { target, element });
    } catch (error) {
      // A missing target means no click occurred. Other failures need reconciliation.
      if (!String(error.cause).includes('does not match any elements') || Date.now() > deadline) {
        throw error;
      }

      await new Promise((accept) => setTimeout(accept, 1000));
    }
  } while (Date.now() <= deadline);

  throw new Error(`The control did not appear: ${element}`);
}

export function browserFill(target, value, name) {
  return browserTool('browser_fill_form', {
    fields: [
      {
        target,
        name,
        type: 'textbox',
        value,
      },
    ],
  });
}
