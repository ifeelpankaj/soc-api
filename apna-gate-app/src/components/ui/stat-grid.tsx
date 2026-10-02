import { StyleSheet, View } from "react-native";

import { spacing } from "@/theme/spacing";

import { StatCard } from "./stat-card";

type StatItem = {
  label: string;
  value: string | number;
  tone?: "default" | "teal" | "amber" | "rose";
};

const toneMap: Record<NonNullable<StatItem["tone"]>, "default" | "teal" | "warning" | "success"> = {
  amber: "warning",
  default: "default",
  rose: "default",
  teal: "teal",
};

export function StatGrid({ stats }: { stats: StatItem[] }) {
  return (
    <View style={styles.grid}>
      {stats.map((stat) => (
        <StatCard
          key={stat.label}
          label={stat.label}
          tone={toneMap[stat.tone ?? "default"]}
          value={stat.value}
        />
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  grid: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.md,
  },
});
