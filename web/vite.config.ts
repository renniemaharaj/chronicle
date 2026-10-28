import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig(({ mode }) => {
  const environment = loadEnv(mode, '..', 'CHRONICLE_');

  const value = (name: string, fallback: string) => process.env[name] ?? environment[name] ?? fallback;

  return {
    plugins: [tailwindcss(), sveltekit()],
    server: {
      host: value('CHRONICLE_DEV_HOST', '127.0.0.1'),
      port: Number(value('CHRONICLE_DEV_PORT', '5173')),
      strictPort: true,
      proxy: {
        '/api': value('CHRONICLE_DEV_API_URL', 'http://127.0.0.1:8080'),
      },
    },
  };
});
