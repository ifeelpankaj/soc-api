import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

type PhotoUploadTarget = {
  societyId: number;
  entryId: number;
};

export function photoUploadKey({ societyId, entryId }: PhotoUploadTarget) {
  return `${societyId}:${entryId}`;
}

type PhotoUploadsState = Record<string, true>;

type PhotoUploadsRootState = {
  photoUploads: PhotoUploadsState;
};

const initialState: PhotoUploadsState = {};

const photoUploadsSlice = createSlice({
  name: "photoUploads",
  initialState,
  reducers: {
    photoUploadStarted(state, action: PayloadAction<PhotoUploadTarget>) {
      state[photoUploadKey(action.payload)] = true;
    },
    photoUploadFinished(state, action: PayloadAction<PhotoUploadTarget>) {
      delete state[photoUploadKey(action.payload)];
    },
  },
});

export const { photoUploadFinished, photoUploadStarted } =
  photoUploadsSlice.actions;

export function selectIsPhotoUploading(
  state: PhotoUploadsRootState,
  target: PhotoUploadTarget,
) {
  return Boolean(state.photoUploads[photoUploadKey(target)]);
}

export default photoUploadsSlice.reducer;
