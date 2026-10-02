import { useEffect, useMemo, useRef, useState } from "react";
import {
  ActivityIndicator,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { SymbolView } from "expo-symbols";
import { useLocalSearchParams, useRouter } from "expo-router";

import { Button, Card } from "@/components/ui";
import { GuardEntryEditSheet } from "@/features/guard/components/guard-entry-edit-sheet";
import { GuardBottomActions } from "@/features/guard/components/guard-bottom-actions";
import { checkInFooterState } from "@/features/guard/guard-footer-state";
import { GuardSubScreen } from "@/features/guard/components/guard-sub-screen";
import {
  canEditVisitorEntry,
  getCheckInSessionKey,
} from "@/features/guard/guard-entry-edit";
import {
  firstParam,
  guardScannerRoute,
  guardWaitingAtGateRoute,
  parseCheckInParams,
} from "@/features/guard/guard-routes";
import {
  type GuardScanOutcome,
  useGuardCheckIn,
} from "@/features/guard/hooks/use-guard-check-in";
import { useGuardFeedback } from "@/features/guard/hooks/use-guard-feedback";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import { VisitorDetailsCard } from "@/features/visitors/components/visitor-details-card";
import { useAppSelector } from "@/redux/hooks";
import { selectIsPhotoUploading } from "@/redux/photoUploadsSlice";
import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

function getSubtitle(outcome: GuardScanOutcome) {
  switch (outcome) {
    case "already_inside":
      return "This visitor is already inside the society.";
    case "just_checked_in":
      return "Check-in complete. The visitor may proceed to entry.";
    case "pending_approval":
      return "This visitor is waiting for resident approval.";
    case "loading":
      return "Review visitor details before allowing entry.";
    default:
      return "Review visitor details before allowing entry.";
  }
}

export default function GuardCheckInScreen() {
  const router = useRouter();
  const feedback = useGuardFeedback();
  const params = useLocalSearchParams<{
    source?: string | string[];
    token?: string | string[];
    entryId?: string | string[];
  }>();
  const { selectedSocietyId } = useGuardScreen();

  const sourceParam = firstParam(params.source);
  const tokenParam = firstParam(params.token);
  const entryIdParam = firstParam(params.entryId);
  const checkInInput = useMemo(
    () =>
      parseCheckInParams({
        source: sourceParam,
        token: tokenParam,
        entryId: entryIdParam,
      }),
    [entryIdParam, sourceParam, tokenParam],
  );
  const checkIn = useGuardCheckIn(checkInInput, selectedSocietyId);
  const [editVisible, setEditVisible] = useState(false);
  const [editedEntry, setEditedEntry] = useState<{
    key: string;
    entry: ModelsVisitorEntry;
  } | null>(null);
  const checkInToastShownRef = useRef(false);
  const wasCheckingInRef = useRef(false);
  const redirectTriggeredRef = useRef(false);

  const sessionKey = useMemo(
    () => (checkInInput ? getCheckInSessionKey(checkInInput) : ""),
    [checkInInput],
  );

  const completed = checkIn.scanOutcome === "just_checked_in" || checkIn.scanOutcome === "already_inside";
  const entry = !completed && editedEntry?.key === sessionKey
    ? editedEntry.entry
    : checkIn.entry;
  const isPhotoUploading = useAppSelector((state) => Platform.OS === "android" &&
    selectIsPhotoUploading(state, { societyId: selectedSocietyId ?? 0, entryId: entry?.id ?? 0 }));

  useEffect(() => {
    checkInToastShownRef.current = false;
    wasCheckingInRef.current = false;
    redirectTriggeredRef.current = false;
  }, [sessionKey]);

  useEffect(() => {
    if (checkIn.isCheckingIn) {
      wasCheckingInRef.current = true;
      return;
    }

    if (
      wasCheckingInRef.current &&
      checkIn.scanOutcome === "just_checked_in" &&
      !checkInToastShownRef.current
    ) {
      checkInToastShownRef.current = true;
      wasCheckingInRef.current = false;
      feedback.showSuccess("Checked in", "The visitor may proceed to entry.");
    }
  }, [checkIn.isCheckingIn, checkIn.scanOutcome, feedback]);

  useEffect(() => {
    if (
      checkIn.entryError?.redirectToWaitingAtGate &&
      !checkIn.entry &&
      !checkIn.isLoadingEntry &&
      !redirectTriggeredRef.current
    ) {
      redirectTriggeredRef.current = true;
      feedback.showError("QR not recognized", "", checkIn.entryError.message);
      router.replace(guardWaitingAtGateRoute());
    }
  }, [
    checkIn.entry,
    checkIn.entryError,
    checkIn.isLoadingEntry,
    feedback,
    router,
  ]);

  useEffect(() => {
    if (
      checkIn.entryError &&
      checkIn.entry &&
      !checkIn.isCheckingIn &&
      checkIn.entryError.kind !== "invalid_params" &&
      checkIn.entryError.kind !== "validation_failed"
    ) {
      feedback.showActionResult(
        { success: false, message: checkIn.entryError.message },
        { errorTitle: "Check-in failed", successTitle: "Checked in" },
      );
    }
  }, [checkIn.entry, checkIn.entryError, checkIn.isCheckingIn, feedback]);

  if (!checkInInput) {
    return (
      <GuardSubScreen title="Check In">
        <View style={styles.content}>
          <Card>
            <Text style={styles.cardTitle}>Invalid check-in link</Text>
            <Text style={styles.cardBody}>
              This check-in route is missing a valid visitor reference. Scan a visitor QR
              or open the entry from Add Entry to continue.
            </Text>
          </Card>
          <Pressable
            accessibilityRole="button"
            style={styles.primaryButton}
            onPress={() => router.replace(guardScannerRoute())}
          >
            <Text style={styles.primaryButtonText}>Go to scanner</Text>
          </Pressable>
        </View>
      </GuardSubScreen>
    );
  }

  const { scanOutcome } = checkIn;
  const footerState = checkInFooterState(scanOutcome, entry?.status, checkIn.isCheckingIn);
  const showEditDetails = canEditVisitorEntry(entry);

  const footer = footerState ? (
    <GuardBottomActions>
      {footerState === "scan" ? (
        <Button title="Scan another" onPress={() => router.replace(guardScannerRoute())} />
      ) : (
        <>
          <Button title="Check In" loading={checkIn.isCheckingIn} onPress={() => void checkIn.checkIn()} />
          {showEditDetails ? (
            <Button title="Edit Details" variant="secondary" disabled={checkIn.isCheckingIn} onPress={() => setEditVisible(true)} />
          ) : null}
        </>
      )}
    </GuardBottomActions>
  ) : undefined;

  return (
    <GuardSubScreen title="Check In" footer={footer}>
      <ScrollView style={{ flex: 1 }} contentContainerStyle={styles.scrollContent}>
        <View style={styles.content}>
          <View style={styles.heroHeader}>
            <View style={styles.heroIconWrap}>
              <SymbolView
                name={{ ios: "shield.fill", android: "security", web: "security" }}
                size={28}
                tintColor={colors.brand.orange}
              />
            </View>
            <Text style={styles.heroTitle}>Visitor Check-In</Text>
            <Text style={styles.subtitle}>{getSubtitle(scanOutcome)}</Text>
          </View>

          {scanOutcome === "loading" ? (
            <View style={styles.inlineLoading}>
              <ActivityIndicator color={colors.guard.teal} />
              <Text style={styles.cardBody}>Loading visitor...</Text>
            </View>
          ) : null}

          {scanOutcome === "error" &&
          checkIn.entryError &&
          !checkIn.entryError.redirectToWaitingAtGate ? (
            <Card style={styles.errorCard}>
              <Text style={styles.errorTitle}>Unable to load visitor</Text>
              <Text style={styles.errorBody}>{checkIn.entryError.message}</Text>
            </Card>
          ) : null}

          {entry ? <VisitorDetailsCard entry={entry} variant="compact" /> : null}
          {isPhotoUploading ? (
            <View accessibilityLiveRegion="polite" style={styles.inlineLoading}>
              <ActivityIndicator color={colors.brand.orange} />
              <Text style={styles.cardBody}>Uploading photo… You can continue check-in.</Text>
            </View>
          ) : null}

          {scanOutcome === "pending_approval" ? (
            <Card style={styles.warningCard}>
              <Text style={styles.warningTitle}>Waiting for approval</Text>
              <Text style={styles.warningBody}>
                {checkIn.disabledReason ??
                  "This visitor is not approved for check-in yet."}
              </Text>
            </Card>
          ) : null}

          {scanOutcome === "blocked" && checkIn.disabledReason ? (
            <Card style={styles.warningCard}>
              <Text style={styles.warningTitle}>Check-in unavailable</Text>
              <Text style={styles.warningBody}>{checkIn.disabledReason}</Text>
            </Card>
          ) : null}

          {scanOutcome === "already_inside" ? (
            <Card style={styles.successCard}>
              <Text style={styles.successTitle}>Already checked in</Text>
              <Text style={styles.successBody}>
                This visitor is already inside. No further action is needed.
              </Text>
            </Card>
          ) : null}

          {scanOutcome === "just_checked_in" ? (
            <Card style={styles.successCard}>
              <Text style={styles.successTitle}>Checked in successfully</Text>
              <Text style={styles.successBody}>
                The visitor may proceed to entry.
              </Text>
            </Card>
          ) : null}

          {checkIn.entryError &&
          entry &&
          scanOutcome !== "just_checked_in" ? (
            <Card style={styles.errorCard}>
              <Text style={styles.errorTitle}>Check-in failed</Text>
              <Text style={styles.errorBody}>{checkIn.entryError.message}</Text>
            </Card>
          ) : null}

        </View>
      </ScrollView>

      <GuardEntryEditSheet
        entry={entry}
        societyId={selectedSocietyId ?? 0}
        visible={editVisible}
        onClose={() => setEditVisible(false)}
        onSaved={(updated) => {
          if (sessionKey) {
            setEditedEntry({ key: sessionKey, entry: updated });
          }
          checkIn.refreshEntry();
          feedback.showSuccess("Details updated", "Visitor information saved.");
        }}
      />
    </GuardSubScreen>
  );
}

const styles = StyleSheet.create({
  cardBody: {
    ...typography.bodySmall,
    color: colors.text.secondary,
    marginTop: spacing.xs,
  },
  cardTitle: {
    ...typography.body,
    color: colors.text.primary,
    fontWeight: "700",
  },
  content: {
    gap: spacing.lg,
    paddingBottom: layout.screenPaddingBottom,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.md,
  },
  errorBody: {
    ...typography.bodySmall,
    color: colors.status.error,
    marginTop: spacing.xs,
  },
  errorCard: {
    backgroundColor: colors.status.errorSoft,
    borderColor: "#fecaca",
  },
  errorTitle: {
    ...typography.body,
    color: colors.status.error,
    fontWeight: "700",
  },
  heroHeader: {
    alignItems: "center",
    gap: spacing.xs,
  },
  heroIconWrap: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.full,
    height: 48,
    justifyContent: "center",
    width: 48,
  },
  heroTitle: {
    ...typography.title,
    color: colors.text.primary,
    fontSize: 22,
    fontWeight: "800",
  },
  inlineLoading: {
    alignItems: "center",
    gap: spacing.sm,
    paddingVertical: spacing["2xl"],
  },
  primaryButton: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: 14,
    paddingVertical: spacing.md,
  },
  primaryButtonText: {
    ...typography.button,
    color: colors.text.inverse,
  },
  scrollContent: {
    flexGrow: 1,
  },
  subtitle: {
    ...typography.bodySmall,
    color: colors.text.secondary,
    textAlign: "center",
  },
  successBody: {
    ...typography.bodySmall,
    color: "#166534",
    marginTop: spacing.xs,
  },
  successCard: {
    backgroundColor: colors.status.successSoft,
    borderColor: "#bbf7d0",
  },
  successTitle: {
    ...typography.body,
    color: "#166534",
    fontWeight: "700",
  },
  warningBody: {
    ...typography.bodySmall,
    color: "#92400e",
    marginTop: spacing.xs,
  },
  warningCard: {
    backgroundColor: colors.status.warningSoft,
    borderColor: "#fde68a",
  },
  warningTitle: {
    ...typography.body,
    color: "#78350f",
    fontWeight: "700",
  },
});
