import { View, StyleSheet } from "react-native";

import {
  DashboardActionTile,
  type DashboardActionTileConfig,
} from "@/components/dashboard/dashboard-action-tile";
import { spacing } from "@/theme/spacing";

type DashboardActionGridProps = {
  actions: DashboardActionTileConfig[];
};

export function DashboardActionGrid({ actions }: DashboardActionGridProps) {
  return (
    <View style={styles.grid}>
      {actions.map((action) => (
        <DashboardActionTile
          key={action.id}
          icon={action.icon}
          subtitle={action.subtitle}
          title={action.title}
          tone={action.tone}
          onPress={action.onPress}
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
