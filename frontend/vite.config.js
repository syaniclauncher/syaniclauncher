import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { resolve } from "node:path";

export default defineConfig({
    plugins: [vue()],
    base: "./",
    resolve: {
        alias: {
            "@": resolve(__dirname, "src"),
            "@wails": resolve(__dirname, "wailsjs"),
        },
    },
    build: { outDir: "dist", emptyOutDir: true },
});
