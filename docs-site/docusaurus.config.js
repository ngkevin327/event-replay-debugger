// @ts-check
/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Replay Platform",
  tagline: "Production incident replay for async backends",
  favicon: "img/favicon.ico",
  url: "https://docs.replay.example",
  baseUrl: "/",
  organizationName: "replay-platform",
  projectName: "replay-platform",
  onBrokenLinks: "throw",
  i18n: { defaultLocale: "en", locales: ["en"] },
  presets: [
    [
      "@docusaurus/preset-classic",
      {
        docs: {
          editUrl:
            "https://github.com/replay-platform/replay/tree/main/docs-site/",
        },
        blog: false,
      },
    ],
  ],
  themeConfig: {
    navbar: {
      title: "Replay",
      items: [
        { type: "docSidebar", sidebarId: "docs", label: "Docs" },
        {
          href: "https://github.com/replay-platform/replay",
          label: "GitHub",
          position: "right",
        },
      ],
    },
  },
};

module.exports = config;
