import { defineConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';

export default defineConfig({
  plugins: [pluginReact()],
  tools: {
    postcss: (config) => {
      // Ensure the plugins array exists
      config.plugins = config.plugins || [];
      config.plugins.push(require('tailwindcss'));
    },
  },
});
