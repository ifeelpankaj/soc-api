export type PhotoUploadLifecycle = {
  cleanup: () => void;
  finish: () => void;
  onFailure: () => void;
  start: () => void;
  upload: () => Promise<unknown>;
};

export async function runPhotoUploadLifecycle(task: PhotoUploadLifecycle) {
  task.start();
  try {
    await task.upload();
  } catch {
    task.onFailure();
  } finally {
    task.cleanup();
    task.finish();
  }
}

