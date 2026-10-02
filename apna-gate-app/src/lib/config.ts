function requireEnv(key: string, value: string | undefined) {
  if (!value) {
    throw new Error(`Missing required environment variable: ${key}`);
  }

  return value;
}

const apiBaseUrl = process.env.EXPO_PUBLIC_API_BASE_URL;
const swaggerDocsUrl = process.env.EXPO_PUBLIC_SWAGGER_URL;
const webBaseUrl = process.env.EXPO_PUBLIC_WEB_BASE_URL;

export const appConfig = {
  get apiBaseUrl() {
    return requireEnv("EXPO_PUBLIC_API_BASE_URL", apiBaseUrl);
  },
  get swaggerDocsUrl() {
    return requireEnv("EXPO_PUBLIC_SWAGGER_URL", swaggerDocsUrl);
  },
  get webBaseUrl() {
    return requireEnv("EXPO_PUBLIC_WEB_BASE_URL", webBaseUrl).replace(/\/+$/, "");
  },
};

export function buildVisitorInviteUrl(societyCode: string, token: string) {
  return `${appConfig.webBaseUrl}/visit/${encodeURIComponent(societyCode)}/invite/${encodeURIComponent(token)}`;
}

export function buildMemberInviteUrl(token: string) {
  return `${appConfig.webBaseUrl}/join/flat/${encodeURIComponent(token)}`;
}
