# Apna Gate App

Mobile application for Apna Gate, built with Expo, React Native, Expo Router, Redux Toolkit, and EAS Build.

This app supports resident, guard, visitor, and profile workflows, and is configured for Android and iOS production builds through EAS environments.

## Tech Stack

- Expo `~57`
- React Native `0.86`
- React `19`
- Expo Router
- Redux Toolkit / RTK Query
- TypeScript
- EAS Build
- Firebase mobile config files for Android/iOS

## Project Structure

```text
apna-gate-app/
|-- app.config.js
|-- eas.json
|-- package.json
|-- scripts/
|   |-- eas.js
|   |-- help.js
|   `-- sync-firebase-configs.js
|-- src/
|-- assets/
`-- README.md
```

Related local config files live outside the app:

```text
apna-gate-config/
|-- google-services.json
|-- GoogleService-Info.plist
`-- smartsociety-...firebase-adminsdk...json
```

The Firebase Admin SDK JSON is server-side only. Do not upload it for mobile EAS builds.

## Requirements

- Node.js compatible with the current Expo SDK
- npm
- Expo CLI through project scripts
- EAS CLI for cloud builds
- Android Studio for local Android builds
- Xcode/macOS for local iOS builds

Install EAS CLI when working with cloud builds:

```powershell
npm install -g eas-cli
eas login
```

## Setup

Install dependencies:

```powershell
npm install
```

Create local environment config:

```powershell
copy .env.example .env.local
```

Or pull the development environment from EAS:

```powershell
npm run eas:pull -- --env development
```

For local native Firebase fallback files, keep these available:

```text
../apna-gate-config/google-services.json
../apna-gate-config/GoogleService-Info.plist
```

## Environment Variables

Required public variables:

```text
EXPO_PUBLIC_API_BASE_URL
EXPO_PUBLIC_SWAGGER_URL
EXPO_PUBLIC_WEB_BASE_URL
```

Optional:

```text
EXPO_PUBLIC_GOOGLE_WEB_CLIENT_ID
```

These `EXPO_PUBLIC_*` values are bundled into the mobile app and are readable by users. Do not put secrets, private keys, passwords, admin tokens, or database credentials in them.

## Firebase Files

Android builds use:

```text
GOOGLE_SERVICES_JSON
```

iOS builds use:

```text
GOOGLE_SERVICE_INFO_PLIST
```

For EAS cloud builds, upload these as EAS file variables:

```powershell
npm run eas:file-var -- --env production --platform android
npm run eas:file-var -- --env production --platform ios
```

Preview first without uploading:

```powershell
npm run eas:file-var -- --env production --platform all --dry-run
```

## Development

Start the Expo development server:

```powershell
npm run start
```

Start with development client:

```powershell
npm run dev
```

Run Android locally:

```powershell
npm run android
```

Run iOS locally:

```powershell
npm run ios
```

Run web:

```powershell
npm run web
```

List all available npm scripts:

```powershell
npm run help
```

## API Codegen

The app has scripts for fetching and generating API client code from Swagger/OpenAPI.

Run the complete codegen flow:

```powershell
npm run api:codegen
```

Individual steps:

```powershell
npm run api:fetch
npm run api:convert
npm run api:generate
```

Make sure `EXPO_PUBLIC_SWAGGER_URL` points to the intended API docs endpoint before regenerating.

## Quality Checks

Run lint:

```powershell
npm run lint
```

Run unit tests:

```powershell
npm run test:unit
```

Note: lint may report existing React Compiler / hooks issues in app source. Fix those before treating a production release as ready.

## EAS Workflow

This project wraps EAS commands through:

```text
scripts/eas.js
```

Detailed EAS help:

```powershell
npm run eas:help
```

Command directions:

```text
eas:pull      EAS -> .env.local
eas:push      .env.local -> EAS string variables
eas:file-var  Firebase config files -> EAS file variables
eas:setup     validate configuration
eas:build     validate + build
```

Pull development env locally:

```powershell
npm run eas:pull -- --env development
```

Push local env values to EAS:

```powershell
npm run eas:push -- --env development
```

Push production env values with confirmation:

```powershell
npm run eas:push -- --env production
```

Use `--force` for CI/non-interactive production pushes:

```powershell
npm run eas:push -- --env production --force
```

Validate production Android setup:

```powershell
npm run eas:setup -- --env production --platform android
```

Validate production iOS setup:

```powershell
npm run eas:setup -- --env production --platform ios
```

Preview a production-config APK build command without starting a cloud build:

```powershell
npm run eas:build -- --env preview --platform android --dry-run
```

## Production Builds

Android production:

```powershell
npm run eas:build:android:prod
```

iOS production:

```powershell
npm run eas:build:ios:prod
```

Android and iOS production together:

```powershell
npm run eas:build:all:prod
```

Before starting production builds, verify:

- `npm run lint` has no release-blocking errors.
- `npm run test:unit` passes.
- EAS production string variables exist.
- EAS production file variables exist for the selected platform.
- Firebase Admin SDK JSON has not been uploaded to the mobile app environment.

## EAS Profiles

`eas.json` maps profiles to explicit EAS environments:

```text
development -> development
preview     -> production
production  -> production
```

Android `development` and `preview` profiles build APKs for internal distribution. The `preview` profile is an APK preview path that loads production config.

## Troubleshooting

If EAS commands fail because you are not logged in:

```powershell
eas login
```

If local Expo behavior is stale:

```powershell
npx expo start --clear
```

For deeper Expo, Metro, and Android cache notes, see:

```text
../apna-gate-docs/clear-expo-cache.md
```

For the full EAS workflow guide, see:

```text
../apna-gate-docs/eas-build-workflow.md
```
