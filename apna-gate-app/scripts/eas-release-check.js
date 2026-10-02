/* global process, require */
if (["preview", "production"].includes(process.env.EAS_BUILD_PROFILE)) {
  require("./check-release-config");
}
