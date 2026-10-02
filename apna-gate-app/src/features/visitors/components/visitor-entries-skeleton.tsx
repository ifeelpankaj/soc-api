import { StyleSheet, View } from "react-native";

import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

export function VisitorEntriesSkeleton() {
  return (
    <View
      accessible
      accessibilityLabel="Loading visitor entries"
      accessibilityState={{ busy: true }}
    >
      <View
        accessibilityElementsHidden
        importantForAccessibility="no-hide-descendants"
        style={styles.list}
      >
        {[0, 1, 2, 3, 4].map((index) => (
          <View key={index} style={styles.card}>
            <View style={styles.titleRow}>
              <View style={[styles.placeholder, styles.avatar]} />
              <View style={styles.copy}>
                <View style={[styles.placeholder, styles.name]} />
                <View style={[styles.placeholder, styles.meta]} />
              </View>
              <View style={[styles.placeholder, styles.badge]} />
            </View>
            <View style={styles.timeline}>
              {[0, 1, 2].map((step) => (
                <View key={step} style={styles.step}>
                  <View style={[styles.placeholder, styles.timeLabel]} />
                  <View style={[styles.placeholder, styles.date]} />
                  <View style={[styles.placeholder, styles.time]} />
                </View>
              ))}
            </View>
          </View>
        ))}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  list: { gap: spacing.xs * 2 },
  card: {
    backgroundColor: colors.surface.card,
    borderColor: "rgba(16, 29, 54, 0.11)",
    borderRadius: radius.sm,
    borderWidth: 1,
    gap: spacing.sm,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingVertical: spacing.md,
  },
  placeholder: {
    backgroundColor: "rgba(226, 232, 240, 0.7)",
    borderRadius: radius.sm,
  },
  titleRow: { flexDirection: "row", alignItems: "center", gap: spacing.sm },
  avatar: { width: 36, height: 36 },
  copy: { flex: 1, gap: spacing.xs },
  name: { width: "70%", height: 15 },
  meta: { width: "90%", height: 12 },
  badge: { width: 64, height: 20 },
  timeline: {
    backgroundColor: colors.surface.secondary,
    borderColor: "rgba(16, 29, 54, 0.08)",
    borderRadius: radius.sm,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
  },
  step: { flex: 1, alignItems: "center", gap: spacing.xs, paddingVertical: spacing.sm },
  timeLabel: { width: "60%", height: 10 },
  date: { width: "80%", height: 13 },
  time: { width: "80%", height: 12 },
});
