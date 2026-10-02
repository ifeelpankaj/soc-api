import { useRouter } from "expo-router";
import { useMemo, useState } from "react";

import { titleize } from "@/features/guard/guard-utils";
import {
  ProfileActionRow,
  ProfileSection,
  ResidenceAccessCard,
  ResidenceSwitchSheet,
} from "@/features/profile";
import { ProfileEditSheet } from "@/features/profile/components/profile-edit-sheet";
import type { ProfileScreenProps } from "@/features/profile/components/profile-screen";
import { ResidentSocietyGate } from "@/features/resident/components/resident-society-gate";
import { useResident } from "@/features/resident/resident-context";
import {
  residentDashboardRoute,
  residentMembersRoute,
  residentVisitorSettingsRoute,
} from "@/features/resident/resident-routes";

function formatFlatLabel(residence?: {
  block?: string | null;
  flat_number?: string | null;
}) {
  if (!residence) {
    return "Flat";
  }

  const parts = [
    residence.block ? `Block ${residence.block}` : null,
    residence.flat_number ? `Flat ${residence.flat_number}` : null,
  ].filter(Boolean);

  return parts.length > 0 ? parts.join(" · ") : "Flat";
}

export function useResidentProfileModel(): ProfileScreenProps {
  const router = useRouter();
  const [switchSheetVisible, setSwitchSheetVisible] = useState(false);
  const [editSheetVisible, setEditSheetVisible] = useState(false);
  const {
    canManageFlatVisitors,
    flatId,
    isLoading,
    isPrimary,
    refetch,
    requiresSelection,
    residences,
    selectResidence,
    selectedResidence,
    user,
  } = useResident();

  const societyName = selectedResidence?.society_name ?? "Your society";
  const flatLabel = formatFlatLabel(selectedResidence);

  const gate = useMemo(
    () => (!isLoading && requiresSelection ? <ResidentSocietyGate /> : undefined),
    [isLoading, requiresSelection],
  );

  const workspaceSlot = useMemo(
    () =>
      selectedResidence ? (
        <ResidenceAccessCard
          flatLabel={flatLabel}
          isPrimary={isPrimary}
          residenceCount={residences.length}
          societyName={societyName}
          status={selectedResidence.status}
          onSwitchPress={() => setSwitchSheetVisible(true)}
        />
      ) : null,
    [flatLabel, isPrimary, residences.length, selectedResidence, societyName],
  );

  const extraSections = useMemo(
    () => (
      <ProfileSection title="Residence">
        <ProfileActionRow
          description={flatLabel}
          icon={{ ios: "house.fill", android: "home", web: "home" }}
          label={societyName}
          onPress={() => {
            if (residences.length > 1) {
              setSwitchSheetVisible(true);
            }
          }}
        />
        {canManageFlatVisitors ? (
          <ProfileActionRow
            icon={{ ios: "slider.horizontal.3", android: "tune", web: "tune" }}
            label="Visitor settings"
            onPress={() => router.push(residentVisitorSettingsRoute())}
          />
        ) : null}
        <ProfileActionRow
          icon={{ ios: "person.2.fill", android: "group", web: "group" }}
          isLast
          label="Flat members"
          onPress={() => router.push(residentMembersRoute())}
        />
      </ProfileSection>
    ),
    [canManageFlatVisitors, flatLabel, residences.length, router, societyName],
  );

  const modals = useMemo(
    () => (
      <>
        <ResidenceSwitchSheet
          residences={residences}
          selectedFlatId={flatId}
          visible={switchSheetVisible}
          onClose={() => setSwitchSheetVisible(false)}
          onSelect={selectResidence}
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
      flatId,
      refetch,
      residences,
      selectResidence,
      switchSheetVisible,
      user,
    ],
  );

  return useMemo(
    () => ({
      defaultName: "Resident",
      extraSections,
      fallbackHomeRoute: residentDashboardRoute(),
      gate,
      isLoading,
      modals,
      onEditProfile: () => setEditSheetVisible(true),
      onRefreshAccess: refetch,
      supportAboutDescription: `${titleize(selectedResidence?.role ?? "resident")} access`,
      user,
      workspaceSlot,
    }),
    [extraSections, gate, isLoading, modals, refetch, selectedResidence?.role, user, workspaceSlot],
  );
}
