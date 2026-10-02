import { Pressable, StyleSheet, Text, View } from "react-native";
import { useRouter } from "expo-router";
import { SymbolView } from "expo-symbols";

import { DashboardSection } from "@/components/dashboard";
import { Row } from "@/components/layout";
import { useResident } from "@/features/resident/resident-context";
import {
  residentMaintenanceBillRoute,
  residentMaintenancePaymentRoute,
  residentMaintenanceRoute,
} from "@/features/resident/resident-routes";
import {
  formatBillingMonth,
  formatMoney,
  getMaintenanceBillPresentation,
} from "@/features/resident/maintenance/maintenance-presentation";
import { formatDateOnly } from "@/features/guard/guard-utils";
import type { MaintenanceBill } from "@/lib/api/maintenance-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

type ResidentMaintenanceDashboardCardProps = {
  bill?: MaintenanceBill | null;
  loading?: boolean;
};

export function ResidentMaintenanceDashboardCard({
  bill,
  loading,
}: ResidentMaintenanceDashboardCardProps) {
  const router = useRouter();
  const { societyId, flatId } = useResident();

  if (loading || !bill || !societyId || !flatId) {
    return null;
  }

  const presentation = getMaintenanceBillPresentation({ bill });
  const amount = formatMoney(bill.total_paise);
  const month = formatBillingMonth(bill.billing_month);

  let headline = `${amount} · ${month}`;
  let detail = bill.due_message ?? `Due ${formatDateOnly(bill.due_date, bill.timezone)}`;
  let actionLabel: string = presentation.actionLabel;
  let onPress = () => router.push(residentMaintenanceBillRoute(bill.id, societyId, flatId));

  if (presentation.state === "PENDING_REVIEW") {
    detail = `${amount} · UTR submitted`;
    actionLabel = "View Details";
  } else if (presentation.state === "PAID") {
    headline = `${amount} · Paid`;
    detail = bill.paid_on
      ? `Paid ${formatDateOnly(bill.paid_on, bill.timezone)}`
      : "Payment verified";
    actionLabel = "View Receipt";
  } else if (presentation.state === "REJECTED") {
    headline = "Payment verification failed";
    detail = `${amount} · ${month}`;
    actionLabel = "Review";
  } else if (presentation.action === "pay") {
    onPress = () => router.push(residentMaintenancePaymentRoute(bill.id, societyId, flatId));
    actionLabel = "Pay Now";
  }

  return (
    <DashboardSection
      actionLabel="Open →"
      title="Maintenance"
      onAction={() => router.navigate(residentMaintenanceRoute())}
    >
      <Pressable
        accessibilityRole="button"
        style={({ pressed }) => [styles.card, pressed && styles.pressed]}
        onPress={onPress}
      >
        <Row align="center" justify="space-between">
          <View style={styles.copy}>
            <Text style={styles.headline}>{headline}</Text>
            <Text style={styles.detail}>{detail}</Text>
          </View>
          <Row align="center" gap="xs">
            <Text style={styles.action}>{actionLabel}</Text>
            <SymbolView
              name={{ ios: "chevron.right", android: "chevron_right", web: "chevron_right" }}
              size={14}
              tintColor={colors.brand.orange}
            />
          </Row>
        </Row>
      </Pressable>
    </DashboardSection>
  );
}

const styles = StyleSheet.create({
  action: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "700",
  },
  card: {
    backgroundColor: colors.surface.card,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.lg,
    borderWidth: StyleSheet.hairlineWidth,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.lg,
    ...shadows.sm,
  },
  copy: {
    flex: 1,
    paddingRight: spacing.md,
  },
  detail: {
    color: colors.text.muted,
    fontSize: 13,
    marginTop: 4,
  },
  headline: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "600",
  },
  pressed: {
    opacity: 0.92,
  },
});
