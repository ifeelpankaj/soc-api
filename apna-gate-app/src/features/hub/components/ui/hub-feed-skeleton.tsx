import { StyleSheet, View } from "react-native";

import { Stack } from "@/components/layout";
import { hubTheme } from "@/features/hub/hub-theme";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

function Block({ height, width }: { height: number; width?: number | `${number}%` }) {
  return <View style={[styles.block, { height, width: width ?? "100%" }]} />;
}

function RowSkeleton() {
  return (
    <View style={styles.row}>
      <Block height={40} width={40} />
      <Stack gap="sm" style={styles.rowCopy}>
        <Block height={14} width="70%" />
        <Block height={12} width="90%" />
        <Block height={12} width="50%" />
      </Stack>
    </View>
  );
}

export function HubFeedSkeleton({ count = 3 }: { count?: number }) {
  return (
    <Stack gap="md">
      {Array.from({ length: count }, (_, index) => (
        <RowSkeleton key={index} />
      ))}
    </Stack>
  );
}

export function HubPostSkeleton() {
  return (
    <Stack gap="lg">
      <View style={styles.card}>
        <Block height={12} width="30%" />
        <Block height={22} width="85%" />
        <Block height={80} />
        <Block height={12} width="40%" />
      </View>
      <Block height={16} width="25%" />
      <View style={styles.comment}>
        <Block height={12} width="35%" />
        <Block height={40} />
      </View>
    </Stack>
  );
}

const styles = StyleSheet.create({
  block: {
    backgroundColor: colors.surface.secondary,
    borderRadius: radius.sm,
  },
  card: {
    ...hubTheme.card,
    gap: spacing.sm,
  },
  comment: {
    ...hubTheme.card,
    gap: spacing.sm,
    padding: spacing.md,
  },
  row: {
    ...hubTheme.card,
    flexDirection: "row",
    gap: spacing.md,
  },
  rowCopy: {
    flex: 1,
  },
});
