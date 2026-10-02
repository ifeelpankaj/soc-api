import { Pressable, StyleSheet, View } from "react-native";
import { Stack, Row } from "@/components/layout";
import { AppText, Badge } from "@/components/ui";
import { colors } from "@/theme/colors";
import { formatDateOnly } from "@/features/guard/guard-utils";
import type { MaintenanceBill } from "@/lib/api/maintenance-api";
import {
  formatBillingMonth,
  formatMoney,
  getMaintenanceBillPresentation,
} from "./maintenance-presentation";
import { spacing } from "@/theme/spacing";

export function BillListItem({
  bill,
  onPress,
}: {
  bill: MaintenanceBill;
  onPress: () => void;
}) {
  const status = getMaintenanceBillPresentation({ bill });
  return (
    <Pressable
      accessibilityRole="button"
      accessibilityLabel={`${formatBillingMonth(bill.billing_month)}, ${status.label}`}
      onPress={onPress}
      style={styles.billRow}
    >
      <Stack gap="sm">
        <Row justify="space-between" style={{ gap: spacing.sm, alignItems: "flex-start" }}>
          <AppText style={[styles.bold, { flex: 1, fontSize: 16 }]}>{formatBillingMonth(bill.billing_month)}</AppText>
          <View style={{ maxWidth: "48%" }}><Badge label={status.label} tone={status.tone} /></View>
        </Row>
        <Row justify="space-between">
          <AppText style={[styles.bold, { fontSize: 22, lineHeight: 28 }]}>{formatMoney(bill.total_paise)}</AppText>
          <AppText color="muted" accessibilityElementsHidden>{"\u203a"}</AppText>
        </Row>
        <AppText variant="caption" color="muted">
          {status.state === "PAID" && bill.paid_on
            ? `Paid ${formatDateOnly(bill.paid_on, "UTC")}`
            : `Due ${formatDateOnly(bill.due_date, "UTC")}`}
        </AppText>
      </Stack>
    </Pressable>
  );
}
export const maintenanceStyles = StyleSheet.create({
  content: { padding: spacing.lg, gap: spacing.lg },
  card: { gap: spacing.md },
  bold: { fontWeight: "700" },
  billRow: {
    backgroundColor: colors.surface.card,
    padding: spacing.lg,
    borderRadius: 18,
    borderWidth: 1,
    borderColor: colors.border.default,
  },
});
const styles = maintenanceStyles;
