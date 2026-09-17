/**
 * @file Oxlint configuration
 */

const { default: createOxlintConfig } = await import("@goauthentik/oxlint-config");

export default createOxlintConfig({ lit: true, react: true, padding: true });
