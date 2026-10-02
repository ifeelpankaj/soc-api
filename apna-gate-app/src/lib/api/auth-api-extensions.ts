import { setCredentials } from "@/redux/authSlice";
import { enhancedApi } from "@/lib/api/enhanced-api";
import type { ModelsUserResponse } from "@/lib/api/generated-api";

type ApiEnvelope<T> = {
  data?: T;
  message?: string;
  success?: boolean;
};

export type UpdateProfileBody = {
  date_of_birth?: string;
  first_name?: string;
  gender?: string;
  language?: string;
  last_name?: string;
  phone_number?: string;
  timezone?: string;
};

export const authApiExtensions = enhancedApi.injectEndpoints({
  overrideExisting: true,
  endpoints: (build) => ({
    patchV1AuthProfile: build.mutation<
      ApiEnvelope<{ user?: ModelsUserResponse }>,
      UpdateProfileBody
    >({
      query: (body) => ({
        url: "/v1/auth/profile",
        method: "PATCH",
        body,
      }),
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          if (data?.data?.user) {
            dispatch(setCredentials({ user: data.data.user }));
          }
        } catch {
          // handled by caller
        }
      },
      invalidatesTags: ["Auth", "Bootstrap"],
    }),
  }),
});

export const { usePatchV1AuthProfileMutation } = authApiExtensions;
