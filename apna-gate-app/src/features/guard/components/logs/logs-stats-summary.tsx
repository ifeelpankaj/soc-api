import { StyleSheet, Text, View } from "react-native";

import { theme } from "@/theme";

const G = theme.guard;

type LogsStatsSummaryProps = {
  checkedOutToday: number;
  pendingCount: number;
  visitorsInside: number;
};

export function LogsStatsSummary({
  checkedOutToday,
  pendingCount,
  visitorsInside,
}: LogsStatsSummaryProps) {
  const items = [
    {
      count: pendingCount,
      label: "Pending",
      tone: {
        bg: theme.status.warningSoft,
        border: "#fde68a",
        count: theme.status.warning,
        text: "#92400e",
      },
    },
    {
      count: visitorsInside,
      label: "Checked In",
      tone: {
        bg: theme.status.successSoft,
        border: "#bbf7d0",
        count: theme.status.success,
        text: "#166534",
      },
    },
    {
      count: checkedOutToday,
      label: "Checked Out",
      tone: {
        bg: "#f8fafc",
        border: theme.border.default,
        count: G.textMuted,
        text: G.textMuted,
      },
    },
  ];

  return (
    <View style={styles.row}>
      {items.map((item) => (
        <View
          key={item.label}
          style={[
            styles.card,
            { backgroundColor: item.tone.bg, borderColor: item.tone.border },
          ]}
        >
          <Text style={[styles.label, { color: item.tone.text }]}>{item.label}</Text>
          <Text style={[styles.count, { color: item.tone.count }]}>{item.count}</Text>
        </View>
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    alignItems: "center",
    borderRadius: 12,
    borderWidth: 1,
    flex: 1,
    gap: 2,
    minWidth: 0,
    paddingHorizontal: 8,
    paddingVertical: 10,
  },
  count: {
    fontSize: 18,
    fontWeight: "700",
  },
  label: {
    fontSize: 11,
    fontWeight: "600",
    textAlign: "center",
  },
  row: {
    flexDirection: "row",
    gap: 8,
  },
});
