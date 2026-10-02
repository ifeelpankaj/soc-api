import { StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import { VisitorPhoto } from "@/features/visitors/photos/visitor-photo";
import { VisitorPhotoPreview } from "@/features/visitors/photos/visitor-photo-viewer";

import { Stack } from "@/components/layout";
import { Card } from "@/components/ui";
import {
  formatDateTime,
  getCompanionDetails,
  getFlatLabel,
  getVisitorDetailRows,
  getVisitorName,
  getVisitorStatusContextMessage,
  getVisitorStatusMeta,
  isVisitorCheckoutOverdue,
  titleize,
} from "@/features/guard/guard-utils";
import { VisitorTimelineRow } from "@/features/visitors/components/visitor-timeline-row";
import type { ModelsVisitorEntry, ModelsVisitorPendingEntry } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type VisitorDetailsCardProps = {
  entry: ModelsVisitorEntry | ModelsVisitorPendingEntry;
  showTimeline?: boolean;
  variant?: "compact" | "default" | "sheet" | "resident";
};

function DetailListRow({
  compact,
  resident,
  label,
  value,
}: {
  compact?: boolean;
  resident?: boolean;
  label: string;
  value: string;
}) {
  return (
    <View style={[styles.detailRow, compact && styles.detailRowCompact, resident && styles.residentDetailRow]}>
      <Text style={[styles.detailRowLabel, resident && styles.residentDetailLabel]}>{label}</Text>
      <Text style={styles.detailRowValue}>{value}</Text>
    </View>
  );
}

function CompanionRow({
  companion,
  index,
}: {
  companion: { name: string; phoneNumber: string };
  index: number;
}) {
  const title = companion.name || `Companion ${index + 1}`;
  const phone = companion.phoneNumber;

  return (
    <View style={styles.companionRow}>
      <View style={styles.companionAvatar}>
        <Text style={styles.companionAvatarText}>
          {title.charAt(0).toUpperCase()}
        </Text>
      </View>
      <View style={styles.companionCopy}>
        <Text numberOfLines={1} style={styles.companionName}>
          {title}
        </Text>
        {phone ? (
          <Text numberOfLines={1} style={styles.companionPhone}>
            {phone}
          </Text>
        ) : null}
      </View>
    </View>
  );
}

export function VisitorDetailsCard({
  entry,
  showTimeline = true,
  variant = "default",
}: VisitorDetailsCardProps) {
  const statusMeta = getVisitorStatusMeta(entry.status);
  const detailRows = getVisitorDetailRows(entry);
  const companionDetails = getCompanionDetails(entry);
  const contextMessage = getVisitorStatusContextMessage(entry);
  const checkoutOverdue = isVisitorCheckoutOverdue(entry);
  const visitorName = getVisitorName(entry);
  const purposeLabel = entry.purpose ? titleize(entry.purpose) : "Visitor";
  const isCompact = variant === "compact";
  const isResident = variant === "resident";

  return (
    <Card
      style={[
        styles.card,
        variant === "sheet" && styles.cardSheet,
        isCompact && styles.cardCompact,
        isResident && styles.residentPage,
      ]}
    >
      <View style={[styles.hero, isCompact && styles.heroCompact]}>
        {!isCompact && entry.visitor?.photo_url?.trim() ? (
          <View style={[styles.photoBanner, isResident && styles.residentPortrait]}>
            <VisitorPhotoPreview entry={entry} />
          </View>
        ) : (
          <View style={isResident && styles.residentAvatar}>
            <VisitorPhoto entry={entry} variant="list" size={isResident ? 88 : undefined} />
          </View>
        )}

        <Stack align="center" gap="xs" style={styles.heroCopy}>
          <Text
            style={[styles.name, isCompact && styles.nameCompact]}
          >
            {visitorName}
          </Text>
          <View
            style={[
              styles.statusBadge,
              { backgroundColor: statusMeta.bg, borderColor: statusMeta.border },
            ]}
          >
            <Text style={[styles.statusText, { color: statusMeta.color }]}>{statusMeta.label}</Text>
          </View>
          {checkoutOverdue ? (
            <View style={[styles.statusBadge, styles.overdueBadge]}>
              <Text style={styles.overdueText}>Overdue checkout</Text>
            </View>
          ) : null}
          <Text style={styles.meta}>
            {purposeLabel} · {getFlatLabel(entry)}
          </Text>
          {entry.created_at ? (
            <Text style={styles.subtle}>Created {formatDateTime(entry.created_at)}</Text>
          ) : null}
        </Stack>
      </View>

      {contextMessage ? (
        <View
          style={[
            styles.contextBanner,
            isResident && styles.residentContext,
            entry.status === "expired" || entry.status === "rejected"
              ? styles.contextBannerMuted
              : checkoutOverdue
                ? styles.contextBannerWarning
                : styles.contextBannerInfo,
          ]}
        >
          <SymbolView
            name={{
              ios: "info.circle.fill",
              android: "info",
              web: "info",
            }}
            size={16}
            tintColor={
              entry.status === "expired" || entry.status === "rejected"
                ? colors.status.error
                : checkoutOverdue
                  ? colors.status.warning
                  : colors.guard.teal
            }
          />
          <Text style={styles.contextText}>{contextMessage}</Text>
        </View>
      ) : null}

      {detailRows.length > 0 ? (
        <View style={styles.section}>
          <Text style={[styles.sectionTitle, isResident && styles.residentSectionTitle]}>Visitor details</Text>
          <View style={[styles.detailList, isResident && styles.openList]}>
            {detailRows.map((row) => (
              <DetailListRow
                compact={isCompact}
                resident={isResident}
                key={row.label}
                label={row.label}
                value={row.value}
              />
            ))}
          </View>
        </View>
      ) : null}

      {companionDetails.length > 0 ? (
        <View style={styles.section}>
          <Text style={[styles.sectionTitle, isResident && styles.residentSectionTitle]}>Companion details</Text>
          <View style={[styles.companionList, isResident && styles.openList]}>
            {companionDetails.map((companion, index) => (
              <CompanionRow
                key={`${companion.name}-${companion.phoneNumber}-${index}`}
                companion={companion}
                index={index}
              />
            ))}
          </View>
        </View>
      ) : null}

      {showTimeline ? (
        <View style={styles.section}>
          <Text style={[styles.sectionTitle, isResident && styles.residentSectionTitle]}>Visit timeline</Text>
          <VisitorTimelineRow entry={entry} />
        </View>
      ) : null}

      {entry.notes ? (
        <View style={[styles.notesBox, isResident && styles.residentNotes]}>
          <Text style={styles.notesLabel}>Notes</Text>
          <Text style={styles.notesValue}>{entry.notes}</Text>
        </View>
      ) : null}
    </Card>
  );
}

const styles = StyleSheet.create({
  residentDetailRow: { paddingHorizontal: 0, paddingVertical: spacing.lg, gap: spacing.sm },
  residentDetailLabel: { textTransform: "none", letterSpacing: 0, fontSize: 13, fontWeight: "500" },
  residentSectionTitle: { textTransform: "none", letterSpacing: -0.3, fontSize: 18, fontWeight: "700", color: colors.text.primary },
  residentContext: { borderWidth: 0, borderRadius: 24 },
  residentNotes: { backgroundColor: "transparent", borderLeftColor: colors.brand.orange, borderLeftWidth: 3, borderRadius: 0, paddingVertical: spacing.xs },
  residentPage: {
    backgroundColor: "transparent",
    borderWidth: 0,
    borderRadius: 0,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xl,
    gap: spacing["3xl"],
    elevation: 0,
    shadowOpacity: 0,
    boxShadow: "none",
  },
  residentPortrait: {
    width: 156,
    height: 180,
    borderRadius: 48,
    borderWidth: 0,
  },
  residentAvatar: {
    borderRadius: 999,
    overflow: "hidden",
    borderWidth: 6,
    borderColor: colors.surface.card,
  },
  openList: {
    backgroundColor: "transparent",
    borderWidth: 0,
    borderRadius: 0,
    padding: 0,
  },
  avatar: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderRadius: radius.full,
    height: 72,
    justifyContent: "center",
    width: 72,
  },
  avatarCompact: {
    height: 52,
    width: 52,
  },
  avatarRing: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.full,
    padding: 4,
  },
  avatarText: {
    color: colors.brand.orange,
    fontSize: 28,
    fontWeight: "800",
  },
  avatarTextCompact: {
    fontSize: 21,
  },
  card: {
    gap: spacing["2xl"],
    padding: spacing["2xl"],
    ...shadows.card,
  },
  cardCompact: {
    gap: spacing.lg,
    padding: spacing.lg,
  },
  cardSheet: {
    padding: spacing.xl,
  },
  companionAvatar: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.full,
    height: 34,
    justifyContent: "center",
    width: 34,
  },
  companionAvatarText: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "800",
  },
  companionCopy: {
    flex: 1,
    gap: 2,
    minWidth: 0,
  },
  companionList: {
    backgroundColor: colors.surface.secondary,
    borderColor: colors.border.default,
    borderRadius: radius.xl,
    borderWidth: StyleSheet.hairlineWidth,
    gap: spacing.xs,
    padding: spacing.sm,
  },
  companionName: {
    ...typography.bodySmall,
    color: colors.text.primary,
    fontWeight: "700",
  },
  companionPhone: {
    ...typography.caption,
    color: colors.text.secondary,
  },
  companionRow: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: StyleSheet.hairlineWidth,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  contextBanner: {
    alignItems: "flex-start",
    borderRadius: radius.xl,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  contextBannerInfo: {
    backgroundColor: colors.guard.tealSoft,
    borderColor: "#99f6e4",
  },
  contextBannerMuted: {
    backgroundColor: colors.status.errorSoft,
    borderColor: "#fecaca",
  },
  contextBannerWarning: {
    backgroundColor: colors.status.warningSoft,
    borderColor: "#fde68a",
  },
  contextText: {
    ...typography.bodySmall,
    color: colors.text.secondary,
    flex: 1,
    lineHeight: 20,
  },
  detailList: {
    backgroundColor: colors.surface.secondary,
    borderColor: colors.border.default,
    borderRadius: radius.xl,
    borderWidth: StyleSheet.hairlineWidth,
    overflow: "hidden",
  },
  detailRow: {
    borderBottomColor: colors.border.default,
    borderBottomWidth: StyleSheet.hairlineWidth,
    gap: 4,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  detailRowCompact: {
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  detailRowLabel: {
    ...typography.caption,
    color: colors.text.muted,
    fontSize: 11,
    fontWeight: "700",
    letterSpacing: 0.8,
    textTransform: "uppercase",
  },
  detailRowValue: {
    ...typography.body,
    color: colors.text.primary,
    fontWeight: "600",
  },
  hero: {
    alignItems: "center",
    gap: spacing.lg,
  },
  photoBanner: {
    width: "100%",
    height: 240,
    backgroundColor: colors.surface.secondary,
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: 1,
    overflow: "hidden",
  },
  heroCompact: {
    gap: spacing.sm,
  },
  heroCopy: {
    alignItems: "center",
  },
  meta: {
    ...typography.bodySmall,
    color: colors.text.secondary,
    textAlign: "center",
  },
  name: {
    ...typography.title,
    color: colors.text.primary,
    fontSize: 26,
    fontWeight: "800",
    textAlign: "center",
  },
  nameCompact: {
    fontSize: 22,
  },
  notesBox: {
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.xl,
    gap: spacing.xs,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.lg,
  },
  notesLabel: {
    color: colors.brand.orange,
    fontSize: 12,
    fontWeight: "700",
    letterSpacing: 0.6,
    textTransform: "uppercase",
  },
  notesValue: {
    ...typography.body,
    color: colors.text.primary,
    lineHeight: 22,
  },
  overdueBadge: {
    backgroundColor: colors.status.warningSoft,
    borderColor: "#fde68a",
  },
  overdueText: {
    color: colors.status.warning,
    fontSize: 12,
    fontWeight: "700",
    letterSpacing: 0.3,
  },
  section: {
    gap: spacing.md,
  },
  sectionTitle: {
    ...typography.eyebrow,
    color: colors.text.muted,
    fontSize: 11,
  },
  statusBadge: {
    borderRadius: radius["2xl"],
    borderWidth: 1,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.xs,
  },
  statusText: {
    fontSize: 12,
    fontWeight: "700",
    letterSpacing: 0.3,
  },
  subtle: {
    ...typography.caption,
    color: colors.text.muted,
    textAlign: "center",
  },
});
