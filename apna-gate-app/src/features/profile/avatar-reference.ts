export function normalizeAvatarReference(value?: string | null) {
  const normalized = value?.trim();
  return normalized || undefined;
}
