import { mkdir, readFile, writeFile } from "node:fs/promises";
import path from "node:path";

import "./load-env";
import { appConfig } from "../lib/config";

const swaggerOutputPath = path.resolve("src/openapi/swagger.json");
const fallbackSwaggerInputPath = path.resolve("../apna-gate-api/docs/swagger.json");

type SwaggerDocument = {
  swagger?: string;
  paths?: Record<string, unknown>;
};

function validateSwaggerDocument(
  document: unknown,
): asserts document is SwaggerDocument {
  if (!document || typeof document !== "object") {
    throw new Error("Swagger source must be a JSON object.");
  }

  const swaggerDocument = document as SwaggerDocument;
  if (swaggerDocument.swagger !== "2.0") {
    throw new Error('Swagger source must include swagger: "2.0".');
  }

  if (
    !swaggerDocument.paths ||
    Object.keys(swaggerDocument.paths).length === 0
  ) {
    throw new Error("Swagger source must include non-empty paths.");
  }
}

function getErrorMessage(error: unknown) {
  return error instanceof Error ? error.message : String(error);
}

async function fetchSwaggerDocument() {
  const response = await fetch(appConfig.swaggerDocsUrl);

  if (!response.ok) {
    throw new Error(
      `Failed to fetch Swagger docs from ${appConfig.swaggerDocsUrl}: ${response.status} ${response.statusText}`,
    );
  }

  const swaggerDocument = await response.json();
  validateSwaggerDocument(swaggerDocument);
  return swaggerDocument;
}

async function readFallbackSwaggerDocument() {
  if (!fallbackSwaggerInputPath) {
    throw new Error("No local Swagger fallback is configured.");
  }

  const swaggerDocument = JSON.parse(
    await readFile(fallbackSwaggerInputPath, "utf8"),
  );
  validateSwaggerDocument(swaggerDocument);
  return swaggerDocument;
}

async function main() {
  let swaggerDocument: SwaggerDocument;
  let source = appConfig.swaggerDocsUrl;

  try {
    swaggerDocument = await fetchSwaggerDocument();
  } catch (error) {
    source = fallbackSwaggerInputPath;
    console.warn(
      `Could not fetch Swagger docs from ${appConfig.swaggerDocsUrl}: ${getErrorMessage(error)}`,
    );
    console.warn(`Using local Swagger fallback -> ${fallbackSwaggerInputPath}`);
    swaggerDocument = await readFallbackSwaggerDocument();
  }

  await mkdir(path.dirname(swaggerOutputPath), { recursive: true });
  await writeFile(
    swaggerOutputPath,
    `${JSON.stringify(swaggerDocument, null, 2)}\n`,
    "utf8",
  );

  console.log(`Fetched Swagger docs from ${source} -> ${swaggerOutputPath}`);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
