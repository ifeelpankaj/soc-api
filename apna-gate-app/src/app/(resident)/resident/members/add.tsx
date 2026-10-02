import { Redirect, useRouter } from "expo-router";
import { useState } from "react";
import { StyleSheet, Text } from "react-native";

import {
  AddMemberFlatBadge,
  AddMemberSheet,
} from "@/features/resident/members/add-member-sheet";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";
import { useResidentMembers } from "@/features/resident/hooks/use-resident-members";
import { useResident } from "@/features/resident/resident-context";
import {
  residentInviteDetailRoute,
  residentMembersRoute,
} from "@/features/resident/resident-routes";
import { useBackAction } from "@/lib/navigation/use-back-action";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

export default function AddFlatMemberScreen() {
  const router = useRouter();
  const handleBack = useBackAction(residentMembersRoute());
  const { canManageFlatMembers, selectedResidence } = useResident();
  const { refetchAll } = useResidentMembers();
  const [isInviteSuccess, setIsInviteSuccess] = useState(false);

  if (!canManageFlatMembers) {
    return <Redirect href={residentMembersRoute()} />;
  }

  const flatNumber = selectedResidence?.flat_number;

  return (
    <AddMemberSheet
      renderLayout={(content, footer) => (
        <ResidentSubScreen
          footer={footer}
          headerExtra={
            isInviteSuccess && flatNumber ? (
              <Text style={styles.headerSubtitle}>Flat {flatNumber}</Text>
            ) : null
          }
          headerTrailing={<AddMemberFlatBadge flatNumber={flatNumber} />}
          title={isInviteSuccess ? "Invite Member" : "Add Member"}
        >
          {content}
        </ResidentSubScreen>
      )}
      onCreated={(inviteId) => {
        void refetchAll();
        if (inviteId) {
          router.replace(residentInviteDetailRoute("member", inviteId));
          return;
        }
        handleBack();
      }}
      onSuccessViewChange={setIsInviteSuccess}
    />
  );
}

const styles = StyleSheet.create({
  headerSubtitle: {
    color: colors.guard.textMuted,
    fontSize: 13,
    fontWeight: "500",
    marginTop: -spacing.xs,
    paddingBottom: spacing.xs,
  },
});
