#!/usr/bin/env node

const fs = require("node:fs");
const path = require("node:path");
const { spawnSync } = require("node:child_process");
const readline = require("node:readline/promises");

const root = path.join(__dirname, "..");
const configRoot = path.join(root, "..", "apna-gate-config");

const ENVIRONMENTS = {
  development: { profile: "development", environment: "development" },
  preview: { profile: "preview", environment: "production" },
  production: { profile: "production", environment: "production" },
};

const PLATFORM_VALUES = ["android", "ios", "all"];
const PLATFORMS = new Set(PLATFORM_VALUES);
const REQUIRED_PUBLIC_VARS = [
  "EXPO_PUBLIC_API_BASE_URL",
  "EXPO_PUBLIC_SWAGGER_URL",
  "EXPO_PUBLIC_WEB_BASE_URL",
];
const REQUIRED_FILE_VARS = {
  android: ["GOOGLE_SERVICES_JSON"],
  ios: ["GOOGLE_SERVICE_INFO_PLIST"],
  all: ["GOOGLE_SERVICES_JSON", "GOOGLE_SERVICE_INFO_PLIST"],
};
const FILE_VAR_DEFAULTS = {
  GOOGLE_SERVICES_JSON: path.join(configRoot, "google-services.json"),
  GOOGLE_SERVICE_INFO_PLIST: path.join(configRoot, "GoogleService-Info.plist"),
};
const FILE_VAR_LABELS = {
  GOOGLE_SERVICES_JSON: "Android Firebase config",
  GOOGLE_SERVICE_INFO_PLIST: "iOS Firebase config",
};
const SUSPICIOUS_PUBLIC_NAME =
  /EXPO_PUBLIC_.*(ADMIN|PASSWORD|PRIVATE|SECRET|TOKEN|CREDENTIAL)/i;

class UsageError extends Error {}

function printHelp() {
  console.log(`Apna Gate EAS helper

Usage:
  node scripts/eas.js help
  node scripts/eas.js setup [--env development|preview|production] [--platform android|ios|all] [--pull]
  node scripts/eas.js pull [--env development|preview|production] [--path .env.local]
  node scripts/eas.js push [--env development|preview|production] [--path .env.local] [--force] [--dry-run]
  node scripts/eas.js file-var [--env development|preview|production] [--platform android|ios|all] [--force] [--dry-run]
  node scripts/eas.js build [--env development|preview|production] [--platform android|ios|all] [--dry-run]

NPM scripts:
  npm run eas:help
    Show this help page.

  npm run eas:setup -- --env production --platform android
    Validate EAS CLI/login, app config, selected EAS environment, and required platform vars.

  npm run eas:pull -- --env development
    Pull EAS environment variables into .env.local. Direction: EAS -> local.

  npm run eas:push -- --env development
    Push .env.local string variables into an EAS environment. Direction: local -> EAS.

  npm run eas:file-var -- --env production --platform android
    Upload mobile Firebase config files as EAS file variables.

  npm run eas:build -- --env preview --platform android --dry-run
    Validate production config for the preview APK profile and print the EAS build command without starting a build.

  npm run eas:build:android:dev
    Build Android using the development profile/environment.

  npm run eas:build:android:preview
    Build Android preview APK using the preview profile with production config.

  npm run eas:build:android:prod
    Build Android using the production profile/environment.

  npm run eas:build:ios:prod
    Build iOS using the production profile/environment.

  npm run eas:build:all:prod
    Build Android and iOS using the production profile/environment.

Directions:
  eas:pull      EAS -> .env.local
  eas:push      .env.local -> EAS string variables
  eas:file-var  Firebase config files -> EAS file variables
  eas:setup     validate configuration
  eas:build     validate + build

Examples:
  npm run eas:setup -- --env production --platform android
  npm run eas:pull -- --env development
  npm run eas:push -- --env development
  npm run eas:file-var -- --env production --platform android --dry-run
  npm run eas:build -- --env preview --platform android --dry-run`);
}

