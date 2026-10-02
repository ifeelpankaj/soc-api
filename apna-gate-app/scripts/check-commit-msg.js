#!/usr/bin/env node

const fs = require("node:fs");

const [, , commitMsgPath] = process.argv;

if (!commitMsgPath) {
  console.error("Missing commit message file path.");
  process.exit(1);
}

const message = fs.readFileSync(commitMsgPath, "utf8").trim();
const firstLine = message.split(/\r?\n/, 1)[0] ?? "";
const allowedTypes = [
  "build",
  "chore",
  "ci",
  "docs",
  "feat",
  "fix",
  "perf",
  "refactor",
  "revert",
  "style",
  "test",
];
const typePattern = allowedTypes.join("|");
const conventionalPattern = new RegExp(
  `^(${typePattern})(\\([a-z0-9-]+\\))?!?: .{1,100}$`,
);

if (
  conventionalPattern.test(firstLine) ||
  firstLine.startsWith("Merge ") ||
  firstLine.startsWith("Revert ") ||
  firstLine.startsWith("fixup!") ||
  firstLine.startsWith("squash!")
) {
  process.exit(0);
}

console.error("Invalid commit message.");
console.error("");
console.error("Use Conventional Commit format:");
console.error("  type(scope): short summary");
console.error("");
console.error(`Allowed types: ${allowedTypes.join(", ")}`);
console.error("");
console.error("Examples:");
console.error("  feat(app): add login screen");
console.error("  fix(app): handle expired session");
console.error("  docs: update EAS workflow");
console.error("");
console.error(`Received: ${firstLine || "<empty>"}`);
process.exit(1);
