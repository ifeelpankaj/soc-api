import type { AppDispatch } from "@/redux/store";
import {
  photoUploadFinished,
  photoUploadStarted,
} from "@/redux/photoUploadsSlice";
import { photoApi } from "@/lib/api/photo-api";
import {
  photoFormData,
  releaseSelectedPhoto,
  type SelectedPhoto,
} from "./photo-picker";
import { runPhotoUploadLifecycle } from "./photo-upload-lifecycle";

type VisitorPhotoUploadTask = {
  societyId: number;
  entryId: number;
  photo: SelectedPhoto;
  photoReference?: string;
  onFailure: () => void;
};

export function startVisitorPhotoUpload(
  dispatch: AppDispatch,
  task: VisitorPhotoUploadTask,
) {
  const target = { societyId: task.societyId, entryId: task.entryId };
  void runPhotoUploadLifecycle({
    start: () => dispatch(photoUploadStarted(target)),
    upload: async () => {
      const body = await photoFormData(task.photo);
      await dispatch(
        photoApi.endpoints.putVisitorPhoto.initiate({ ...target, body, photoReference: task.photoReference }),
      ).unwrap();
    },
    onFailure: task.onFailure,
    cleanup: () => releaseSelectedPhoto(task.photo),
    finish: () => dispatch(photoUploadFinished(target)),
  });
}
