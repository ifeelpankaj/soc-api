import {
  createContext,
  type PropsWithChildren,
  useContext,
  useMemo,
} from "react";

import {
  type ModelsSocietyMemberResponse,
  type ModelsUserResponse,
  useGetV1BootstrapQuery,
} from "@/lib/api/generated-api";
import { saveWorkspace } from "@/features/auth/auth-storage";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import { setSelectedSociety } from "@/redux/appSlice";
import { PhotoAccessProvider } from "@/features/visitors/photos/photo-access";

type GuardSocietyContextValue = {
  isLoading: boolean;
  memberships: ModelsSocietyMemberResponse[];
  selectedMembership?: ModelsSocietyMemberResponse;
  selectedSocietyId?: number;
  requiresSelection: boolean;
  selectSociety: (societyId: number) => void;
  refetch: () => void;
  user?: ModelsUserResponse | null;
};

const GuardSocietyContext = createContext<GuardSocietyContextValue | null>(
  null,
);

function isGuardMembership(membership: ModelsSocietyMemberResponse) {
  return (
    membership.status === "active" &&
    typeof membership.society_id === "number" &&
    membership.role === "staff"
  );
}

export function GuardSocietyProvider({ children }: PropsWithChildren) {
  const dispatch = useAppDispatch();
  const selectedSocietyId = useAppSelector(
    (state) => state.app.selectedSocietyId,
  );
  const { data, isLoading, refetch } = useGetV1BootstrapQuery(undefined);

  const memberships = useMemo(
    () => data?.data?.memberships?.filter(isGuardMembership) ?? [],
    [data?.data?.memberships],
  );

  const selectedMembership =
    memberships.length === 1
      ? memberships[0]
      : memberships.find(
          (membership) => membership.society_id === selectedSocietyId,
        );

  const value = useMemo<GuardSocietyContextValue>(
    () => ({
      isLoading,
      memberships,
      selectedMembership,
      selectedSocietyId: selectedMembership?.society_id,
      requiresSelection: memberships.length > 1 && !selectedMembership,
      selectSociety: (societyId: number) => {
        dispatch(setSelectedSociety(societyId));
        void saveWorkspace({ societyId });
      },
      refetch,
      user: data?.data?.user ?? null,
    }),
    [
      data?.data?.user,
      dispatch,
      isLoading,
      memberships,
      refetch,
      selectedMembership,
    ],
  );

  return (
    <GuardSocietyContext.Provider value={value}>
      <PhotoAccessProvider societyId={value.selectedSocietyId} context="guard">
        {children}
      </PhotoAccessProvider>
    </GuardSocietyContext.Provider>
  );
}

export function useGuardSociety() {
  const context = useContext(GuardSocietyContext);

  if (!context) {
    throw new Error("useGuardSociety must be used inside GuardSocietyProvider");
  }

  return context;
}
