/* global __dirname */

const path = require("path");
// Default paths for Google services configuration files
const defaultGoogleServicesJson = path.join(
  __dirname,
  "..",
  "apna-gate-config",
  "google-services.json",
);
const defaultGoogleServiceInfoPlist = path.join(
  __dirname,
  "..",
  "apna-gate-config",
  "GoogleService-Info.plist",
);

module.exports = {
  expo: {
    name: "Apna Gate",
    slug: "apna-gate",
    version: "1.0.0",
    orientation: "portrait",
    icon: "./assets/appIcons/ios/AppIcon~ios-marketing.png",
    scheme: "app",
    userInterfaceStyle: "light",
    ios: {
      icon: "./assets/appIcons/ios/AppIcon~ios-marketing.png",
      bundleIdentifier: "com.apnagate",
      googleServicesFile:
        process.env.GOOGLE_SERVICE_INFO_PLIST ?? defaultGoogleServiceInfoPlist,
      supportsTablet: true,
      config: {
        usesNonExemptEncryption: false,
      },
      infoPlist: {
        NSCameraUsageDescription:
          "This app uses the camera to take visitor photos and scan QR codes.",
        NSPhotoLibraryUsageDescription:
          "This app uses the photo library to choose visitor photos and QR codes.",
      },
    },
    android: {
      package: "com.apnagate",
      googleServicesFile:
        process.env.GOOGLE_SERVICES_JSON ?? defaultGoogleServicesJson,
      softwareKeyboardLayoutMode: "pan",
      predictiveBackGestureEnabled: false,
      permissions: [
        "android.permission.CAMERA",
        "android.permission.POST_NOTIFICATIONS",
        "android.permission.VIBRATE",
      ],
    },
    web: {
      output: "static",
      bundler: "metro",
      favicon: "./assets/appIcons/web/favicon.ico",
    },
    plugins: [
      "expo-router",
      "expo-image",
      "expo-web-browser",
      [
        "expo-dev-client",
        {
          launchMode: "most-recent",
        },
      ],
      [
        "expo-splash-screen",
        {
          backgroundColor: "#F8F6F2",
          image: "./assets/images/public/splash-icon.png",
          imageWidth: 160,
          android: {
            image: "./assets/images/public/splash-icon.png",
            imageWidth: 160,
          },
          ios: {
            image: "./assets/images/public/splash-icon.png",
            imageWidth: 160,
          },
        },
      ],
      [
        "expo-secure-store",
        {
          configureAndroidBackup: true,
          faceIDPermission:
            "Allow $(PRODUCT_NAME) to access Face ID for secure authentication.",
        },
      ],
      [
        "expo-notifications",
        {
          color: "#ff6a1a",
          defaultChannel: "default",
        },
      ],
      [
        "expo-camera",
        {
          cameraPermission:
            "Allow $(PRODUCT_NAME) to use the camera for QR scanning and visitor photos.",
          recordAudioAndroid: false,
          barcodeScannerEnabled: true,
        },
      ],
      [
        "expo-image-picker",
        {
          photosPermission:
            "Allow $(PRODUCT_NAME) to choose visitor photos and QR codes from your gallery.",
          cameraPermission: "Allow $(PRODUCT_NAME) to take visitor photos and scan QR codes.",
          microphonePermission: false,
        },
      ],
      "./plugins/with-generated-app-icons",
    ],
    experiments: {
      typedRoutes: true,
      reactCompiler: true,
    },
    extra: {
      router: {},
      eas: {
        projectId: "dfcebf37-fee6-451b-af9f-ef6f7f4ad196",
      },
    },
    owner: "apna-gate",
  },
};
