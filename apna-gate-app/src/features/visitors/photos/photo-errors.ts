export type PhotoUploadError = {
  status: number | string;
  code?: string;
  requestId?: string;
  category:
    | "permission"
    | "session"
    | "size"
    | "format"
    | "state"
    | "network"
    | "service";
  message: string;
  stage?: PhotoOperationStage;
  nativeMessage?: string;
};

export type PhotoOperationStage =
  | "prepare_file"
  | "build_form_data"
  | "dispatch"
  | "response";

export class PhotoOperationError extends Error {
  constructor(
    public readonly stage: PhotoOperationStage,
    message: string,
    public readonly userMessage: string,
  ) {
    super(message);
    this.name = "PhotoOperationError";
  }
}

export function photoUploadError(
  error: unknown,
  requestId?: string,
): PhotoUploadError {
  const input = (error && typeof error === "object" ? error : {}) as Record<
    string,
    unknown
  >;
  if (typeof input.category === "string" && typeof input.message === "string")
    return input as PhotoUploadError;
  const data = (
    input.data && typeof input.data === "object" ? input.data : {}
  ) as Record<string, unknown>;
  const details = (
    data.error && typeof data.error === "object" ? data.error : data
  ) as Record<string, unknown>;
  const code = typeof details.code === "string" ? details.code : undefined;
  const status =
    typeof input.status === "number" || typeof input.status === "string"
      ? input.status
      : "LOCAL_ERROR";
  const stage =
    typeof input.stage === "string"
      ? (input.stage as PhotoOperationStage)
      : status === "FETCH_ERROR" || status === "TIMEOUT_ERROR"
        ? "dispatch"
        : typeof input.status === "number"
          ? "response"
          : undefined;
  const nativeMessage =
    typeof input.error === "string"
      ? input.error
      : typeof input.message === "string"
        ? input.message
        : undefined;
  let category: PhotoUploadError["category"] = "service";
  let message = "Photo upload failed. Please try again.";
  if (code === "IMAGE_TOO_LARGE" || status === 413) {
    category = "size";
    message = "Choose a smaller photo (maximum 5 MiB and 25 megapixels).";
  } else if (
    code === "IMAGE_INVALID" ||
    code === "UNSUPPORTED_MEDIA_TYPE" ||
    status === 415
  ) {
    category = "format";
    message = "Choose a valid JPEG or PNG photo, or take a new photo.";
  } else if (code === "IMAGE_FORBIDDEN" || status === 403) {
    category = "permission";
    message = "You don't have permission to update this photo.";
  } else if (status === 401 || code === "UNAUTHORIZED") {
    category = "session";
    message = "Your session expired. Please sign in again.";
  } else if (code === "IMAGE_ENTRY_STATE") {
    category = "state";
    message = "This entry no longer allows photo changes.";
  } else if (status === "FETCH_ERROR" || status === "TIMEOUT_ERROR") {
    category = "network";
    message = "Couldn't upload the photo. Check your connection and try again.";
  } else if (input instanceof PhotoOperationError) {
    message = input.userMessage;
  } else if (
    code === "IMAGE_UNAVAILABLE" ||
    code === "IMAGE_BUSY" ||
    code === "IMAGE_PROVIDER_FAILED"
  ) {
    message =
      "Photo uploads are temporarily unavailable. Please try again shortly.";
  } else if (code === "IMAGE_PERSISTENCE_FAILED") {
    message =
      "Couldn't confirm the photo was saved. Refresh before trying again.";
  } else if (code === "IMAGE_TARGET_NOT_FOUND") {
    message = "This entry is no longer available. Refresh your entries.";
  }
  return { status, code, requestId, category, message, stage, nativeMessage };
}

export function logPhotoFailure(error: PhotoUploadError) {
  const { status, code, requestId, category, stage } = error;
  console.warn("Photo operation failed", {
    status,
    code,
    requestId,
    category,
    stage,
  });
}
