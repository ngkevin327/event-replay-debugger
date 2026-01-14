export const a11yRules = {
  runOnly: {
    type: "tag" as const,
    values: ["wcag2a", "wcag2aa", "wcag21aa"],
  },
  rules: {
    "color-contrast": { enabled: true },
    "document-title": { enabled: true },
    region: { enabled: true },
  },
};

export const seriousImpact = ["serious", "critical"];
