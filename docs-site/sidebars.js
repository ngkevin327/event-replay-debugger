/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  docs: [
    "intro",
    "quickstart",
    "agent-config",
    "api-reference",
    {
      type: "category",
      label: "Playbooks",
      items: ["playbooks/fintech-duplicate-payout"],
    },
  ],
};

module.exports = sidebars;
