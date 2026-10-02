export type PhotoVariant = "list" | "avatar" | "detail" | "original";
export type PhotoContext = "guard" | `resident:${number}`;
export type PhotoQuery = {
  userId: number;
  societyId: number;
  entryId: number;
  context: PhotoContext;
  variant: PhotoVariant;
  photoReference: string;
};

export const MAX_PHOTO_BYTES = 5 * 1024 * 1024;
export const MAX_PHOTO_PIXELS = 25_000_000;
export const PHOTO_URL_REFRESH_MARGIN_MS = 3 * 60 * 1000;

export function photoCacheKey(query: PhotoQuery) {
  return JSON.stringify([
    query.userId,
    query.societyId,
    query.context,
    query.entryId,
    query.photoReference,
    query.variant,
  ]);
}

export function photoPath(
  query: Pick<PhotoQuery, "societyId" | "entryId" | "context">,
) {
  const prefix = `/v1/societies/${query.societyId}`;
  const flat = query.context.startsWith("resident:")
    ? `/flats/${query.context.slice("resident:".length)}`
    : "";
  return `${prefix}${flat}/visitor-entries/${query.entryId}/photo`;
}

export function photoRefreshDelay(expiresAt: string, now = Date.now()) {
  return Math.max(
    0,
    Date.parse(expiresAt) - now - PHOTO_URL_REFRESH_MARGIN_MS,
  );
}

export function photoNeedsRefresh(expiresAt?: string, now = Date.now()) {
  if (!expiresAt) return true;
  const expires = Date.parse(expiresAt);
  return !Number.isFinite(expires) || expires <= now + PHOTO_URL_REFRESH_MARGIN_MS;
}

export function photoPreparationNeeded(
  mime: string,
  size: number,
  width: number,
  height: number,
) {
  return (
    !["image/jpeg", "image/png"].includes(mime) ||
    size > MAX_PHOTO_BYTES ||
    width * height > MAX_PHOTO_PIXELS
  );
}

export function photoDisplayState(
  hasPhoto: boolean,
  loading: boolean,
  failed: boolean,
) {
  return failed || (!hasPhoto && !loading)
    ? "fallback"
    : loading
      ? "loading"
      : "photo";
}
