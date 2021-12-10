const { tailwindExtractor } = require('tailwindcss/lib/lib/purgeUnusedStyles');

const config = {
  mode: 'aot', // jit
  purge: {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    options: {
      defaultExtractor: (content) => [
        // If this stops working, please open an issue at https://github.com/svelte-add/svelte-add/issues rather than bothering Tailwind Labs about it
        ...tailwindExtractor(content),
        // Match Svelte class: directives (https://github.com/tailwindlabs/tailwindcss/discussions/1731)
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        ...[...content.matchAll(/(?:class:)*([\w\d-/:%.]+)/gm)].map(([_match, group, ..._rest]) => group),
      ],
    },
    safelist: [/^svelte-[\d\w]+$/],
  },
  corePlugins: {
    container: false,
  },
  theme: {
    extend: {},
  },
  plugins: [],
};

module.exports = config;
