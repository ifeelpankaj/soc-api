#!/usr/bin/env node

const {
  syncAndroidIcons,
  syncIosIcons,
  syncWebIcons,
  validateIconPack,
} = require("./app-icon-sync-utils");

const projectRoot = process.cwd();

function logResult(label, result) {
  if (result.skipped) {
    console.log(`${label}: skipped - ${result.reason}`);
    return;
  }

  console.log(`${label}: synced ${result.target}`);
}

try {
  validateIconPack(projectRoot);
  logResult("Android", syncAndroidIcons(projectRoot));
  logResult("iOS", syncIosIcons(projectRoot));
  logResult("Web", syncWebIcons(projectRoot));
} catch (error) {
  console.error(error instanceof Error ? error.message : error);
  process.exit(1);
}
