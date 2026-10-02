/* global process, require, console */
const { load } = require("@expo/env");
process.env.NODE_ENV ??= "production";
load(process.cwd());
const config = require("../app.config.js").expo;
const problems = [];
for (const key of ["EXPO_PUBLIC_API_BASE_URL", "EXPO_PUBLIC_WEB_BASE_URL"]) {
  try {
    const url = new URL(process.env[key]);
    const host = url.hostname.toLowerCase();
    if (
      url.protocol !== "https:" ||
      url.username ||
      url.password ||
      host === "localhost" ||
      host.endsWith(".local") ||
      host.endsWith(".localhost") ||
      host.includes(":") ||
      /^\d+\.\d+\.\d+\.\d+$/.test(host)
    ) {
      problems.push(
        `${key} must use a public HTTPS hostname without embedded credentials.`,
      );
    }
  } catch {
    problems.push(`${key} must contain a valid production URL.`);
  }
}
if (config.userInterfaceStyle !== "light")
  problems.push("The release appearance must match the light theme.");
if (!config.android?.package || !config.ios?.bundleIdentifier)
  problems.push("Both platform identifiers are required.");
if (problems.length) {
  for (const problem of problems) console.error(problem);
  process.exitCode = 1;
} else console.log("Release URL, appearance and app identifier checks passed.");
