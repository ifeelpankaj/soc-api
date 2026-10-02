import { StyleSheet, Text, View } from "react-native";

import { hubTheme } from "@/features/hub/hub-theme";
import { spacing } from "@/theme/spacing";

type HubPageHeaderProps = {
  subtitle?: string;
};

export function HubPageHeader({ subtitle }: HubPageHeaderProps) {
  if (!subtitle) {
    return null;
  }
  return (
    <View style={styles.wrap}>
      <Text style={styles.subtitle}>{subtitle}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  subtitle: {
    ...hubTheme.typography.pageSubtitle,
    marginTop: -spacing.sm,
  },
  wrap: {
    marginBottom: spacing.xs,
  },
});
