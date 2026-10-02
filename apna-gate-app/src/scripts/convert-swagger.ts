import { mkdir, readFile, writeFile } from "node:fs/promises";
import { createRequire } from "node:module";
import path from "node:path";

const require = createRequire(import.meta.url);
const swagger2openapi = require("swagger2openapi") as {
  convertObj: (
    swagger: unknown,
    options: Record<string, unknown>,
  ) => Promise<{ openapi: OpenApiDocument }>;
};

const swaggerInputPath = path.resolve("src/openapi/swagger.json");
const openApiOutputPath = path.resolve("src/openapi/openapi.json");

type SwaggerDocument = {
  swagger?: string;
  paths?: Record<string, unknown>;
};

type OpenApiDocument = {
  openapi?: string;
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
    throw new Error(`${swaggerInputPath} must include swagger: "2.0".`);
  }

  if (
    !swaggerDocument.paths ||
    Object.keys(swaggerDocument.paths).length === 0
  ) {
    throw new Error(`${swaggerInputPath} must include non-empty paths.`);
  }
}

function validateOpenApiDocument(
  document: unknown,
): asserts document is OpenApiDocument {
  if (!document || typeof document !== "object") {
    throw new Error("OpenAPI output must be a JSON object.");
  }

  const openApiDocument = document as OpenApiDocument;
  if (!openApiDocument.openapi?.startsWith("3.")) {
    throw new Error('OpenAPI output must include openapi: "3.x".');
  }

  if (
    !openApiDocument.paths ||
    Object.keys(openApiDocument.paths).length === 0
  ) {
    throw new Error("OpenAPI output must include non-empty paths.");
  }
}

async function main() {
  const swaggerDocument = JSON.parse(await readFile(swaggerInputPath, "utf8"));
  validateSwaggerDocument(swaggerDocument);

  const { openapi } = await swagger2openapi.convertObj(swaggerDocument, {
    patch: true,
    warnOnly: true,
  });
  validateOpenApiDocument(openapi);

  await mkdir(path.dirname(openApiOutputPath), { recursive: true });
  await writeFile(
    openApiOutputPath,
    `${JSON.stringify(openapi, null, 2)}\n`,
    "utf8",
  );

  console.log(`Converted Swagger 2.0 -> ${openApiOutputPath}`);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
