import { expect, it } from 'vitest';
import { safeSourceURL } from './links';

it('permits source links without allowing imported projects to execute URL schemes', () => {
  expect(safeSourceURL('https://example.org/archive')).toBe('https://example.org/archive');
  expect(safeSourceURL('javascript:alert(1)')).toBeUndefined();
  expect(safeSourceURL('data:text/html,<script>alert(1)</script>')).toBeUndefined();
});
