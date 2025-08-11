import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import UnoCSS from "unocss/vite";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));

// https://vite.dev/config/
export default defineConfig({
  plugins: [UnoCSS(), react()],
  resolve: {
    alias: { "@/": resolve(__dirname, "src") + "/" },
  },
});
