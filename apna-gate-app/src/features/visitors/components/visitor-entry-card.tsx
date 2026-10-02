import { Pressable, StyleSheet, Text, View } from "react-native";
import { VisitorPhoto } from "@/features/visitors/photos/visitor-photo";

import { Row, Stack } from "@/components/layout";
import { Button, Card } from "@/components/ui";
import {
  formatDateTime,
  getFlatLabel,
  getVisitorName,
  getVisitorStatusContextMessage,
  getVisitorStatusMeta,
  titleize,
  type WaitingDurationTone,
} from "@/features/visitors/visitor-utils";
import { VisitorTimelineRow } from "@/features/visitors/components/visitor-timeline-row";
import type {
  ModelsVisitorEntry,
  ModelsVisitorPendingEntry,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type VisitorEntryCardProps = {
  entry: ModelsVisitorEntry | ModelsVisitorPendingEntry;
  primaryActionLabel?: string;
  secondaryActionLabel?: string;
  onPrimaryAction?: () => void;
  onSecondaryAction?: () => void;
  onPress?: () => void;
  loading?: boolean;
  loadingEntryId?: number;
  waitingLabel?: string;
  waitingTone?: WaitingDurationTone;
};

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
        backgroundColor: "#ecfdf5",
        borderColor: "#a7f3d0",
        color: colors.status.success,
      };
  }
}

export function VisitorEntryCard({
  entry,
  loading,
  loadingEntryId,
  onPress,
  onPrimaryAction,
  onSecondaryAction,
  primaryActionLabel,
  secondaryActionLabel,
  waitingLabel,
  waitingTone,
}: VisitorEntryCardProps) {
  const isLoading = loading || loadingEntryId === entry.id;
  const statusMeta = getVisitorStatusMeta(entry.status);
  const waitStyle = waitingToneStyle(waitingTone);
  const hasActions = Boolean(primaryActionLabel || secondaryActionLabel);
  const visitorName = getVisitorName(entry);
  const contextMessage = getVisitorStatusContextMessage(entry);
  const phone = entry.visitor?.phone_number;
  const email = entry.visitor?.email;

  const header = (
    <View style={styles.header}>
      <VisitorPhoto entry={entry} />

      <Stack gap="xs" style={styles.copy}>
        <Row align="flex-start" gap="sm" justify="space-between">
          <View style={styles.titleBlock}>
            <Text style={styles.visitorName}>{visitorName}</Text>
            <Text style={styles.visitDetail}>
              {entry.purpose ? titleize(entry.purpose) : "Visitor"} · {getFlatLabel(entry)}
            </Text>
          </View>
          <View style={styles.badges}>
            {waitingLabel ? (
              <View
                style={[
                  styles.statusBadge,
                  {
                    backgroundColor: waitStyle.backgroundColor,
                    borderColor: waitStyle.borderColor,
                  },
                ]}
              >
                <Text style={[styles.statusText, { color: waitStyle.color }]}>{waitingLabel}</Text>
              </View>
            ) : null}
            <View
              style={[
                styles.statusBadge,
                { backgroundColor: statusMeta.bg, borderColor: statusMeta.border },
              ]}
            >
              <Text style={[styles.statusText, { color: statusMeta.color }]}>{statusMeta.label}</Text>
            </View>
          </View>
        </Row>

        <Text style={styles.timestamp}>
          {entry.checked_in_at
            ? `Checked in ${formatDateTime(entry.checked_in_at)}`
            : entry.expected_at
              ? `Expected ${formatDateTime(entry.expected_at)}`
              : entry.approved_at
                ? `Approved ${formatDateTime(entry.approved_at)}`
                : entry.created_at
                  ? `Created ${formatDateTime(entry.created_at)}`
                  : "Gate record"}
        </Text>

        {phone || email || entry.vehicle_number ? (
          <View style={styles.quickFacts}>
            {phone ? <Text style={styles.quickFact}>Mobile · {phone}</Text> : null}
            {email ? <Text style={styles.quickFact}>Email · {email}</Text> : null}
            {entry.vehicle_number ? (
              <Text style={styles.quickFact}>Vehicle · {entry.vehicle_number}</Text>
            ) : null}
          </View>
        ) : null}

        {contextMessage &&
        (entry.status === "expired" ||
          entry.status === "rejected" ||
          entry.status === "waiting_approval") ? (
          <Text style={styles.contextLine}>{contextMessage}</Text>
        ) : null}

        <VisitorTimelineRow entry={entry} />
      </Stack>
    </View>
  );

  const body = (
    <>
      {onPress && hasActions ? (
        <Pressable accessibilityRole="button" onPress={onPress}>
          {header}
          {entry.notes ? <Text style={styles.notes}>{entry.notes}</Text> : null}
        </Pressable>
      ) : (
        <>
          {header}
          {entry.notes ? <Text style={styles.notes}>{entry.notes}</Text> : null}
        </>
      )}
    </>
  );

  const actions = hasActions ? (
    <Row align="center" gap="md" justify="flex-start">
      {secondaryActionLabel ? (
        <View style={styles.actionSlot}>
          <Button
            compact
            fullWidth
            title={secondaryActionLabel}
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
            loading={isLoading}
            title={primaryActionLabel}
            onPress={onPrimaryAction}
          />
        </View>
      ) : null}
    </Row>
  ) : null;

  const card = (
    <Card style={styles.card}>
      {body}
      {actions}
    </Card>
  );

  if (onPress && !hasActions) {
    return (
      <Pressable accessibilityRole="button" onPress={onPress}>
        {card}
      </Pressable>
    );
  }

  return card;
}

const styles = StyleSheet.create({
  actionSlot: {
    flex: 1,
    minWidth: 0,
  },
  avatar: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderRadius: radius.full,
    height: 52,
    justifyContent: "center",
    width: 52,
  },
  avatarRing: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.full,
    padding: 3,
  },
  avatarText: {
    color: colors.brand.orange,
    fontSize: 22,
    fontWeight: "800",
  },
  badges: {
    alignItems: "flex-end",
    gap: spacing.xs,
  },
  card: {
    gap: spacing.lg,
    padding: spacing.xl,
    ...shadows.card,
  },
  contextLine: {
    ...typography.bodySmall,
    color: colors.text.secondary,
    lineHeight: 20,
  },
  copy: {
    flex: 1,
    minWidth: 0,
  },
  header: {
    flexDirection: "row",
    gap: spacing.lg,
  },
  notes: {
    ...typography.bodySmall,
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.lg,
    color: colors.text.primary,
    lineHeight: 20,
    marginTop: spacing.md,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  quickFact: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  quickFacts: {
    gap: 2,
  },
  statusBadge: {
    borderRadius: radius["2xl"],
    borderWidth: 1,
    paddingHorizontal: spacing.md,
    paddingVertical: 4,
  },
  statusText: {
    fontSize: 11,
    fontWeight: "700",
  },
  timestamp: {
    ...typography.caption,
    color: colors.text.muted,
  },
  titleBlock: {
    flex: 1,
    gap: 2,
    minWidth: 0,
  },
  visitDetail: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  visitorName: {
    ...typography.subtitle,
    color: colors.text.primary,
    fontSize: 20,
    fontWeight: "800",
  },
});
