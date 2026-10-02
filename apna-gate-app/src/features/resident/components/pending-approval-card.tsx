import { SymbolView } from "expo-symbols";
import { VisitorPhoto } from "@/features/visitors/photos/visitor-photo";
import { StyleSheet, Text, View } from "react-native";

import { Button, Card } from "@/components/ui";
import {
  formatDateTime,
  getFlatLabel,
  getVisitorName,
  getVisitorPhone,
} from "@/features/visitors/visitor-utils";
import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

type PendingApprovalCardProps = {
  entry: ModelsVisitorEntry;
  disabled?: boolean;
  loadingAction?: "approve" | "reject";
  onApprove?: () => void;
  onReject?: () => void;
};

function remainingLabel(value?: string) {
  if (!value) return "-";
  const minutes = Math.max(
    0,
    Math.round((new Date(value).getTime() - Date.now()) / 60000),
  );
  if (minutes < 60) return `in ${minutes}m`;
  return `in ${Math.floor(minutes / 60)}h ${minutes % 60}m`;
}

export function PendingApprovalCard({
  entry,
  disabled,
  loadingAction,
  onApprove,
  onReject,
}: PendingApprovalCardProps) {
  const name = getVisitorName(entry);
  const phone = getVisitorPhone(entry);

  return (
    <Card style={styles.card}>
      <View style={styles.topRow}>
        <VisitorPhoto entry={entry} size={68} />
        <View style={styles.identity}>
          <Text numberOfLines={1} style={styles.name}>
            {name}
          </Text>
          <View style={styles.metaRow}>
            <SymbolView
              name={{ ios: "person", android: "person", web: "person" }}
              size={16}
              tintColor={colors.text.secondary}
            />
            <Text style={styles.meta}>Guest</Text>
            <Text style={styles.dot}>|</Text>
            <Text style={styles.meta}>{getFlatLabel(entry)}</Text>
          </View>
          <View style={styles.metaRow}>
            <SymbolView
              name={{
                ios: "calendar",
                android: "calendar_month",
                web: "calendar_month",
              }}
              size={16}
              tintColor={colors.text.secondary}
            />
            <Text style={styles.meta}>
              {formatDateTime(entry.created_at || entry.expected_at)}
            </Text>
          </View>
          {phone ? (
            <View style={styles.metaRow}>
              <SymbolView
                name={{ ios: "phone", android: "phone", web: "phone" }}
                size={16}
                tintColor={colors.text.secondary}
              />
              <Text style={styles.meta}>{phone}</Text>
            </View>
          ) : null}
        </View>
        <View style={styles.pendingBadge}>
          <Text style={styles.pendingText}>Pending</Text>
        </View>
      </View>

      <View style={styles.reasonBox}>
        <SymbolView
          name={{
            ios: "hourglass",
            android: "hourglass_empty",
            web: "hourglass_empty",
          }}
          size={22}
          tintColor={colors.brand.orange}
        />
        <View style={styles.reasonCopy}>
          <Text style={styles.reasonTitle}>Waiting for resident approval</Text>
          <Text style={styles.reasonText}>
            {entry.notes || "The visitor can enter after approval."}
          </Text>
        </View>
      </View>

      <View style={styles.checkoutBox}>
        <View>
          <Text style={styles.checkoutLabel}>EXPECTED CHECKOUT</Text>
          <Text style={styles.checkoutValue}>
            {formatDateTime(entry.expected_checkout_at)}
          </Text>
        </View>
        <View style={styles.divider} />
        <View style={styles.remaining}>
          <SymbolView
            name={{ ios: "clock", android: "schedule", web: "schedule" }}
            size={20}
            tintColor={colors.text.secondary}
          />
          <Text style={styles.meta}>
            {remainingLabel(entry.expected_checkout_at)}
          </Text>
        </View>
      </View>

      <View style={styles.actions}>
        <View style={styles.actionSlot}>
          <Button
            compact
            fullWidth
            disabled={disabled}
            loading={loadingAction === "reject"}
            title="Reject"
            variant="secondary"
            onPress={onReject}
          />
        </View>
        <View style={styles.actionSlot}>
          <Button
            compact
            fullWidth
            disabled={disabled}
            loading={loadingAction === "approve"}
            title="Approve"
            onPress={onApprove}
          />
        </View>
      </View>
    </Card>
  );
}

const styles = StyleSheet.create({
  actionSlot: { flex: 1, minWidth: 0 },
  actions: { flexDirection: "row", gap: spacing.md },
  avatar: {
    alignItems: "center",
    backgroundColor: "#fff0e8",
    borderRadius: radius.full,
    height: 68,
    justifyContent: "center",
    width: 68,
  },
  avatarText: { color: colors.brand.orange, fontSize: 34, fontWeight: "800" },
  card: { gap: spacing.lg, padding: spacing.lg, ...shadows.card },
  checkoutBox: {
    alignItems: "center",
    borderColor: colors.border.default,
    borderRadius: radius.md,
    borderWidth: 1,
    flexDirection: "row",
    justifyContent: "space-between",
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  checkoutLabel: {
    color: colors.text.secondary,
    fontSize: 11,
    fontWeight: "700",
  },
  checkoutValue: {
    color: colors.text.primary,
    fontSize: 16,
    fontWeight: "700",
    marginTop: 4,
  },
  divider: { backgroundColor: colors.border.default, height: 42, width: 1 },
  dot: { color: colors.text.secondary, marginHorizontal: 2 },
  identity: { flex: 1, gap: spacing.xs, minWidth: 0 },
  meta: { color: colors.text.secondary, fontSize: 14 },
  metaRow: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.sm,
    minWidth: 0,
  },
  name: { color: colors.text.primary, fontSize: 20, fontWeight: "800" },
  pendingBadge: {
    backgroundColor: "#fffaf0",
    borderColor: "#f6d37a",
    borderRadius: radius.md,
    borderWidth: 1,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  pendingText: { color: "#b45309", fontSize: 14, fontWeight: "600" },
  reasonBox: {
    alignItems: "flex-start",
    backgroundColor: "#fff7f1",
    borderRadius: radius.md,
    flexDirection: "row",
    gap: spacing.md,
    padding: spacing.md,
  },
  reasonCopy: { flex: 1, gap: 2 },
  reasonText: { color: colors.text.secondary, fontSize: 14, lineHeight: 20 },
  reasonTitle: { color: colors.text.primary, fontSize: 15, fontWeight: "700" },
  remaining: { alignItems: "center", flexDirection: "row", gap: spacing.sm },
  topRow: { alignItems: "flex-start", flexDirection: "row", gap: spacing.md },
});
