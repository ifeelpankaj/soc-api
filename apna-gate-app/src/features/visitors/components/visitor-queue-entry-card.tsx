import {
  ActivityIndicator,
  Pressable,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { SymbolView, type SymbolViewProps } from "expo-symbols";
import { VisitorPhoto } from "@/features/visitors/photos/visitor-photo";

import {
  formatDateOnly,
  formatTimeOfDay,
  getFlatLocationParts,
  getVisitorName,
  getVisitorPhone,
  getVisitorStatusContextMessage,
  getVisitorStatusMeta,
  titleize,
  type WaitingDurationTone,
} from "@/features/visitors/visitor-utils";
import type {
  ModelsVisitorEntry,
  ModelsVisitorPendingEntry,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { shadows } from "@/theme/shadows";

const CARD_BORDER = "rgba(16, 29, 54, 0.11)";

type VisitorQueueEntryCardProps = {
  variant?: "default" | "resident";
  entry: ModelsVisitorEntry | ModelsVisitorPendingEntry;
  disabled?: boolean;
  loadingAction?: "primary" | "secondary" | "tertiary";
  primaryActionLabel?: string;
  secondaryActionLabel?: string;
  tertiaryActionLabel?: string;
  onPress?: () => void;
  onPrimaryAction?: () => void;
  onSecondaryAction?: () => void;
  onTertiaryAction?: () => void;
  waitingLabel?: string;
  waitingTone?: WaitingDurationTone;
};

function getFlatParts(entry: ModelsVisitorEntry | ModelsVisitorPendingEntry) {
  return getFlatLocationParts(entry);
}

function getArrivalParts(
  entry: ModelsVisitorEntry | ModelsVisitorPendingEntry,
) {
  const timestamp =
    entry.checked_in_at ??
    entry.approved_at ??
    entry.created_at ??
    entry.expected_at;

  const label = entry.checked_in_at
    ? "Checked in"
    : entry.approved_at
      ? "Approved"
      : entry.created_at
        ? "Requested"
        : "Expected";

  return {
    date: formatDateOnly(timestamp),
    label,
    time: formatTimeOfDay(timestamp),
  };
}

function waitingToneStyle(tone?: WaitingDurationTone) {
  switch (tone) {
    case "warning":
      return {
        backgroundColor: colors.status.warningSoft,
        borderColor: "#fde68a",
        color: colors.status.warning,
      };
    case "error":
      return {
        backgroundColor: colors.status.errorSoft,
        borderColor: "#fecdd3",
        color: colors.status.error,
      };
    default:
      return {
        backgroundColor: colors.brand.orangeSoft,
        borderColor: "#FED7AA",
        color: colors.brand.orange,
      };
  }
}

function InfoColumn({
  icon,
  label,
  primary,
  secondary,
}: {
  icon: Extract<NonNullable<SymbolViewProps["name"]>, object>;
  label: string;
  primary: string;
  secondary?: string | null;
}) {
  return (
    <View style={styles.infoColumn}>
      <View style={styles.infoIconWrap}>
        <SymbolView
          name={{ ios: icon.ios, android: icon.android, web: icon.web }}
          size={15}
          tintColor={colors.brand.orange}
        />
      </View>
      <View style={styles.infoCopy}>
        <Text style={styles.infoLabel}>{label}</Text>
        <Text numberOfLines={1} style={styles.infoPrimary}>
          {primary}
        </Text>
        {secondary ? (
          <Text numberOfLines={1} style={styles.infoSecondary}>
            {secondary}
          </Text>
        ) : null}
      </View>
    </View>
  );
}

export function VisitorQueueEntryCard({
  variant = "default",
  entry,
  disabled = false,
  loadingAction,
  onPress,
  onPrimaryAction,
  onSecondaryAction,
  onTertiaryAction,
  primaryActionLabel,
  secondaryActionLabel,
  tertiaryActionLabel,
  waitingLabel,
  waitingTone,
}: VisitorQueueEntryCardProps) {
  const visitorName = getVisitorName(entry);
  const isResident = variant === "resident";
  const purpose = entry.purpose ? titleize(entry.purpose) : "Guest";
  const statusMeta = getVisitorStatusMeta(entry.status);
  const contextMessage = getVisitorStatusContextMessage(entry);
  const { number: flatNumber, wing } = getFlatParts(entry);
  const arrival = getArrivalParts(entry);
  const waitStyle = waitingToneStyle(waitingTone);
  const badgeStyle = waitingLabel
    ? waitStyle
    : {
        backgroundColor: statusMeta.bg,
        borderColor: statusMeta.border,
        color: statusMeta.color,
      };
  const badgeLabel = waitingLabel ?? statusMeta.label;
  const hasActions = Boolean(
    primaryActionLabel || secondaryActionLabel || tertiaryActionLabel,
  );

  return (
    <View style={[styles.card, isResident && styles.residentCard]}>
      {isResident ? (
        <View style={styles.residentStatusRow}>
          <View style={[styles.statusDot, { backgroundColor: badgeStyle.color }]} />
          <Text style={[styles.residentStatusText, { color: badgeStyle.color }]}>{badgeLabel}</Text>
        </View>
      ) : null}
      <View style={styles.headerRow}>
        <View style={isResident && styles.residentAvatar}>
          <VisitorPhoto entry={entry} size={isResident ? 60 : 48} />
        </View>

        <View style={styles.headerCopy}>
          <Text style={[styles.visitorName, isResident && styles.residentName]}>{visitorName}</Text>
          <View style={styles.tagRow}>
            <View style={styles.tag}>
              <SymbolView
                name={{ ios: "person.fill", android: "person", web: "person" }}
                size={11}
                tintColor={colors.guard.textMuted}
              />
              <Text style={styles.tagText}>{purpose}</Text>
            </View>
            <View style={styles.tag}>
              <Text numberOfLines={1} style={styles.tagText}>
                {flatNumber}
                {wing ? ` (${wing})` : ""}
              </Text>
            </View>
          </View>
          {entry.service_provider ? (
            <Text style={styles.contextText}>{entry.service_provider}</Text>
          ) : null}
          {getVisitorPhone(entry) ? (
            <Text style={styles.contextText}>{getVisitorPhone(entry)}</Text>
          ) : null}
        </View>

        {!isResident ? <View
          style={[
            styles.statusBadge,
            {
              backgroundColor: badgeStyle.backgroundColor,
              borderColor: badgeStyle.borderColor,
            },
          ]}
        >
          <SymbolView
            name={{ ios: "clock", android: "schedule", web: "schedule" }}
            size={10}
            tintColor={badgeStyle.color}
          />
          <Text style={[styles.statusText, { color: badgeStyle.color }]}>
            {badgeLabel}
          </Text>
        </View> : null}
      </View>

      <View style={styles.divider} />

      <View style={styles.infoRow}>
        <InfoColumn
          icon={{
            ios: "calendar",
            android: "calendar_today",
            web: "calendar_today",
          }}
          label={arrival.label}
          primary={arrival.time || "-"}
          secondary={arrival.date}
        />
        <View style={styles.infoDivider} />
        <InfoColumn
          icon={{
            ios: "building.2.fill",
            android: "apartment",
            web: "apartment",
          }}
          label="Flat"
          primary={flatNumber}
          secondary={wing}
        />
      </View>

      {contextMessage ? (
        <View style={[styles.contextBox, isResident && styles.residentContext]}>
          <View style={styles.contextIcon}>
            <SymbolView
              name={{ ios: "info.circle.fill", android: "info", web: "info" }}
              size={14}
              tintColor={colors.brand.orange}
            />
          </View>
          <Text style={styles.contextText}>{contextMessage}</Text>
        </View>
      ) : null}

      {hasActions ? (
        <View style={styles.actionsWrap}>
          <View style={styles.actionsRow}>
            {secondaryActionLabel ? (
              <Pressable
                accessibilityRole="button"
                disabled={disabled || Boolean(loadingAction)}
                style={({ pressed }) => [
                  styles.secondaryButton,
                  isResident && styles.residentSecondary,
                  disabled && styles.disabled,
                  pressed &&
                    !disabled &&
                    !loadingAction &&
                    styles.buttonPressed,
                ]}
                onPress={onSecondaryAction}
              >
                {loadingAction === "secondary" ? (
                  <ActivityIndicator color={colors.guard.text} size="small" />
                ) : (
                  <Text style={styles.secondaryButtonText}>
                    {secondaryActionLabel}
                  </Text>
                )}
              </Pressable>
            ) : null}
            {primaryActionLabel ? (
              <Pressable
                accessibilityRole="button"
                disabled={disabled || Boolean(loadingAction)}
                style={({ pressed }) => [
                  styles.primaryButton,
                  isResident && styles.residentPrimary,
                  disabled && styles.disabled,
                  !secondaryActionLabel &&
                    !tertiaryActionLabel &&
                    styles.primaryButtonFull,
                  pressed &&
                    !disabled &&
                    !loadingAction &&
                    styles.buttonPressed,
                ]}
                onPress={onPrimaryAction}
              >
                {loadingAction === "primary" ? (
                  <ActivityIndicator color={colors.text.inverse} size="small" />
                ) : (
                  <>
                    <SymbolView
                      name={{
                        ios: "checkmark.circle.fill",
                        android: "check_circle",
                        web: "check_circle",
                      }}
                      size={16}
                      tintColor={colors.text.inverse}
                    />
                    <Text style={styles.primaryButtonText}>
                      {primaryActionLabel}
                    </Text>
                  </>
                )}
              </Pressable>
            ) : null}
          </View>
          {tertiaryActionLabel ? (
            <Pressable
              accessibilityRole="button"
              disabled={disabled || Boolean(loadingAction)}
              style={({ pressed }) => [
                styles.tertiaryButton,
                disabled && styles.disabled,
                pressed && !disabled && !loadingAction && styles.buttonPressed,
              ]}
              onPress={onTertiaryAction}
            >
              {loadingAction === "tertiary" ? (
                <ActivityIndicator color={colors.guard.text} size="small" />
              ) : (
                <>
                  <SymbolView
                    name={{
                      ios: "bell",
                      android: "notifications",
                      web: "notifications",
                    }}
                    size={15}
                    tintColor={colors.guard.text}
                  />
                  <Text style={styles.tertiaryButtonText}>
                    {tertiaryActionLabel}
                  </Text>
                </>
              )}
            </Pressable>
          ) : null}
        </View>
      ) : null}

      {onPress ? (
        <Pressable
          accessibilityRole="button"
          style={styles.detailsLink}
          onPress={onPress}
        >
          <Text style={styles.detailsText}>View details</Text>
          <SymbolView
            name={{
              ios: "chevron.right",
              android: "chevron_right",
              web: "chevron_right",
            }}
            size={12}
            tintColor={colors.guard.textMuted}
          />
        </Pressable>
      ) : null}
    </View>
  );
}

export function VisitorQueueEntryDivider() {
  return <View style={styles.listDivider} />;
}

const styles = StyleSheet.create({
  residentCard: {
    borderWidth: 0,
    borderRadius: 28,
    padding: spacing.xl,
    gap: spacing.lg,
    ...shadows.card,
  },
  residentStatusRow: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
  },
  statusDot: { width: 7, height: 7, borderRadius: 4 },
  residentStatusText: { fontSize: 12, fontWeight: "700", flexShrink: 1 },
  residentAvatar: { borderRadius: 30, overflow: "hidden" },
  residentName: { fontSize: 21, letterSpacing: -0.5 },
  residentContext: {
    borderWidth: 0,
    backgroundColor: "transparent",
    paddingHorizontal: 0,
    paddingVertical: 0,
  },
  residentPrimary: { borderRadius: 999, minHeight: 52, ...shadows.brand },
  residentSecondary: {
    borderRadius: 999,
    minHeight: 52,
    borderWidth: 0,
    backgroundColor: colors.surface.secondary,
  },
  disabled: { opacity: 0.55 },
  actionsRow: {
    flexDirection: "row",
    gap: spacing.sm,
  },
  actionsWrap: {
    gap: spacing.sm,
  },
  avatar: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: 999,
    height: 40,
    justifyContent: "center",
    width: 40,
  },
  avatarText: {
    color: colors.brand.orange,
    fontSize: 17,
    fontWeight: "800",
  },
  buttonPressed: {
    opacity: 0.9,
  },
  card: {
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.lg,
    borderWidth: 1,
    gap: spacing.sm,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingVertical: spacing.md,
  },
  contextBox: {
    alignItems: "flex-start",
    backgroundColor: colors.surface.secondary,
    borderColor: "rgba(16, 29, 54, 0.06)",
    borderRadius: radius.sm,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.sm,
  },
  contextIcon: {
    marginTop: 1,
  },
  contextText: {
    color: colors.guard.textMuted,
    flex: 1,
    fontSize: 12,
    lineHeight: 17,
  },
  detailsLink: {
    minHeight: 44,
    alignItems: "center",
    flexDirection: "row",
    gap: 4,
    justifyContent: "center",
    paddingTop: 2,
  },
  detailsText: {
    color: colors.guard.textMuted,
    fontSize: 12,
    fontWeight: "600",
  },
  divider: {
    backgroundColor: "rgba(16, 29, 54, 0.08)",
    height: 1,
  },
  headerCopy: {
    flex: 1,
    gap: 4,
    minWidth: 0,
  },
  headerRow: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: spacing.sm,
  },
  infoColumn: {
    alignItems: "center",
    flex: 1,
    flexDirection: "row",
    gap: spacing.sm,
    minWidth: 0,
  },
  infoCopy: {
    flex: 1,
    gap: 1,
    minWidth: 0,
  },
  infoDivider: {
    alignSelf: "stretch",
    backgroundColor: "rgba(16, 29, 54, 0.08)",
    marginVertical: 2,
    width: 1,
  },
  infoIconWrap: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.sm,
    height: 30,
    justifyContent: "center",
    width: 30,
  },
  infoLabel: {
    color: colors.guard.textMuted,
    fontSize: 10,
    fontWeight: "600",
  },
  infoPrimary: {
    color: colors.brand.navy,
    fontSize: 13,
    fontWeight: "700",
  },
  infoRow: {
    flexDirection: "row",
    gap: spacing.sm,
  },
  infoSecondary: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "500",
  },
  listDivider: {
    height: spacing.xs,
  },
  primaryButton: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: radius.sm,
    flex: 1,
    flexDirection: "row",
    gap: spacing.xs,
    justifyContent: "center",
    minHeight: 48,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.sm,
  },
  primaryButtonFull: {
    flex: 1,
  },
  primaryButtonText: {
    flexShrink: 1,
    textAlign: "center",
    color: colors.text.inverse,
    fontSize: 12,
    fontWeight: "700",
  },
  secondaryButton: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.sm,
    borderWidth: 1,
    flex: 1,
    flexDirection: "row",
    gap: spacing.xs,
    justifyContent: "center",
    minHeight: 48,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.sm,
  },
  secondaryButtonText: {
    flexShrink: 1,
    textAlign: "center",
    color: colors.guard.text,
    fontSize: 12,
    fontWeight: "600",
  },
  statusBadge: {
    alignItems: "center",
    borderRadius: radius["2xl"],
    borderWidth: 1,
    flexDirection: "row",
    gap: 4,
    paddingHorizontal: 8,
    paddingVertical: 4,
  },
  statusText: {
    fontSize: 10,
    fontWeight: "700",
  },
  tag: {
    alignItems: "center",
    backgroundColor: colors.surface.secondary,
    borderRadius: 999,
    flexDirection: "row",
    gap: 4,
    maxWidth: "100%",
    paddingHorizontal: 8,
    paddingVertical: 3,
  },
  tagRow: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 6,
  },
  tagText: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "500",
  },
  tertiaryButton: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.sm,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.xs,
    justifyContent: "center",
    minHeight: 48,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.sm,
  },
  tertiaryButtonText: {
    flexShrink: 1,
    textAlign: "center",
    color: colors.guard.text,
    fontSize: 12,
    fontWeight: "600",
  },
  visitorName: {
    color: colors.brand.navy,
    fontSize: 16,
    fontWeight: "700",
  },
});
