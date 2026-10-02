import { StyleSheet, View } from "react-native";
import { AppText } from "./app-text";
import { spacing } from "@/theme/spacing";
import { colors } from "@/theme/colors";

export function InfoRow({ label, value }: { label: string; value: string }) {
  return (
    <View style={styles.row}>
      <AppText variant="bodySmall" color="muted" style={styles.label}>
        {label}
      </AppText>
      <AppText variant="bodySmall" selectable style={styles.value}>
        {value}
      </AppText>
    </View>
  );
}
const styles = StyleSheet.create({
  row: {
    flexDirection: "row",
    gap: spacing.md,
    justifyContent: "space-between",
    paddingVertical: spacing.sm,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.border.default,
  },
  label: { flex: 1 },
  value: { flex: 1, textAlign: "right", fontWeight: "600" },
});
