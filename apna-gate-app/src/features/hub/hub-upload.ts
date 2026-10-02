import * as ImagePicker from "expo-image-picker";
import { Platform } from "react-native";

import {
  photoFormData,
  type SelectedPhoto,
} from "@/features/visitors/photos/photo-picker";

export {
  defaultCategoryId,
  HUB_MAX_ATTACHMENTS,
  hubImageSizeErrorMessage,
  normalizeAttachmentIds,
  validateHubImageMime,
  validateHubImageSize,
} from "@/features/hub/hub-upload-validation";

import {
  HUB_MAX_ATTACHMENTS,
  hubImageSizeErrorMessage,
  validateHubImageMime,
  validateHubImageSize,
} from "@/features/hub/hub-upload-validation";

export type HubPendingAttachment = {
  localUri: string;
  uploadId?: number;
  uploading?: boolean;
  error?: string;
};

export async function pickHubPhotos(maxCount: number): Promise<SelectedPhoto[]> {
  if (maxCount <= 0) {
    return [];
  }
  const permission = await ImagePicker.requestMediaLibraryPermissionsAsync();
  if (!permission.granted) {
    throw new Error("Photo library access is needed to add images.");
  }
  const result = await ImagePicker.launchImageLibraryAsync({
    mediaTypes: ["images"],
    allowsMultipleSelection: Platform.OS !== "web" && maxCount > 1,
    selectionLimit: maxCount,
    quality: 0.85,
  });
  if (result.canceled || !result.assets.length) {
    return [];
  }
  const photos: SelectedPhoto[] = [];
  for (const asset of result.assets.slice(0, maxCount)) {
    const mime =
      asset.mimeType?.toLowerCase() ??
      (/\.png$/i.test(asset.uri) ? "image/png" : "image/jpeg");
    if (!validateHubImageMime(mime)) {
      throw new Error(hubImageSizeErrorMessage());
    }
    let size = asset.fileSize ?? asset.file?.size ?? 0;
    if (!size && Platform.OS !== "web") {
      try {
        const { File } = await import("expo-file-system");
        const file = new File(asset.uri);
        size = file.size;
      } catch {
        size = 0;
      }
    }
    if (size && !validateHubImageSize(size)) {
      throw new Error(hubImageSizeErrorMessage());
    }
    photos.push({
      uri: asset.uri,
      name: mime === "image/png" ? "hub.png" : "hub.jpg",
      type: mime,
      file: asset.file,
    });
  }
  return photos;
}

export async function hubUploadFormData(photo: SelectedPhoto): Promise<FormData> {
  return photoFormData(photo);
}
