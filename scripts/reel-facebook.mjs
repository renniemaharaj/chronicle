import { browserTool, browserClick, browserValue, browserWait, browserFill } from './reel-browser.mjs';

export async function facebookIdentity() {
  await browserTool('browser_navigate', { url: 'https://www.facebook.com/me/' });

  return browserWait(
    () => ({
      id: new URL(location.href).searchParams.get('id'),
    }),
    (value) => Boolean(value.id),
    30,
  );
}

function composerState() {
  const form = document.querySelector('form[aria-label="Reels"], [role=form][aria-label="Reels"]');
  const dialog = form?.closest('[role=dialog]');

  return {
    text: dialog?.innerText ?? '',
    description: form?.querySelector('[role=textbox]')?.innerText,
    ai: form?.querySelector('[role=switch]')?.getAttribute('aria-checked'),
    postEnabled: Boolean(
      dialog?.querySelector('[role=button][aria-label="Post"]:not([aria-disabled="true"])'),
    ),
    nextEnabled: [...document.querySelectorAll('[role=button]')].some(
      (e) => e.innerText === 'Next' && e.getAttribute('aria-disabled') !== 'true',
    ),
  };
}

export async function prepareFacebook(part, metadata) {
  await browserTool('browser_navigate', { url: 'https://www.facebook.com/reels/create/' });
  await browserWait(
    () =>
      [...document.querySelectorAll('[role=button]')].some((element) =>
        element.innerText.includes('Add video'),
      ),
    Boolean,
  );
  await browserClick('role=button[name="Add video or drag and drop"]', 'Attach the recorded episode');
  await browserTool('browser_file_upload', { paths: [part.file] });
  await browserWait(composerState, (value) => value.nextEnabled);
  await browserClick('role=button[name="Next"]');
  await browserWait(composerState, (value) => value.text.includes('Edit reel'));
  await configureFacebook(metadata);
}

export async function configureFacebook(metadata) {
  await browserClick('role=button[name="Next"]');
  await browserWait(composerState, (value) => value.text.includes('Reel settings'));
  await browserFill('role=textbox[name="Describe your reel..."]', metadata.description, 'Reel description');
  await configureFacebookAudience(metadata);
}

export async function configureFacebookAudience(metadata) {
  await browserTool('browser_press_key', { key: 'Escape' });

  let state = await browserWait(composerState, (value) => value.ai !== undefined);
  if (!state.text.includes('Anyone on or off Facebook')) {
    await browserClick('role=button[name="Friends Your Facebook friends"]');
    await browserClick('role=radio[name="Public Anyone on or off Facebook"]');
    await browserClick('role=button[name="Done"]');
  }

  state = await browserValue(composerState);
  if (state.ai !== 'true') {
    await browserClick('role=switch[name="Add AI label"]');
  }

  state = await browserValue(composerState);
  if (
    !state.description?.includes(metadata.title) ||
    !state.text.includes('Anyone on or off Facebook') ||
    state.ai !== 'true'
  ) {
    throw new Error('The reel caption, public audience or synthetic narration label was not confirmed.');
  }

  await browserWait(composerState, (value) => value.postEnabled, 900);
}

export async function facebookReceipt(title, knownURLs, account) {
  const profile = `https://www.facebook.com/profile.php?id=${account.id}&sk=reels_tab`;
  const known = new Set(knownURLs.map((url) => url.match(/\/reel\/(\d+)/)?.[1]));
  const deadline = Date.now() + 480000;
  do {
    await browserTool('browser_navigate', { url: profile });

    const candidates = await browserWait(
      () =>
        [...document.links]
          .filter((e) => e.getAttribute('aria-label') === 'Reel tile preview')
          .map((e) => e.href),
      (links) => links.length > 0,
      45,
    );
    for (const candidate of candidates) {
      const id = candidate.match(/\/reel\/(\d+)/)?.[1];
      if (!id || known.has(id)) {
        continue;
      }

      await browserTool('browser_navigate', { url: candidate });

      const observed = await browserWait(
        () => ({
          url: location.href,
          text: [...document.querySelectorAll('main, [role=main]')].map((e) => e.innerText).join('\n'),
          public: Boolean(document.querySelector('[aria-label="Shared with Public"]')),
        }),
        (value) => value.text.includes(account.name),
        45,
      );
      if (
        observed.text.includes(title.slice(0, 64)) &&
        observed.text.includes(account.name) &&
        observed.public
      ) {
        return `https://www.facebook.com/reel/${id}`;
      }
    }

    await new Promise((accept) => setTimeout(accept, 5000));
  } while (Date.now() < deadline);

  throw new Error('No confirmed public reel found. Reconcile the submitted episode before retrying.');
}
