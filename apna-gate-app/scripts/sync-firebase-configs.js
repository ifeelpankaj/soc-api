const fs = require("node:fs");
const path = require("node:path");

const root = path.join(__dirname, "..");
const configsDir = path.join(root, "..", "configs");
const targets = [
  {
    source: path.join(configsDir, "google-services.json"),
    destination: path.join(root, "google-services.json"),
  },
  {
    source: path.join(configsDir, "GoogleService-Info.plist"),
    destination: path.join(root, "GoogleService-Info.plist"),
  },
];

for (const { source, destination } of targets) {
  fs.copyFileSync(source, destination);
  console.log(`Synced ${source} -> ${destination}`);
}

console.log("Firebase config sync complete.");
