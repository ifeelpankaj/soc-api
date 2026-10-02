import { Pressable, ScrollView, StyleSheet, Text, View } from "react-native";

import { SubscriptionExpiredBanner } from "@/components/dashboard";
import { Button, Card } from "@/components/ui";
import { getApiMessage, getVisitorActionErrorMessage, isSubscriptionError } from "@/features/auth/api-error";
import { titleize } from "@/features/guard/guard-utils";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";
import {
  VisitorPurposeSettingCard,
  accessModeToPatch,
} from "@/features/resident/components/visitor-purpose-setting-card";
import { useResidentFeedback } from "@/features/resident/hooks/use-resident-feedback";
import { useResident } from "@/features/resident/resident-context";
import type { ModelsFlatVisitorSettingsResponse } from "@/lib/api/generated-api";
import {
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsQuery,
  usePatchV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsPurposeMutation,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsResetMutation,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

function ErrorBanner({
  message,
  onRetry,
}: {
  message: string;
  onRetry: () => void;
}) {
  return (
    <Pressable accessibilityRole="button" onPress={onRetry} style={styles.errorBanner}>
      <Text style={styles.errorMessage}>{message}</Text>
      <Text style={styles.errorAction}>Retry</Text>
    </Pressable>
  );
}

function societyModeExplanation(mode: string, inherits?: boolean) {
  switch (mode) {
    case "hybrid":
      return "Your society lets each flat choose how guests, delivery, and other visitors are handled.";
    case "mandatory":
      return "Your society requires resident approval for all visitors. These settings are shown for reference.";
    case "optional":
      return "Your society allows visitors by default. These settings are shown for reference.";
    default:
      return inherits
        ? "This flat follows your society's visitor policy."
        : "Contact your society admin if you need to change the visitor policy.";
  }
}

export default function ResidentVisitorSettingsScreen() {
  const { flatId, canManageFlatVisitors, isLoading, requiresSelection, societyId } = useResident();
  const feedback = useResidentFeedback();
  const contextQuery = useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextQuery(
    { societyId: societyId ?? 0, flatId: flatId ?? 0 },
    { skip: !societyId || !flatId },
  );
  const settingsQuery = useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsQuery(
    { societyId: societyId ?? 0, flatId: flatId ?? 0 },
    { skip: !societyId || !flatId || !canManageFlatVisitors },
  );
  const [patchSetting] =
    usePatchV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsPurposeMutation();
  const [resetSettings, resetState] =
    usePostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorSettingsResetMutation();

  const context = contextQuery.data?.data?.context;
  const settings =
    settingsQuery.data?.data?.visitor_settings ?? context?.visitor_settings ?? [];
  const approvalMode = context?.society_approval_mode ?? "mandatory";
  const isHybrid = approvalMode === "hybrid";
  const canEdit = isHybrid && canManageFlatVisitors;
  const isQueryLoading =
    (contextQuery.isLoading && !contextQuery.data) ||
    (canManageFlatVisitors && settingsQuery.isLoading && !settingsQuery.data);
  const failedQuery = [contextQuery, ...(canManageFlatVisitors ? [settingsQuery] : [])].find(
    (query) => query.isError,
  );
  const failedError = failedQuery?.error ?? null;
  const isSubscriptionBlocked = isSubscriptionError(failedError);
  const hasLoadError = !!failedQuery && !isSubscriptionBlocked;

  const refetchAll = () => {
    if (!contextQuery.isUninitialized) {
      void contextQuery.refetch();
    }
    if (canManageFlatVisitors && !settingsQuery.isUninitialized) {
      void settingsQuery.refetch();
    }
  };

  const updateSetting = async (
    purpose: NonNullable<ModelsFlatVisitorSettingsResponse["purpose"]>,
    patch: { approval_required?: boolean; is_enabled?: boolean },
  ) => {
    if (!canEdit || !societyId || !flatId) {
      return;
    }

    try {
      await patchSetting({
        societyId,
        flatId,
        purpose,
        modelsUpdateFlatVisitorSettingRequest: patch,
      }).unwrap();
      refetchAll();
      feedback.showSuccess("Settings updated", "Visitor settings saved.");
    } catch (error) {
      feedback.showError(
        "Update failed",
        error,
        getVisitorActionErrorMessage(error, "Please try again."),
      );
    }
  };

  const handleReset = async () => {
    if (!canEdit || !societyId || !flatId) {
      return;
    }

    try {
      await resetSettings({ societyId, flatId }).unwrap();
      refetchAll();
      feedback.showSuccess("Settings reset", "Flat visitor settings were restored to defaults.");
    } catch (error) {
      feedback.showError(
        "Reset failed",
        error,
        getVisitorActionErrorMessage(error, "Please try again."),
      );
    }
  };

  return (
    <ResidentSubScreen title="Visitor Settings">
      {isLoading || isQueryLoading ? (
        <View style={styles.loadingWrap}>
          <Text style={styles.loadingText}>Loading visitor settings...</Text>
        </View>
      ) : requiresSelection || !societyId || !flatId ? null : (
        <ScrollView contentContainerStyle={styles.scrollContent}>
          <View style={styles.content}>
            <Text style={styles.pageSubtitle}>
              {canEdit
                ? "For each visitor type, pick what should happen when someone arrives at your flat."
                : societyModeExplanation(approvalMode, context?.inherits_society_mode)}
            </Text>

            {isSubscriptionBlocked ? <SubscriptionExpiredBanner /> : null}

            {hasLoadError ? (
              <ErrorBanner
                message={getApiMessage(failedError, "Unable to load visitor settings.")}
                onRetry={refetchAll}
              />
            ) : null}

            {!isSubscriptionBlocked && !canEdit ? (
              <Card style={styles.modeCard}>
                <Text style={styles.modeLabel}>Society policy</Text>
                <Text style={styles.modeValue}>{titleize(approvalMode)}</Text>
                <Text style={styles.modeDescription}>
                  {societyModeExplanation(approvalMode, context?.inherits_society_mode)}
                </Text>
              </Card>
            ) : null}

            {!isSubscriptionBlocked && settings.length > 0 ? (
              <View style={styles.settingsList}>
                {settings.map((setting) => (
                  <VisitorPurposeSettingCard
                    key={`setting-${setting.purpose}`}
                    editable={canEdit}
                    setting={setting}
                    onChange={(mode) => {
                      if (setting.purpose) {
                        void updateSetting(setting.purpose, accessModeToPatch(mode));
                      }
                    }}
                  />
                ))}
              </View>
            ) : !isSubscriptionBlocked ? (
              <Card>
                <Text style={styles.readOnlyBody}>No visitor settings are configured yet.</Text>
              </Card>
            ) : null}

            {!isSubscriptionBlocked && canEdit ? (
              <Button
                title="Reset to defaults"
                variant="secondary"
                loading={resetState.isLoading}
                onPress={handleReset}
              />
            ) : null}

            {!isSubscriptionBlocked && !canEdit && isHybrid ? (
              <Card>
                <Text style={styles.readOnlyTitle}>View only</Text>
                <Text style={styles.readOnlyBody}>
                  Only the flat owner or primary resident can change these settings.
                </Text>
              </Card>
            ) : null}
          </View>
        </ScrollView>
      )}
    </ResidentSubScreen>
  );
}

const styles = StyleSheet.create({
  content: {
    gap: spacing.xl,
  },
  errorAction: {
    color: "#b91c1c",
    fontSize: 13,
    fontWeight: "600",
  },
  errorBanner: {
    alignItems: "center",
    backgroundColor: colors.status.errorSoft,
    borderRadius: radius.xl,
    flexDirection: "row",
    justifyContent: "space-between",
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  errorMessage: {
    color: "#991b1b",
    flex: 1,
    fontSize: 13,
    fontWeight: "500",
    paddingRight: spacing.md,
  },
  loadingText: {
    ...typography.bodySmall,
    color: colors.text.secondary,
    textAlign: "center",
  },
  loadingWrap: {
    flex: 1,
    justifyContent: "center",
    padding: spacing["3xl"],
  },
  modeCard: {
    gap: spacing.sm,
  },
  modeDescription: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  modeLabel: {
    ...typography.eyebrow,
    color: colors.text.muted,
  },
  modeValue: {
    ...typography.subtitle,
    color: colors.text.primary,
    fontWeight: "700",
    textTransform: "capitalize",
  },
  pageSubtitle: {
    ...typography.body,
    color: colors.text.secondary,
  },
  readOnlyBody: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  readOnlyTitle: {
    ...typography.body,
    color: colors.text.primary,
    fontWeight: "700",
  },
  scrollContent: {
    paddingBottom: layout.screenPaddingBottom,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.sm,
  },
  settingsList: {
    gap: spacing.md,
  },
});
