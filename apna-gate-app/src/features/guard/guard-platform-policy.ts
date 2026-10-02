export function supportsGuardCompanions(platform: string, purpose?: string) {
  return purpose === "guest" || (platform === "android" && purpose === "cab");
}

export function usesCabOptionalCompanions(platform: string, purpose?: string) {
  return platform === "android" && purpose === "cab";
}
