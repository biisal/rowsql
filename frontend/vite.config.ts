import react from '@vitejs/plugin-react-swc';
import path from 'path';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig , type PluginOption } from 'vite';
import { heyApiPlugin } from '@hey-api/vite-plugin';

// import { visualizer } from 'rollup-plugin-visualizer';
// https://vite.dev/config/
export default defineConfig({
  // base: '/admin',
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor-react': ['react', 'react-dom'],
          'vendor-utils': ['axios', 'zod', '@tanstack/react-query'],
        },
      },
    },
  },
  plugins: [
    react(),
    tailwindcss(),
    heyApiPlugin({
      config: {
        input: 'http://localhost:8000/openapi.json', // sign up at app.heyapi.dev
        output: 'src/client',
        plugins: [
        '@hey-api/sdk',
        '@tanstack/react-query', 
        'zod', 
        ]
      },
    }) as PluginOption,
    // visualizer({ open: true, filename: 'bundle-stats.html', gzipSize: true }),
  ],
  server: {
    port: 3000,
    host: '0.0.0.0',
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
});
