import { getApiErrorCode, getFriendlyApiMessage } from "@/features/auth/api-error";

export function hubUploadErrorMessage(error: unknown): string {
  const code = getApiErrorCode(error);
  if (code === "IMAGE_UNAVAILABLE") {
    return "Photo uploads are not set up on this server (enable ImageKit in API .env). You can still publish without photos.";
  }
  return getFriendlyApiMessage(error, "Could not upload this photo. Remove it or try again.");
}

export function hubPublishErrorMessage(error: unknown): string {
  const code = getApiErrorCode(error);
  if (code === "HUB_FORBIDDEN") {
    return "You cannot publish this post. Use Community (not Announcements) and avoid admin-only options.";
  }
  return getFriendlyApiMessage(error, "Please try again.");
}
