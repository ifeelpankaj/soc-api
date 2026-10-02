import type { ConfigFile } from "@rtk-query/codegen-openapi";

const config: ConfigFile = {
  schemaFile: "./src/openapi/openapi.json",
  apiFile: "./src/redux/queries/baseApi.ts",
  apiImport: "baseApi",
  outputFile: "./src/lib/api/generated-api.ts",
  exportName: "generatedApi",
  hooks: true,
  tag: true,
};

export default config;
