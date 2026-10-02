import { StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import type { SymbolViewProps } from "expo-symbols";

import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type VisitorLogSummaryStripProps = {
  checkedIn: number;
  checkedOut: number;
  pending: number;
  title: string;
};

function SummaryCard({
  icon,
  label,
  tone,
  value,
}: {
  icon: SymbolViewProps["name"];
  label: string;
  tone: "orange" | "green" | "neutral";
  value: number;
}) {
  const toneStyle =
    tone === "green"
      ? {
          bg: colors.dashboard.statGreenSoft,
          border: "#BBF7D0",
          color: colors.dashboard.statGreen,
          iconBg: "#DCFCE7",
        }
      : tone === "orange"
        ? {
            bg: colors.dashboard.statOrangeSoft,
            border: "#FED7AA",
            color: colors.dashboard.statOrange,
            iconBg: "#FFEDD5",
          }
        : {
            bg: colors.dashboard.statNeutralSoft,
            border: colors.guard.border,
            color: colors.dashboard.statNeutral,
            iconBg: "#EEF2F7",
          };

  return (
    <View style={[styles.card, { backgroundColor: toneStyle.bg, borderColor: toneStyle.border }]}>
      <View style={[styles.iconWrap, { backgroundColor: toneStyle.iconBg }]}>
        <SymbolView
          name={icon}
          size={16}
          tintColor={toneStyle.color}
        />
      </View>
      <View style={styles.copy}>
        <Text numberOfLines={1} style={[styles.label, { color: toneStyle.color }]}>
          {label}
        </Text>
        <Text style={[styles.value, { color: toneStyle.color }]}>{value}</Text>
      </View>
    </View>
  );
}

export function VisitorLogSummaryStrip({
  checkedIn,
  checkedOut,
  pending,
  title,
}: VisitorLogSummaryStripProps) {
  return (
    <View style={styles.wrap}>
      <Text style={styles.title}>{title}</Text>
      <View style={styles.row}>
        <SummaryCard
          icon={{ ios: "clock.fill", android: "schedule", web: "schedule" }}
          label="Pending"
          tone="orange"
          value={pending}
        />
        <SummaryCard
          icon={{ ios: "checkmark.circle.fill", android: "check_circle", web: "check_circle" }}
          label="Checked In"
          tone="green"
          value={checkedIn}
        />
        <SummaryCard
          icon={{ ios: "arrow.right.square.fill", android: "logout", web: "logout" }}
          label="Checked Out"
          tone="neutral"
          value={checkedOut}
        />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    alignItems: "center",
    borderRadius: radius.sm,
    borderWidth: 1,
    flex: 1,
    flexDirection: "row",
    gap: spacing.sm,
    minHeight: 52,
    minWidth: 0,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.sm,
  },
  copy: {
    flex: 1,
    justifyContent: "center",
    minWidth: 0,
  },
  iconWrap: {
    alignItems: "center",
    borderRadius: radius.sm,
    height: 34,
    justifyContent: "center",
    width: 34,
  },
  label: {
    fontSize: 10,
    fontWeight: "600",
  },
  row: {
    flexDirection: "row",
    gap: spacing.sm,
  },
  title: {
    color: colors.guard.text,
    fontSize: 13,
    fontWeight: "700",
  },
  value: {
    fontSize: 18,
    fontWeight: "800",
    lineHeight: 22,
  },
  wrap: {
    gap: spacing.sm,
    paddingHorizontal: layout.screenPaddingHorizontal,
  },
});
