import react from "@vitejs/plugin-react-swc";
import path from "path";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig, type PluginOption } from "vite";
import { heyApiPlugin } from "@hey-api/vite-plugin";

// https://vite.dev/config/
export default defineConfig({
  // base: '/admin',
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          "vendor-react": ["react", "react-dom"],
          "vendor-utils": ["@tanstack/react-query", "zod"],
        },
      },
    },
  },
  plugins: [
    react(),
    tailwindcss(),
    heyApiPlugin({
      config: {
        input: "./openapi.json", // sign up at app.heyapi.dev
        output: "src/client",
        plugins: ["@hey-api/sdk", "@tanstack/react-query", "zod"],
      },
    }) as PluginOption,
  ],
  server: {
    port: 3000,
    host: "0.0.0.0",
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
