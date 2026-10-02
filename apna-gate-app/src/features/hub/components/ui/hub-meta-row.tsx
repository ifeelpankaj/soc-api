import { StyleSheet, Text, View } from "react-native";

import { hubTheme } from "@/features/hub/hub-theme";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type HubMetaRowProps = {
  parts: (string | null | undefined)[];
  categoryLabel?: string;
};

export function HubMetaRow({ parts, categoryLabel }: HubMetaRowProps) {
  const line = parts.filter(Boolean).join(" · ");
  return (
    <View style={styles.row}>
      {line ? <Text style={styles.meta}>{line}</Text> : null}
      {categoryLabel ? (
        <View style={styles.badge}>
          <Text style={styles.badgeText}>{categoryLabel}</Text>
        </View>
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  badge: {
    backgroundColor: colors.brand.navySoft,
    borderRadius: radius.sm,
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
  },
  badgeText: {
    color: colors.brand.navy,
    fontSize: 11,
    fontWeight: "600",
  },
  meta: {
    ...hubTheme.typography.meta,
    flex: 1,
  },
  row: {
    alignItems: "center",
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.sm,
    marginTop: spacing.md,
  },
});
