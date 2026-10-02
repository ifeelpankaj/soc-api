import type { HubCategory } from "@/lib/api/hub-api";

export const HUB_MAX_ATTACHMENTS = 5;
export const HUB_IMAGE_MAX_BYTES = 5 * 1024 * 1024;
export const HUB_PDF_MAX_BYTES = 10 * 1024 * 1024;
export const HUB_ALLOWED_IMAGE_MIMES = new Set(["image/jpeg", "image/png"]);

export function normalizeAttachmentIds(ids: number[] | undefined): number[] {
  const unique = [
    ...new Set(
      (ids ?? []).filter((id) => Number.isInteger(id) && id > 0),
    ),
  ];
  if (unique.length > HUB_MAX_ATTACHMENTS) {
    throw new Error(`You can attach up to ${HUB_MAX_ATTACHMENTS} files.`);
  }
  return unique;
}

export function validateHubImageMime(mime: string | undefined): boolean {
  if (!mime) {
    return false;
  }
  return HUB_ALLOWED_IMAGE_MIMES.has(mime.toLowerCase());
}

export function validateHubImageSize(bytes: number): boolean {
  return bytes > 0 && bytes <= HUB_IMAGE_MAX_BYTES;
}

export function hubImageSizeErrorMessage(): string {
  return "Each photo must be 5 MB or smaller (JPEG or PNG).";
}

export function defaultCategoryId(
  categories: HubCategory[] | undefined,
): number | undefined {
  if (!categories?.length) {
    return undefined;
  }
  const general =
    categories.find((c) => c.code === "general") ??
    categories.find((c) => c.name.toLowerCase() === "general");
  return general?.id ?? categories[0]?.id;
}