function parseArgs(argv) {
  const [command, ...rest] = argv;
  const options = {
    command,
    env: undefined,
    envProvided: false,
    platform: undefined,
    platformProvided: false,
    path: ".env.local",
    pathProvided: false,
    filePaths: {},
    dryRun: false,
    force: false,
    pull: false,
  };

  if (command === "--help" || command === "-h") {
    options.help = true;
    return options;
  }

  for (let index = 0; index < rest.length; index += 1) {
    const arg = rest[index];

    if (arg === "--help" || arg === "-h") {
      options.help = true;
    } else if (arg === "--dry-run") {
      options.dryRun = true;
    } else if (arg === "--force") {
      options.force = true;
    } else if (arg === "--pull") {
      options.pull = true;
    } else if (arg === "--env") {
      options.env = requireValue(rest, ++index, "--env");
      options.envProvided = true;
    } else if (arg.startsWith("--env=")) {
      options.env = arg.slice("--env=".length);
      options.envProvided = true;
    } else if (arg === "--platform") {
      options.platform = requireValue(rest, ++index, "--platform");
      options.platformProvided = true;
    } else if (arg.startsWith("--platform=")) {
      options.platform = arg.slice("--platform=".length);
      options.platformProvided = true;
    } else if (arg === "--path") {
      options.path = requireValue(rest, ++index, "--path");
      options.pathProvided = true;
    } else if (arg.startsWith("--path=")) {
      options.path = arg.slice("--path=".length);
      options.pathProvided = true;
    } else if (arg === "--android-file") {
      options.filePaths.GOOGLE_SERVICES_JSON = requireValue(
        rest,
        ++index,
        "--android-file",
      );
    } else if (arg.startsWith("--android-file=")) {
      options.filePaths.GOOGLE_SERVICES_JSON = arg.slice("--android-file=".length);
    } else if (arg === "--ios-file") {
      options.filePaths.GOOGLE_SERVICE_INFO_PLIST = requireValue(
        rest,
        ++index,
        "--ios-file",
      );
    } else if (arg.startsWith("--ios-file=")) {
      options.filePaths.GOOGLE_SERVICE_INFO_PLIST = arg.slice("--ios-file=".length);
    } else {
      throw new UsageError(`Unknown argument: ${arg}`);
    }
  }

  return options;
}

function requireValue(args, index, flag) {
  const value = args[index];

  if (!value || value.startsWith("--")) {
    throw new UsageError(`Missing value for ${flag}.`);
  }

  return value;
}

function easExecutable() {
  return process.platform === "win32" ? "eas.cmd" : "eas";
}

function runEas(args, options = {}) {
  if (process.platform === "win32") {
    return runEasOnWindows(args, options);
  }

  const result = spawnSync(easExecutable(), args, {
    cwd: root,
    encoding: "utf8",
    stdio: options.stdio ?? "pipe",
  });

  if (result.error) {
    throw result.error;
  }

  return result;
}

function runEasOnWindows(args, options = {}) {
  const command = [easExecutable(), ...args.map(formatWindowsShellArg)].join(" ");
  const result = spawnSync("cmd.exe", ["/d", "/s", "/c", command], {
    cwd: root,
    encoding: "utf8",
    stdio: options.stdio ?? "pipe",
  });

  if (result.error) {
    throw result.error;
  }

  return result;
}

function formatCommand(args) {
  return `eas ${args.map(formatArg).join(" ")}`;
}

function formatArg(arg) {
  return /\s/.test(arg) ? `"${arg.replace(/"/g, '\\"')}"` : arg;
}

function formatWindowsShellArg(arg) {
  if (/^[A-Za-z0-9_./:-]+$/.test(arg)) {
    return arg;
  }

  return `"${arg.replace(/"/g, '\\"')}"`;
}

function canPrompt() {
  return Boolean(process.stdin.isTTY && process.stdout.isTTY);
}

function createPrompter() {
  return readline.createInterface({
    input: process.stdin,
    output: process.stdout,
  });
}

async function promptChoice(rl, message, choices, defaultValue) {
  const defaultIndex = Math.max(choices.indexOf(defaultValue), 0);
  console.log(message);
  choices.forEach((choice, index) => {
    const marker = index === defaultIndex ? "*" : " ";
    console.log(`  ${index + 1}. ${choice} ${marker}`);
  });

  const answer = await rl.question(`Select [${defaultIndex + 1}]: `);
  const selected = answer.trim() ? Number(answer.trim()) : defaultIndex + 1;

  if (!Number.isInteger(selected) || selected < 1 || selected > choices.length) {
    throw new UsageError(`Invalid selection for ${message}`);
  }

  return choices[selected - 1];
}

async function promptText(rl, message, defaultValue) {
  const answer = await rl.question(`${message} [${defaultValue}]: `);
  return answer.trim() || defaultValue;
}

