import { useState } from "react";
import { Alert, Pressable, ScrollView, StyleSheet, View } from "react-native";
import { SymbolView } from "expo-symbols";
import { useLocalSearchParams, useRouter } from "expo-router";

import { Button, EmptyState, LoadingState } from "@/components/ui";
import { GuardEntryEditSheet } from "@/features/guard/components/guard-entry-edit-sheet";
import { GuardBottomActions } from "@/features/guard/components/guard-bottom-actions";
import { visitorDetailsFooterState } from "@/features/guard/guard-footer-state";
import { GuardSubScreen } from "@/features/guard/components/guard-sub-screen";
import { GuardVisitorDetailView } from "@/features/guard/components/guard-visitor-detail-view";
import { canEditVisitorEntry } from "@/features/guard/guard-entry-edit";
import { guardCheckInRoute } from "@/features/guard/guard-routes";
import { useGuardActions } from "@/features/guard/hooks/use-guard-actions";
import { useGuardFeedback } from "@/features/guard/hooks/use-guard-feedback";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import { useGetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdQuery } from "@/lib/api/generated-api";
import { useAppSelector } from "@/redux/hooks";
import { selectIsPhotoUploading } from "@/redux/photoUploadsSlice";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

function firstParam(value?: string | string[]) {
  return Array.isArray(value) ? value[0] : value;
}

export default function GuardEntryDetailRoute() {
  const router = useRouter();
  const { entryId: entryIdParam } = useLocalSearchParams<{
    entryId?: string | string[];
  }>();
  const entryId = Number(firstParam(entryIdParam));
  const { selectedSocietyId } = useGuardScreen();
  const actions = useGuardActions(selectedSocietyId ?? 0);
  const feedback = useGuardFeedback();
  const [editVisible, setEditVisible] = useState(false);
  const isPhotoUploading = useAppSelector((state) =>
    selectedSocietyId && Number.isFinite(entryId)
      ? selectIsPhotoUploading(state, {
          societyId: selectedSocietyId,
          entryId,
        })
      : false,
  );

  const query = useGetV1SocietiesBySocietyIdVisitorEntriesAndEntryIdQuery(
    { societyId: selectedSocietyId ?? 0, entryId },
    { skip: !selectedSocietyId || !Number.isFinite(entryId) || entryId <= 0 },
  );

  const entry = query.data?.data?.entry;
  const editable = canEditVisitorEntry(entry);
  const canCheckIn = entry?.status === "approved";

  const handleCheckOut = async () => {
    if (!entry?.id) {
      return;
    }

    const result = await actions.checkOutEntry(entry.id);
    feedback.showActionResult(result, {
      successTitle: "Checked out",
      errorTitle: "Checkout failed",
    });

    if (result.success) {
      void query.refetch();
    }
  };

  const handleCheckIn = async () => {
    if (!entry?.id) {
      return;
    }

    router.push(guardCheckInRoute({ source: "entry", entryId: entry.id }));
  };

  const handleMenuPress = () => {
    const menuOptions: {
      text: string;
      onPress?: () => void;
      style?: "cancel" | "destructive";
    }[] = [{ text: "Cancel", style: "cancel" }];

    if (editable) {
      menuOptions.unshift({
        text: "Edit Details",
        onPress: () => setEditVisible(true),
      });
    }

    if (canCheckIn) {
      menuOptions.unshift({
        text: "Review & Check In",
        onPress: () => void handleCheckIn(),
      });
    }

    if (entry?.status === "checked_in") {
      menuOptions.unshift({
        text: "Check Out",
        onPress: () => void handleCheckOut(),
      });
    }

    if (menuOptions.length === 1) {
      Alert.alert(
        "Visitor Actions",
        "No actions are available for this visitor.",
      );
      return;
    }

    Alert.alert("Visitor Actions", undefined, menuOptions);
  };

  const headerTrailing = (
    <Pressable
      accessibilityLabel="More actions"
      accessibilityRole="button"
      hitSlop={8}
      style={({ pressed }) => [
        styles.menuButton,
        pressed && styles.menuButtonPressed,
      ]}
      onPress={handleMenuPress}
    >
      <SymbolView
        name={{ ios: "ellipsis", android: "more_vert", web: "more_vert" }}
        size={18}
        tintColor={colors.guard.text}
      />
    </Pressable>
  );

  const footerState = visitorDetailsFooterState(entry?.status);
  const footer = entry && footerState ? (
    <GuardBottomActions>
      {footerState === "check_out" ? (
        <Button title="Check Out" loading={actions.activeEntryId === entry.id} onPress={() => void handleCheckOut()} />
      ) : (
        <>
          <Button title="Check In" onPress={() => void handleCheckIn()} />
          {editable ? <Button title="Edit Details" variant="secondary" onPress={() => setEditVisible(true)} /> : null}
        </>
      )}
    </GuardBottomActions>
  ) : undefined;

  return (
    <GuardSubScreen headerTrailing={headerTrailing} title="Visitor Details" footer={footer}>
      {query.isLoading ? (
        <LoadingState message="Loading visitor details" />
      ) : !entry ? (
        <View style={styles.emptyWrap}>
          <EmptyState
            actionLabel="Retry"
            message="This visitor entry could not be loaded."
            title="Entry not found"
            onAction={() => void query.refetch()}
          />
        </View>
      ) : (
        <ScrollView
          style={{ flex: 1 }}
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          <GuardVisitorDetailView
            entry={entry}
            isPhotoUploading={isPhotoUploading}
          />

        </ScrollView>
      )}

      <GuardEntryEditSheet
        entry={entry}
        societyId={selectedSocietyId ?? 0}
        visible={editVisible}
        onClose={() => setEditVisible(false)}
        onSaved={() => {
          setEditVisible(false);
          void query.refetch();
          feedback.showSuccess("Details updated", "Visitor information saved.");
        }}
      />
    </GuardSubScreen>
  );
}

const styles = StyleSheet.create({
  emptyWrap: {
    padding: spacing.lg,
  },
  menuButton: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius["2xl"],
    borderWidth: 1,
    height: 44,
    justifyContent: "center",
    width: 44,
  },
  menuButtonPressed: {
    opacity: 0.85,
  },
  scrollContent: {
    flexGrow: 1,
  },
});
