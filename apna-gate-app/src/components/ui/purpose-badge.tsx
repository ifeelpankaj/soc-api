import { StyleSheet, Text, View } from "react-native";

import type { ModelsVisitorPurpose } from "@/lib/api/generated-api";
import { titleize } from "@/features/guard/guard-utils";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

export function PurposeBadge({ purpose }: { purpose?: ModelsVisitorPurpose | null }) {
  return (
    <View style={styles.wrapper}>
      <Text style={styles.label}>{titleize(purpose)}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  label: {
    ...typography.caption,
    color: colors.text.secondary,
    fontWeight: "700",
  },
  wrapper: {
    alignSelf: "flex-start",
    backgroundColor: colors.guard.sectionBg,
    borderColor: colors.border.default,
    borderRadius: radius["2xl"],
    borderWidth: 1,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.xs,
  },
});
