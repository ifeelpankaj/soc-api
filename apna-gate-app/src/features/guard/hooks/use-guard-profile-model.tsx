import { useMemo, useState } from "react";

import { GuardSocietyGate } from "@/features/guard/components/guard-society-gate";
import { useGuardSociety } from "@/features/guard/guard-context";
import { guardHomeRoute } from "@/features/guard/guard-routes";
import { SocietyAccessCard, SocietySwitchSheet } from "@/features/profile";
import { ProfileEditSheet } from "@/features/profile/components/profile-edit-sheet";
import type { ProfileScreenProps } from "@/features/profile/components/profile-screen";
import { useGetV1SocietiesBySocietyIdQuery } from "@/lib/api/generated-api";

export function useGuardProfileModel(): ProfileScreenProps {
  const {
    isLoading,
    memberships,
    refetch,
    requiresSelection,
    selectSociety,
    selectedMembership,
    selectedSocietyId,
    user,
  } = useGuardSociety();
  const [switchSheetVisible, setSwitchSheetVisible] = useState(false);
  const [editSheetVisible, setEditSheetVisible] = useState(false);

  const societyQuery = useGetV1SocietiesBySocietyIdQuery(
    { societyId: selectedSocietyId ?? 0 },
    { skip: !selectedSocietyId },
  );

  const societyName =
    societyQuery.data?.data?.society?.name ??
    (selectedSocietyId ? `Society #${selectedSocietyId}` : "No society selected");

  const gate = useMemo(
    () =>
      !isLoading && (memberships.length === 0 || requiresSelection) ? (
        <GuardSocietyGate />
      ) : undefined,
    [isLoading, memberships.length, requiresSelection],
  );

  const workspaceSlot = useMemo(
    () =>
      selectedSocietyId ? (
        <SocietyAccessCard
          membershipCount={memberships.length}
          role={selectedMembership?.role}
          societyName={societyName}
          status={selectedMembership?.status}
          onSwitchPress={() => setSwitchSheetVisible(true)}
        />
      ) : null,
    [
      memberships.length,
      selectedMembership?.role,
      selectedMembership?.status,
      selectedSocietyId,
      societyName,
    ],
  );

  const modals = useMemo(
    () => (
      <>
        <SocietySwitchSheet
          memberships={memberships}
          selectedSocietyId={selectedSocietyId}
          visible={switchSheetVisible}
          onClose={() => setSwitchSheetVisible(false)}
          onSelect={selectSociety}
        />
        <ProfileEditSheet
          user={user}
          visible={editSheetVisible}
          onClose={() => setEditSheetVisible(false)}
          onSaved={() => void refetch()}
        />
      </>
    ),
    [
      editSheetVisible,
      memberships,
      refetch,
      selectSociety,
      selectedSocietyId,
      switchSheetVisible,
      user,
    ],
  );

  return useMemo(
    () => ({
      defaultName: "Guard",
      fallbackHomeRoute: guardHomeRoute(),
      gate,
      isLoading,
      modals,
      onEditProfile: () => setEditSheetVisible(true),
      onRefreshAccess: refetch,
      user,
      workspaceSlot,
    }),
    [gate, isLoading, modals, refetch, user, workspaceSlot],
  );
}
