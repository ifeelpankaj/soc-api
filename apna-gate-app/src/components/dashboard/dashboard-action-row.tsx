import { View, StyleSheet } from "react-native";

import {
  DashboardActionTile,
  type DashboardActionTileConfig,
} from "@/components/dashboard/dashboard-action-tile";
import { spacing } from "@/theme/spacing";

type DashboardActionRowProps = {
  actions: DashboardActionTileConfig[];
  columns?: number;
};

export function DashboardActionRow({ actions, columns }: DashboardActionRowProps) {
  if (columns) {
    return (
      <View style={{ gap: spacing.md }}>
        {Array.from({ length: Math.ceil(actions.length / columns) }, (_, row) => (
          <DashboardActionRow key={row} actions={actions.slice(row * columns, (row + 1) * columns)} />
        ))}
      </View>
    );
  }
  return (
    <View style={styles.row}>
      {actions.map((action) => (
        <View key={action.id} style={styles.tileSlot}>
          <DashboardActionTile
            badgeCount={action.badgeCount}
            icon={action.icon}
            subtitle={action.subtitle}
            title={action.title}
            tone={action.tone}
            onPress={action.onPress}
          />
        </View>
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  row: {
    flexDirection: "row",
    gap: spacing.md,
  },
  tileSlot: {
    flex: 1,
    minWidth: 0,
  },
});