async function confirm(rl, message, defaultYes = false) {
  const suffix = defaultYes ? "(Y/n)" : "(y/N)";
  const answer = await rl.question(`${message} ${suffix} `);
  const normalized = answer.trim().toLowerCase();

  if (!normalized) {
    return defaultYes;
  }

  return normalized === "y" || normalized === "yes";
}

function requireKnownOptions(options) {
  if (
    !options.command ||
    !["help", "setup", "pull", "push", "file-var", "build"].includes(
      options.command,
    )
  ) {
    throw new UsageError(
      "Command must be one of: help, setup, pull, push, file-var, build.",
    );
  }

  if (options.env && !ENVIRONMENTS[options.env]) {
    throw new UsageError(
      `Unknown environment "${options.env}". Use development, preview, or production.`,
    );
  }

  if (options.platform && !PLATFORMS.has(options.platform)) {
    throw new UsageError(
      `Unknown platform "${options.platform}". Use android, ios, or all.`,
    );
  }
}

async function prepareOptions(options) {
  const needsEnv = ["setup", "pull", "push", "file-var", "build"].includes(
    options.command,
  );
  const needsPlatform = ["setup", "file-var", "build"].includes(options.command);
  const needsPath = ["pull", "push"].includes(options.command);
  const needsPrompt =
    (needsEnv && !options.env) ||
    (needsPlatform && !options.platform) ||
    (needsPath && !options.pathProvided && options.command === "push") ||
    needsFileVarPathPrompt(options);

  if (!needsPrompt) {
    enforceRequiredOptions(options);
    return options;
  }

  if (!canPrompt()) {
    enforceRequiredOptions(options);
    return options;
  }

  const rl = createPrompter();

  try {
    if (needsEnv && !options.env) {
      options.env = await promptChoice(
        rl,
        "Select EAS environment:",
        Object.keys(ENVIRONMENTS),
        "development",
      );
    }

    if (needsPlatform && !options.platform) {
      options.platform = await promptChoice(
        rl,
        "Select platform:",
        PLATFORM_VALUES,
        "android",
      );
    }

    if (options.command === "push" && !options.pathProvided) {
      options.path = await promptText(rl, "Env file path", options.path);
    }

    if (options.command === "file-var") {
      for (const variable of getRequiredFileVars(options.platform)) {
        if (!options.filePaths[variable]) {
          const relativeDefault = path.relative(root, FILE_VAR_DEFAULTS[variable]);
          options.filePaths[variable] = await promptText(
            rl,
            `${FILE_VAR_LABELS[variable]} path`,
            relativeDefault,
          );
        }
      }
    }
  } finally {
    rl.close();
  }

  enforceRequiredOptions(options);
  return options;
}

function needsFileVarPathPrompt(options) {
  if (options.command !== "file-var" || !options.platform || !canPrompt()) {
    return false;
  }

  return getRequiredFileVars(options.platform).some(
    (variable) => !options.filePaths[variable],
  );
}

function enforceRequiredOptions(options) {
  if (!options.env) {
    throw new UsageError(
      "Missing --env. Use --env development, --env preview, or --env production.",
    );
  }

  if (["setup", "file-var", "build"].includes(options.command) && !options.platform) {
    throw new UsageError("Missing --platform. Use --platform android, ios, or all.");
  }
}

function validateProjectFiles(env) {
  const easJsonPath = path.join(root, "eas.json");
  const appConfigPath = path.join(root, "app.config.js");

  if (!fs.existsSync(easJsonPath)) {
    throw new Error("Missing eas.json in the app root.");
  }

  if (!fs.existsSync(appConfigPath)) {
    throw new Error("Missing app.config.js in the app root.");
  }

  const easJson = JSON.parse(fs.readFileSync(easJsonPath, "utf8"));
  const target = ENVIRONMENTS[env];
  const profile = target.profile;
  const environment = target.environment;
  const buildProfile = easJson.build?.[profile];

  if (!buildProfile) {
    throw new Error(`Missing eas.json build profile "${profile}".`);
  }

  if (buildProfile.environment !== environment) {
    throw new Error(
      `eas.json profile "${profile}" must set "environment": "${environment}".`,
    );
  }

  return { environment, profile };
}

function checkEasCli() {
  const version = runEas(["--version"]);

  if (version.status !== 0) {
    throw new Error(
      "EAS CLI is not available. Install it with: npm install -g eas-cli",
    );
  }
}

