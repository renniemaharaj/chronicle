import { browserTool, browserClick, browserWait, browserFill } from './reel-browser.mjs';

export function youtubeState() {
  const dialogs = [...document.querySelectorAll('[role=dialog]')].filter((e) => e.offsetWidth);

  return {
    text: dialogs.map((e) => e.innerText.replace(/\s+/g, ' ')).join('\n'),
    links: dialogs.flatMap((e) => [...e.querySelectorAll('a')].map((a) => a.href)),
    title: document.querySelector(
      '[aria-label="Add a title that describes your video (type @ to mention a channel)"]',
    )?.innerText,
    description: document.querySelector(
      '[aria-label="Tell viewers about your video (type @ to mention a channel)"]',
    )?.innerText,
  };
}

export async function openYouTubeUpload(part) {
  await browserClick('role=button[name="Create"]');
  await browserClick('role=menuitem[name="Upload videos"]');
  await browserClick('role=button[name="Select files"]');
  await browserTool('browser_file_upload', { paths: [part.file] });

  return browserWait(youtubeState, (value) => {
    if (value.text.includes('Daily upload limit reached')) {
      throw new Error('YouTube daily upload limit reached; verification or a limit reset is required.');
    }

    return value.links.some((link) => /shorts\/[\w-]{11}/.test(link));
  });
}

export async function describeYouTube(metadata) {
  await browserFill(
    'role=textbox[name="Add a title that describes your video (type @ to mention a channel)"]',
    metadata.title,
    'Short title',
  );
  await browserFill(
    'role=textbox[name="Tell viewers about your video (type @ to mention a channel)"]',
    metadata.description,
    'Short description',
  );
  await browserClick('role=radio[name="No, it\'s not made for kids"]');
  await browserClick('role=button[name="Show advanced settings"]');
  await browserClick('role=radio[name="Yes, AI was used"]');
  await browserClick('role=radio[name="No, my video doesn’t include paid promotion"]');
}

export async function readyYouTube() {
  await browserClick('role=button[name="Next"]');
  await browserWait(youtubeState, (value) => value.text.includes('Use cards and an end screen'));
  await browserClick('role=button[name="Next"]');
  await browserWait(youtubeState, (value) => value.text.includes('No issues found'), 900);
  await browserClick('role=button[name="Next"]');
  await browserWait(youtubeState, (value) => value.text.includes('Choose when to publish'));
  await browserClick('role=radio[name="Public"]');
}

export async function confirmedYouTube(title, url) {
  const state = await browserWait(youtubeState, (value) => value.text.includes('Video published'), 240);
  if (!state.text.includes(title) || !state.links.some((link) => link.startsWith(url))) {
    throw new Error('The YouTube confirmation does not match this episode.');
  }

  return state;
}
