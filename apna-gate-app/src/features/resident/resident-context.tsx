import {
  createContext,
  type PropsWithChildren,
  useContext,
  useMemo,
} from "react";

import {
  type ModelsFlatResidentResponse,
  type ModelsUserResponse,
  useGetV1BootstrapQuery,
} from "@/lib/api/generated-api";
import { setSelectedFlat } from "@/redux/appSlice";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import { saveWorkspace } from "@/features/auth/auth-storage";
import { PhotoAccessProvider } from "@/features/visitors/photos/photo-access";
import {
  canManageFlatMembers,
  canManageFlatVisitors,
} from "@/features/resident/resident-access";

type ResidentContextValue = {
  isLoading: boolean;
  residences: ModelsFlatResidentResponse[];
  selectedResidence?: ModelsFlatResidentResponse;
  societyId?: number;
  societyCode?: string;
  flatId?: number;
  isPrimary: boolean;
  canManageFlatVisitors: boolean;
  canManageFlatMembers: boolean;
  requiresSelection: boolean;
  selectResidence: (flatId: number) => void;
  refetch: () => void;
  user?: ModelsUserResponse | null;
};

const ResidentContext = createContext<ResidentContextValue | null>(null);

function isActiveResidence(residence: ModelsFlatResidentResponse) {
  return residence.status === "active" && typeof residence.flat_id === "number";
}

export function ResidentProvider({ children }: PropsWithChildren) {
  const dispatch = useAppDispatch();
  const selectedFlatId = useAppSelector((state) => state.app.selectedFlatId);
  const { data, isLoading, refetch } = useGetV1BootstrapQuery(undefined);

  const residences = useMemo(
    () => data?.data?.residences?.filter(isActiveResidence) ?? [],
    [data?.data?.residences],
  );

  const selectedResidence =
    residences.length === 1
      ? residences[0]
      : residences.find((residence) => residence.flat_id === selectedFlatId);

  const value = useMemo<ResidentContextValue>(
    () => ({
      isLoading,
      residences,
      selectedResidence,
      societyId: selectedResidence?.society_id,
      societyCode: selectedResidence?.society_code,
      flatId: selectedResidence?.flat_id,
      isPrimary: selectedResidence?.is_primary === true,
      canManageFlatVisitors: canManageFlatVisitors(selectedResidence),
      canManageFlatMembers: canManageFlatMembers(selectedResidence),
      requiresSelection: residences.length > 1 && !selectedResidence,
      selectResidence: (flatId: number) => {
        dispatch(setSelectedFlat(flatId));
        void saveWorkspace({ flatId });
      },
      refetch,
      user: data?.data?.user ?? null,
    }),
    [
      data?.data?.user,
      dispatch,
      isLoading,
      refetch,
      residences,
      selectedResidence,
    ],
  );

  return (
    <ResidentContext.Provider value={value}>
      <PhotoAccessProvider
        key={`${value.societyId}:${value.flatId}`}
        societyId={value.societyId}
        context={`resident:${value.flatId ?? 0}`}
      >
        {children}
      </PhotoAccessProvider>
    </ResidentContext.Provider>
  );
}

export function useResident() {
  const context = useContext(ResidentContext);

  if (!context) {
    throw new Error("useResident must be used inside ResidentProvider");
  }

  return context;
}