function checkEasLogin() {
  const whoami = runEas(["whoami"]);

  if (whoami.status !== 0) {
    throw new Error("You are not logged in to EAS. Run: eas login");
  }
}

function listEasEnvNames(env) {
  const result = runEas(["env:list", "--environment", env, "--format", "short"]);

  if (result.status !== 0) {
    throw new Error(
      [
        `Could not list EAS environment variables for "${env}".`,
        result.stderr?.trim(),
      ]
        .filter(Boolean)
        .join("\n"),
    );
  }

  const output = `${result.stdout ?? ""}\n${result.stderr ?? ""}`;
  const names = new Set();

  for (const name of [...REQUIRED_PUBLIC_VARS, ...REQUIRED_FILE_VARS.all]) {
    const pattern = new RegExp(`(^|[^A-Z0-9_])${name}([^A-Z0-9_]|$)`);
    if (pattern.test(output)) {
      names.add(name);
    }
  }

  return names;
}

function getRequiredVars(platform) {
  return [...REQUIRED_PUBLIC_VARS, ...REQUIRED_FILE_VARS[platform]];
}

function getRequiredFileVars(platform) {
  return REQUIRED_FILE_VARS[platform];
}

function validateEasEnvironment(env, platform) {
  const names = listEasEnvNames(env);
  const required = getRequiredVars(platform);
  const missing = required.filter((name) => !names.has(name));

  if (missing.length > 0) {
    const publicMissing = missing.filter((name) => name.startsWith("EXPO_PUBLIC_"));
    const fileMissing = missing.filter((name) => !name.startsWith("EXPO_PUBLIC_"));
    const guidance = [];

    if (publicMissing.length > 0) {
      guidance.push(
        "Create missing public vars with:",
        ...publicMissing.map(
          (name) =>
            `  eas env:set --environment ${env} --name ${name} --value <value> --visibility plaintext`,
        ),
      );
    }

    if (fileMissing.length > 0) {
      guidance.push(
        "Upload missing Firebase file vars with:",
        ...fileMissing.map(
          (name) =>
            `  npm run eas:file-var -- --env ${env} --platform ${platformForFileVar(name)}`,
        ),
      );
    }

    throw new Error(
      [`Missing EAS env vars for ${env}/${platform}: ${missing.join(", ")}`, ...guidance].join(
        "\n",
      ),
    );
  }
}

function platformForFileVar(name) {
  return name === "GOOGLE_SERVICES_JSON" ? "android" : "ios";
}

function validate(options) {
  checkEasCli();
  checkEasLogin();
  const { environment, profile } = validateProjectFiles(options.env);
  validateEasEnvironment(environment, options.platform);
  return { environment, profile };
}

function resolveProjectPath(filePath) {
  return path.isAbsolute(filePath) ? filePath : path.resolve(root, filePath);
}

function validateEnvFileForPush(envPath) {
  const absolutePath = resolveProjectPath(envPath);

  if (!fs.existsSync(absolutePath)) {
    throw new Error(`Env file does not exist: ${envPath}`);
  }

  const contents = fs.readFileSync(absolutePath, "utf8");
  const suspicious = [];
  const fileVars = [];

  for (const line of contents.split(/\r?\n/)) {
    const trimmed = line.trim();

    if (!trimmed || trimmed.startsWith("#") || !trimmed.includes("=")) {
      continue;
    }

    const name = trimmed.split("=", 1)[0].trim();

    if (SUSPICIOUS_PUBLIC_NAME.test(name)) {
      suspicious.push(name);
    }

    if (REQUIRED_FILE_VARS.all.includes(name)) {
      fileVars.push(name);
    }
  }

  if (suspicious.length > 0) {
    throw new Error(
      [
        `Refusing to push suspicious public variable names: ${suspicious.join(", ")}`,
        "EXPO_PUBLIC_* values are bundled into the mobile app and readable by users.",
      ].join("\n"),
    );
  }

  if (fileVars.length > 0) {
    throw new Error(
      [
        `Do not put Firebase file variables in ${envPath}: ${fileVars.join(", ")}`,
        "Use npm run eas:file-var instead.",
      ].join("\n"),
    );
  }
}

