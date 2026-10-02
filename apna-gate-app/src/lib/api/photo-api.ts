import { enhancedApi } from "@/lib/api/enhanced-api";
import {
  logPhotoFailure,
  photoUploadError,
} from "@/features/visitors/photos/photo-errors";

import {
  photoCacheKey,
  photoPath,
  type PhotoQuery,
  type PhotoVariant,
} from "@/features/visitors/photos/photo-utils";

function uploadFailure(error: unknown, meta: unknown) {
  const response = (meta as { response?: Response } | undefined)?.response;
  const failure = photoUploadError(
    error,
    response?.headers.get("X-Request-ID") ?? undefined,
  );
  logPhotoFailure(failure);
  return failure;
}

export type PhotoView = {
  url: string;
  expires_at: string;
  variant: PhotoVariant;
};
type PhotoResponse = { success: boolean; data: PhotoView; message?: string };

const visitorDataInvalidationTags = [
  "Visitor Entries",
  "GuardDesk",
  "VisitorPending",
  "VisitorWaitingAtGate",
  "VisitorExpectedGuests",
  "FlatVisitorContext",
  "FlatVisitorEntries",
] as const;

type VisitorPhotoUpload = {
  societyId: number;
  entryId: number;
  body: FormData;
  photoReference?: string;
};

export const photoApi = enhancedApi
  .enhanceEndpoints({ addTagTypes: ["VisitorPhoto", "ProfilePhoto"] })
  .injectEndpoints({
    endpoints: (build) => ({
      getProfilePhoto: build.query<
        PhotoResponse,
        { userId: number; variant: PhotoVariant }
      >({
        query: ({ variant }) => ({
          url: "/v1/auth/profile/avatar",
          params: { variant },
        }),
        providesTags: ["ProfilePhoto"],
        keepUnusedDataFor: 60,
      }),
      putProfilePhoto: build.mutation<PhotoResponse, { body: FormData }>({
        query: ({ body }) => ({
          url: "/v1/auth/profile/avatar",
          method: "PUT",
          body,
        }),
        transformErrorResponse: uploadFailure,
        invalidatesTags: (result) =>
          result ? ["ProfilePhoto", "Auth", "Bootstrap"] : [],
      }),
      getVisitorPhoto: build.query<PhotoResponse, PhotoQuery>({
        query: ({ photoReference: _photoReference, ...query }) => ({
          url: photoPath(query),
          params: { variant: query.variant },
        }),
        serializeQueryArgs: ({ queryArgs }) =>
          `visitorPhoto:${photoCacheKey(queryArgs)}`,
        providesTags: (_result, _error, query) => [
          { type: "VisitorPhoto", id: query.photoReference },
        ],
        keepUnusedDataFor: 15 * 60,
      }),
      putVisitorPhoto: build.mutation<PhotoResponse, VisitorPhotoUpload>({
        query: ({ body, photoReference: _photoReference, ...target }) => ({
          url: photoPath({ ...target, context: "guard" }),
          method: "PUT",
          body,
        }),
        transformErrorResponse: uploadFailure,
        // A visitor may appear in multiple entries and authorization contexts;
        // the stored reference identifies all queries showing that same photo.
        invalidatesTags: (result, _error, argument) =>
          result
            ? [
                ...visitorDataInvalidationTags,
                ...(argument.photoReference
                  ? [
                      {
                        type: "VisitorPhoto" as const,
                        id: argument.photoReference,
                      },
                    ]
                  : []),
              ]
            : [],
      }),
    }),
  });

export const { useGetVisitorPhotoQuery, usePutVisitorPhotoMutation } = photoApi;
export const { useGetProfilePhotoQuery, usePutProfilePhotoMutation } = photoApi;
