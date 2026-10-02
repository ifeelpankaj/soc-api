import {
  Modal,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { SymbolView } from "expo-symbols";
import { SafeAreaView } from "react-native-safe-area-context";

import { AppStatusBar } from "@/components/layout/app-status-bar";
import { Row } from "@/components/layout";
import { Button } from "@/components/ui";
import { BottomActions } from "@/components/ui/bottom-actions";
import { VisitorDetailsCard } from "@/features/visitors/components/visitor-details-card";
import type {
  ModelsVisitorEntry,
  ModelsVisitorPendingEntry,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type VisitorDetailSheetProps = {
  variant?: "default" | "resident";
  entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry | null;
  loading?: boolean;
  loadingAction?: "primary" | "secondary";
  onClose: () => void;
  onPrimaryAction?: () => void;
  onSecondaryAction?: () => void;
  primaryActionLabel?: string;
  secondaryActionLabel?: string;
  visible: boolean;
};

export function VisitorDetailSheet({
  variant = "default",
  entry,
  loading,
  loadingAction,
  onClose,
  onPrimaryAction,
  onSecondaryAction,
  primaryActionLabel,
  secondaryActionLabel,
  visible,
}: VisitorDetailSheetProps) {
  if (!entry) {
    return null;
  }

  const onBehalf = entry.metadata?.approved_on_behalf === true;

  return (
    <Modal
      animationType="slide"
      presentationStyle="pageSheet"
      statusBarTranslucent={false}
      visible={visible}
      onRequestClose={onClose}
    >
      <SafeAreaView
        edges={
          primaryActionLabel || secondaryActionLabel
            ? ["top", "left", "right"]
            : ["top", "left", "right", "bottom"]
        }
        style={styles.screen}
      >
        <AppStatusBar />

        <View style={styles.header}>
          <View style={styles.headerSide} />
          <Text style={styles.headerTitle}>Visitor Details</Text>
          <Pressable
            accessibilityLabel="Close"
            hitSlop={12}
            style={styles.headerSide}
            onPress={onClose}
          >
            <SymbolView
              name={{ ios: "xmark", android: "close", web: "close" }}
              size={18}
              tintColor={colors.text.secondary}
            />
          </Pressable>
        </View>

        <ScrollView
          contentContainerStyle={styles.content}
          showsVerticalScrollIndicator={false}
        >
          <VisitorDetailsCard entry={entry} variant={variant === "resident" ? "resident" : "sheet"} />

          {onBehalf ? (
            <View style={styles.auditBanner}>
              <SymbolView
                name={{
                  ios: "person.crop.circle.badge.checkmark",
                  android: "verified_user",
                  web: "verified_user",
                }}
                size={16}
                tintColor={colors.status.warning}
              />
              <Text style={styles.auditLine}>
                Approved on behalf of{" "}
                {String(
                  entry.metadata?.on_behalf_of_resident_name ?? "flat owner",
                )}
              </Text>
            </View>
          ) : null}
        </ScrollView>

        {primaryActionLabel || secondaryActionLabel ? (
          <BottomActions>
            <Row align="center" gap="md">
              {secondaryActionLabel ? (
                <View style={styles.actionSlot}>
                  <Button
                    compact
                    fullWidth
                    disabled={loading}
                    loading={loadingAction === "secondary"}
                    title={secondaryActionLabel}
                    style={variant === "resident" && styles.residentAction}
                    variant="secondary"
                    onPress={onSecondaryAction}
                  />
                </View>
              ) : null}
              {primaryActionLabel ? (
                <View style={styles.actionSlot}>
                  <Button
                    compact
                    fullWidth
                    disabled={loading}
                    loading={loading && loadingAction !== "secondary"}
                    title={primaryActionLabel}
                    style={variant === "resident" && styles.residentAction}
                    onPress={onPrimaryAction}
                  />
                </View>
              ) : null}
            </Row>
          </BottomActions>
        ) : null}
      </SafeAreaView>
    </Modal>
  );
}

const styles = StyleSheet.create({
  residentAction: { borderRadius: 999, minHeight: 52 },
  actionSlot: {
    flex: 1,
  },
  auditBanner: {
    alignItems: "center",
    backgroundColor: colors.status.warningSoft,
    borderRadius: radius.lg,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  auditLine: {
    color: colors.status.warning,
    flex: 1,
    fontSize: 13,
    fontWeight: "600",
    lineHeight: 18,
  },
  content: {
    gap: spacing.lg,
    paddingBottom: spacing["3xl"],
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.md,
  },
  footer: {
    backgroundColor: colors.surface.card,
    borderTopColor: colors.border.default,
    borderTopWidth: StyleSheet.hairlineWidth,
    paddingBottom: spacing.lg,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.lg,
    ...shadows.sm,
  },
  header: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderBottomColor: colors.border.default,
    borderBottomWidth: StyleSheet.hairlineWidth,
    flexDirection: "row",
    justifyContent: "space-between",
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  headerSide: {
    alignItems: "center",
    height: 32,
    justifyContent: "center",
    width: 32,
  },
  headerTitle: {
    ...typography.subtitle,
    color: colors.text.primary,
    fontWeight: "700",
  },
  screen: {
    backgroundColor: colors.guard.screenBg,
    flex: 1,
  },
});
