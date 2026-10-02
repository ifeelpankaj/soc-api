import { StyleSheet, View } from "react-native";

import { Stack } from "@/components/layout";
import { spacing } from "@/theme/spacing";

import {
  DashboardOverviewStat,
  type DashboardOverviewStatConfig,
} from "./dashboard-overview-stat";

type DashboardOverviewGridProps = {
  onStatPress?: (id: string) => void;
  stats: DashboardOverviewStatConfig[];
};

export function DashboardOverviewGrid({ onStatPress, stats }: DashboardOverviewGridProps) {
  const rows = [stats.slice(0, 2), stats.slice(2, 4)];

  return (
    <Stack gap="md">
      {rows.map((row, rowIndex) => (
        <View key={`overview-row-${rowIndex}`} style={styles.row}>
          {row.map((stat) => (
            <View key={stat.id} style={styles.cell}>
              <DashboardOverviewStat
                icon={stat.icon}
                label={stat.label}
                subtext={stat.subtext}
                tone={stat.tone}
                value={stat.value}
                onPress={onStatPress ? () => onStatPress(stat.id) : undefined}
              />
            </View>
          ))}
        </View>
      ))}
    </Stack>
  );
}

const styles = StyleSheet.create({
  cell: {
    alignSelf: "stretch",
    flex: 1,
    minWidth: 0,
  },
  row: {
    alignItems: "stretch",
    flexDirection: "row",
    gap: spacing.md,
  },
});
