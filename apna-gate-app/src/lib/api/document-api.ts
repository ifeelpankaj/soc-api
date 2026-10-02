import { enhancedApi } from "@/lib/api/enhanced-api";

// Binary reads still pass through the app's token refresh and error handling.
// Keep Redux results serializable and clear document mutation results after use.
export async function binaryResponse(response: Response) {
  if (!response.ok) return response.json();
  return Array.from(new Uint8Array(await response.arrayBuffer()));
}
export const documentApi = enhancedApi.injectEndpoints({
  endpoints: (build) => ({
    readDocument: build.mutation<number[], { path: string }>({
      query: ({ path }) => ({
        url: path,
        method: "GET",
        responseHandler: binaryResponse,
      }),
    }),
  }),
});
export const { useReadDocumentMutation } = documentApi;
