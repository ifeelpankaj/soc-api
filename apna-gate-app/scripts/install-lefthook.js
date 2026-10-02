#!/usr/bin/env node

const { spawnSync } = require("node:child_process");
const fs = require("node:fs");
const path = require("node:path");

const root = path.join(__dirname, "..");

function run(command, args) {
  return spawnSync(command, args, {
    cwd: root,
    encoding: "utf8",
    stdio: "pipe",
  });
}

const gitCheck = run("git", ["rev-parse", "--show-toplevel"]);

if (gitCheck.status !== 0) {
  console.log("Skipping Lefthook install: this folder is not inside a Git repository.");
  process.exit(0);
}

const lefthook = path.join(root, "node_modules", "lefthook", "bin", "index.js");

if (!fs.existsSync(lefthook)) {
  console.error("Local Lefthook binary is missing. Run npm install first.");
  process.exit(1);
}

const install = spawnSync(process.execPath, [lefthook, "install"], {
  cwd: root,
  encoding: "utf8",
  stdio: "inherit",
});

if (install.error) {
  console.error(install.error.message);
  process.exit(1);
}

process.exit(install.status ?? 0);
