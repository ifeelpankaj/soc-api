const { withDangerousMod } = require("@expo/config-plugins");
const {
  syncAndroidIcons,
  syncIosIcons,
  validateIconPack,
} = require("../scripts/app-icon-sync-utils");

function withGeneratedAppIcons(config) {
  config = withDangerousMod(config, [
    "android",
    async (config) => {
      validateIconPack(config.modRequest.projectRoot);
      syncAndroidIcons(config.modRequest.projectRoot, {
        androidProjectRoot: config.modRequest.platformProjectRoot,
      });
      return config;
    },
  ]);

  config = withDangerousMod(config, [
    "ios",
    async (config) => {
      validateIconPack(config.modRequest.projectRoot);
      syncIosIcons(config.modRequest.projectRoot, {
        iosProjectRoot: config.modRequest.platformProjectRoot,
      });
      return config;
    },
  ]);

  return config;
}

module.exports = withGeneratedAppIcons;
