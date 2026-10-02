/* global __dirname */

const fs = require("fs");
const path = require("path");

const ANDROID_DENSITY_DIRS = [
  "mipmap-mdpi",
  "mipmap-hdpi",
  "mipmap-xhdpi",
  "mipmap-xxhdpi",
  "mipmap-xxxhdpi",
  "mipmap-anydpi-v26",
];

const LAUNCHER_RESOURCE_RE =
  /^ic_launcher(?:_round|_foreground|_background|_monochrome)?\.(?:png|webp|xml|avif)$/;

const WEB_ICON_FILES = [
  "favicon.ico",
  "apple-touch-icon.png",
  "icon-192.png",
  "icon-512.png",
  "icon-192-maskable.png",
  "icon-512-maskable.png",
];

function getProjectRoot(projectRoot) {
  return projectRoot ?? path.resolve(__dirname, "..");
}

function getAppIconRoot(projectRoot) {
  return path.join(getProjectRoot(projectRoot), "assets", "appIcons");
}

function assertDir(dirPath, label) {
  if (!fs.existsSync(dirPath) || !fs.statSync(dirPath).isDirectory()) {
    throw new Error(`Missing ${label}: ${dirPath}`);
  }
}

function assertFile(filePath, label) {
  if (!fs.existsSync(filePath) || !fs.statSync(filePath).isFile()) {
    throw new Error(`Missing ${label}: ${filePath}`);
  }
}

function copyFileWithDirs(sourceFile, targetFile) {
  fs.mkdirSync(path.dirname(targetFile), { recursive: true });
  fs.copyFileSync(sourceFile, targetFile);
}

function copyDirectory(sourceDir, targetDir) {
  assertDir(sourceDir, "source directory");
  fs.mkdirSync(targetDir, { recursive: true });

  for (const entry of fs.readdirSync(sourceDir, { withFileTypes: true })) {
    const sourcePath = path.join(sourceDir, entry.name);
    const targetPath = path.join(targetDir, entry.name);

    if (entry.isDirectory()) {
      copyDirectory(sourcePath, targetPath);
    } else if (entry.isFile()) {
      copyFileWithDirs(sourcePath, targetPath);
    }
  }
}

function validateAndroidIconPack(projectRoot) {
  const androidResRoot = path.join(getAppIconRoot(projectRoot), "android", "res");
  assertDir(androidResRoot, "generated Android icon res directory");

  for (const densityDir of ANDROID_DENSITY_DIRS) {
    const sourceDir = path.join(androidResRoot, densityDir);
    assertDir(sourceDir, `generated Android ${densityDir} directory`);
  }

  assertFile(
    path.join(androidResRoot, "mipmap-anydpi-v26", "ic_launcher.xml"),
    "generated Android adaptive icon XML",
  );

  return androidResRoot;
}

function getIosIconSourceDir(projectRoot) {
  const iosRoot = path.join(getAppIconRoot(projectRoot), "ios");
  const nested = path.join(iosRoot, "AppIcon.appiconset");

  if (fs.existsSync(nested)) {
    assertFile(path.join(nested, "Contents.json"), "generated iOS AppIcon Contents.json");
    return nested;
  }

  assertFile(path.join(iosRoot, "Contents.json"), "generated iOS AppIcon Contents.json");
  return iosRoot;
}

function validateIconPack(projectRoot) {
  validateAndroidIconPack(projectRoot);
  getIosIconSourceDir(projectRoot);

  const webRoot = path.join(getAppIconRoot(projectRoot), "web");
  assertDir(webRoot, "generated web icon directory");

  for (const fileName of WEB_ICON_FILES) {
    assertFile(path.join(webRoot, fileName), `generated web ${fileName}`);
  }
}

function removeLauncherResources(mipmapDir) {
  if (!fs.existsSync(mipmapDir)) {
    return;
  }

  for (const entry of fs.readdirSync(mipmapDir, { withFileTypes: true })) {
    if (entry.isFile() && LAUNCHER_RESOURCE_RE.test(entry.name)) {
      fs.unlinkSync(path.join(mipmapDir, entry.name));
    }
  }
}

function syncAndroidIcons(projectRoot, options = {}) {
  const root = getProjectRoot(projectRoot);
  const androidProjectRoot = options.androidProjectRoot ?? path.join(root, "android");

  if (!fs.existsSync(androidProjectRoot)) {
    return { skipped: true, reason: "Android native folder does not exist." };
  }

  const sourceResRoot = validateAndroidIconPack(root);
  const targetResRoot = path.join(androidProjectRoot, "app", "src", "main", "res");
  assertDir(targetResRoot, "Android native res directory");

  for (const densityDir of ANDROID_DENSITY_DIRS) {
    const sourceDir = path.join(sourceResRoot, densityDir);
    const targetDir = path.join(targetResRoot, densityDir);

    fs.mkdirSync(targetDir, { recursive: true });
    removeLauncherResources(targetDir);
    copyDirectory(sourceDir, targetDir);
  }

  const adaptiveDir = path.join(targetResRoot, "mipmap-anydpi-v26");
  const launcherXml = path.join(adaptiveDir, "ic_launcher.xml");
  const roundXml = path.join(adaptiveDir, "ic_launcher_round.xml");

  if (fs.existsSync(launcherXml)) {
    copyFileWithDirs(launcherXml, roundXml);
  }

  return { skipped: false, target: targetResRoot };
}

function findIosAppIconTarget(projectRoot, options = {}) {
  const root = getProjectRoot(projectRoot);
  const iosProjectRoot = options.iosProjectRoot ?? path.join(root, "ios");

  if (!fs.existsSync(iosProjectRoot)) {
    return undefined;
  }

  const candidates = [];

  function walk(dir) {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      const current = path.join(dir, entry.name);
      if (!entry.isDirectory()) {
        continue;
      }

      if (entry.name === "Images.xcassets") {
        candidates.push(path.join(current, "AppIcon.appiconset"));
      } else {
        walk(current);
      }
    }
  }

  walk(iosProjectRoot);
  return candidates[0];
}

function syncIosIcons(projectRoot, options = {}) {
  const sourceDir = getIosIconSourceDir(projectRoot);
  const targetDir = options.targetDir ?? findIosAppIconTarget(projectRoot, options);

  if (!targetDir) {
    return { skipped: true, reason: "iOS native AppIcon.appiconset target was not found." };
  }

  fs.rmSync(targetDir, { recursive: true, force: true });
  copyDirectory(sourceDir, targetDir);

  return { skipped: false, target: targetDir };
}

function syncWebIcons(projectRoot, options = {}) {
  const root = getProjectRoot(projectRoot);
  const sourceDir = path.join(getAppIconRoot(root), "web");
  const targetDir = options.webPublicRoot ?? path.join(root, "public");

  for (const fileName of WEB_ICON_FILES) {
    copyFileWithDirs(path.join(sourceDir, fileName), path.join(targetDir, fileName));
  }

  return { skipped: false, target: targetDir };
}

module.exports = {
  ANDROID_DENSITY_DIRS,
  WEB_ICON_FILES,
  getIosIconSourceDir,
  syncAndroidIcons,
  syncIosIcons,
  syncWebIcons,
  validateIconPack,
};
