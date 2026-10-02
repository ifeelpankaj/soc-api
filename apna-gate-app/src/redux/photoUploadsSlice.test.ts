import assert from "node:assert/strict";
import test from "node:test";
import reducer, {
  photoUploadFinished,
  photoUploadKey,
  photoUploadStarted,
  selectIsPhotoUploading,
} from "./photoUploadsSlice";

test("photo uploads are keyed by society and entry", () => {
  assert.equal(photoUploadKey({ societyId: 12, entryId: 34 }), "12:34");
});

test("concurrent uploads are tracked and finished independently", () => {
  const first = { societyId: 1, entryId: 10 };
  const second = { societyId: 1, entryId: 11 };
  let state = reducer(undefined, photoUploadStarted(first));
  state = reducer(state, photoUploadStarted(second));

  assert.deepEqual(state, { "1:10": true, "1:11": true });
  assert.equal(selectIsPhotoUploading({ photoUploads: state }, first), true);

  state = reducer(state, photoUploadFinished(first));
  assert.deepEqual(state, { "1:11": true });
  assert.equal(selectIsPhotoUploading({ photoUploads: state }, first), false);

  state = reducer(state, photoUploadFinished(second));
  assert.deepEqual(state, {});
});
