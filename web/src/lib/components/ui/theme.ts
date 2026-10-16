import type { ThemeConfig } from 'flowbite-svelte';

export const chronicleTheme: ThemeConfig = {
  button: { base: 'gap-2 text-[length:var(--font-control)]' },
  input: {
    input: 'min-w-0 rounded-lg text-[length:var(--font-control)] bg-[var(--surface)] text-[var(--text)] dark:bg-[var(--surface)] dark:text-[var(--text)]',
  },
  select: {
    base: 'min-w-0',
    select: 'pe-9 text-[length:var(--font-control)] bg-[var(--surface)] text-[var(--text)] dark:bg-[var(--surface)] dark:text-[var(--text)]',
  },
  label: 'text-[length:var(--font-control)] text-gray-800 dark:text-gray-100',
  checkbox: { div: 'gap-2 text-[length:var(--font-control)]', base: 'shrink-0' },
  modal: {
    base: 'flex max-h-[88dvh] flex-col border border-gray-200 bg-[var(--surface-raised)] text-[var(--text)] shadow-2xl dark:border-gray-700 dark:bg-[var(--surface-raised)] dark:text-[var(--text)]',
    header: 'border-b border-gray-200 px-6 py-5 dark:border-gray-700',
    body: 'min-h-0 overflow-auto p-6',
  },
  card: { base: 'bg-[var(--surface)] dark:bg-[var(--surface)]' },
};
