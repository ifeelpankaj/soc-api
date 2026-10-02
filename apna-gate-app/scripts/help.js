#!/usr/bin/env node

const fs = require("node:fs");
const path = require("node:path");

const root = path.join(__dirname, "..");
const packageJson = JSON.parse(
  fs.readFileSync(path.join(root, "package.json"), "utf8"),
);

const descriptions = {
  help: "List all npm scripts and what they do.",
  prepare: "Install Lefthook Git hooks after npm install when Git is available.",
  start: "Start the Expo development server.",
  prebuild: "Generate native iOS/Android project files from Expo config.",
  phone: "Build and run the Android app on a connected device/emulator.",
  dev: "Start Expo for a development-client build.",
  android: "Run the Android app locally with Expo.",
  ios: "Run the iOS app locally with Expo.",
  web: "Start the app in web mode.",
  lint: "Run Expo ESLint checks.",
  "test:unit": "Run unit tests.",
  "hooks:install": "Install Lefthook Git hooks when Git is available.",
  "hooks:validate": "Validate lefthook.yml.",
  "api:fetch": "Fetch the Swagger/OpenAPI schema.",
  "api:convert": "Convert the fetched Swagger/OpenAPI schema.",
  "api:generate": "Generate API client code.",
  "api:codegen": "Fetch, convert, and generate API client code.",
  "setup:firebase-local": "Copy local Firebase config files into the app root.",
  "eas:help": "Show detailed EAS wrapper help.",
  "eas:setup": "Validate EAS login, project config, env vars, and file vars.",
  "eas:pull": "Pull EAS environment variables into .env.local.",
  "eas:push": "Push .env.local string variables into an EAS environment.",
  "eas:file-var": "Upload Firebase config files as EAS file variables.",
  "eas:build": "Validate and run an EAS build for a selected env/platform.",
  "eas:build:android:dev": "Build Android using the development EAS profile.",
  "eas:build:android:preview": "Build Android preview APK using production config.",
  "eas:build:android:prod": "Build Android using the production EAS profile.",
  "eas:build:ios:prod": "Build iOS using the production EAS profile.",
  "eas:build:all:prod": "Build Android and iOS using the production EAS profile.",
};

const scripts = packageJson.scripts ?? {};
const names = Object.keys(scripts);
const width = Math.max(...names.map((name) => name.length));

console.log("Apna Gate app npm scripts");
console.log("");

for (const name of names) {
  const description = descriptions[name] ?? "No description available.";
  console.log(`  npm run ${name.padEnd(width)}  ${description}`);
}

console.log("");
console.log("For EAS-specific details, run:");
console.log("  npm run eas:help");