async function confirmProduction(options, summaryLines) {
  const environment = ENVIRONMENTS[options.env]?.environment ?? options.env;

  if (environment !== "production" || options.force || options.dryRun) {
    return;
  }

  if (!canPrompt()) {
    throw new UsageError("Production changes require --force in non-interactive mode.");
  }

  const rl = createPrompter();

  try {
    console.log("");
    console.log("Target environment: production");
    for (const line of summaryLines) {
      console.log(line);
    }

    const ok = await confirm(rl, "Continue?", false);

    if (!ok) {
      throw new Error("Cancelled.");
    }
  } finally {
    rl.close();
  }
}

function runPull(options) {
  checkEasCli();
  checkEasLogin();
  const environment = ENVIRONMENTS[options.env].environment;

  const args = [
    "env:pull",
    "--environment",
    environment,
    "--path",
    options.path,
  ];

  console.log(`Pulling EAS ${environment} env to ${options.path}`);
  const result = runEas(args, { stdio: "inherit" });

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

async function runPush(options) {
  const environment = ENVIRONMENTS[options.env].environment;
  const args = [
    "env:push",
    "--environment",
    environment,
    "--path",
    options.path,
  ];

  if (options.force) {
    args.push("--force");
  }

  validateEnvFileForPush(options.path);

  console.log(`Environment : ${environment}`);
  console.log(`Source      : ${options.path}`);

  if (options.dryRun) {
    console.log("\nWould execute:");
    console.log(formatCommand(args));
    return;
  }

  await confirmProduction(options, [`Source: ${options.path}`]);
  checkEasCli();
  checkEasLogin();

  const result = runEas(args, { stdio: "inherit" });

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

async function runFileVar(options) {
  const environment = ENVIRONMENTS[options.env].environment;
  const uploads = getRequiredFileVars(options.platform).map((name) => ({
    name,
    source: options.filePaths[name] ?? path.relative(root, FILE_VAR_DEFAULTS[name]),
  }));

  for (const upload of uploads) {
    const absolutePath = resolveProjectPath(upload.source);

    if (!fs.existsSync(absolutePath)) {
      throw new Error(`${FILE_VAR_LABELS[upload.name]} does not exist: ${upload.source}`);
    }
  }

  for (const upload of uploads) {
    const args = [
      "env:set",
      "--environment",
      environment,
      "--name",
      upload.name,
      "--value",
      upload.source,
      "--type",
      "file",
      "--visibility",
      "sensitive",
    ];

    console.log(`Environment : ${environment}`);
    console.log(`Platform    : ${options.platform}`);
    console.log(`Variable    : ${upload.name}`);
    console.log(`Source      : ${upload.source}`);

    if (options.dryRun) {
      console.log("\nWould execute:");
      console.log(formatCommand(args));
      console.log("");
      continue;
    }

    await confirmProduction(options, [
      `Variable: ${upload.name}`,
      `Source: ${upload.source}`,
    ]);

    checkEasCli();
    checkEasLogin();

    const result = runEas(args, { stdio: "inherit" });

    if (result.status !== 0) {
      process.exit(result.status ?? 1);
    }

    console.log("");
  }
}

function runSetup(options) {
  const { environment, profile } = validate(options);

  console.log("EAS setup validation passed.");
  console.log(`Environment : ${environment}`);
  console.log(`Platform    : ${options.platform}`);
  console.log(`Profile     : ${profile}`);

  if (options.pull) {
    runPull(options);
  }
}

function runBuild(options) {
  const { environment, profile } = validate(options);
  const buildArgs = ["build", "--platform", options.platform, "--profile", profile];

  console.log("EAS setup validation passed.");
  console.log(`Environment : ${environment}`);
  console.log(`Platform    : ${options.platform}`);
  console.log(`Profile     : ${profile}`);

  if (options.dryRun) {
    console.log("\nWould execute:");
    console.log(formatCommand(buildArgs));
    return;
  }

  const result = runEas(buildArgs, { stdio: "inherit" });

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}

async function main() {
  let options;

  try {
    options = parseArgs(process.argv.slice(2));

    if (options.help || options.command === "help") {
      printHelp();
      return;
    }

    requireKnownOptions(options);
    await prepareOptions(options);

    if (options.command === "setup") {
      runSetup(options);
    } else if (options.command === "pull") {
      runPull(options);
    } else if (options.command === "push") {
      await runPush(options);
    } else if (options.command === "file-var") {
      await runFileVar(options);
    } else {
      runBuild(options);
    }
  } catch (error) {
    console.error(error.message);

    if (error instanceof UsageError) {
      console.error("");
      printHelp();
    }

    process.exit(1);
  }
}

main();
