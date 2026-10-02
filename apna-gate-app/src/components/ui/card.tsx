import { StyleSheet, View, type ViewProps } from "react-native";

import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

export function Card({ style, ...props }: ViewProps) {
  return <View style={[styles.card, style]} {...props} />;
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: 1,
    padding: spacing.xl,
    ...shadows.card,
  },
});
