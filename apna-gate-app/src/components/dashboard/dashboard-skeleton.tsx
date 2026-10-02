import { StyleSheet, View, type ViewStyle } from "react-native";

import { Row, Stack } from "@/components/layout";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type DashboardSkeletonSection = "header" | "hero" | "actions" | "stats" | "activity";

type DashboardSkeletonProps = {
  sections?: DashboardSkeletonSection[];
};

function Skeleton({ style }: { style?: ViewStyle }) {
  return <View style={[styles.skeleton, style]} />;
}

export function DashboardSkeleton({
  sections = ["header", "hero", "actions", "stats", "activity"],
}: DashboardSkeletonProps) {
  return (
    <Stack gap="2xl">
      {sections.includes("header") ? (
        <Stack gap="lg">
          <Row justify="space-between">
            <Skeleton style={styles.skeletonBrand} />
            <Row gap="sm">
              <Skeleton style={styles.skeletonIconButton} />
              <Skeleton style={styles.skeletonAvatar} />
            </Row>
          </Row>
          <Skeleton style={styles.skeletonGreeting} />
          <Skeleton style={styles.skeletonLocation} />
          <Skeleton style={styles.skeletonAnnouncement} />
        </Stack>
      ) : null}

      {sections.includes("hero") || sections.includes("actions") ? (
        <Stack gap="md">
          {sections.includes("hero") ? <Skeleton style={styles.skeletonHero} /> : null}
          {sections.includes("actions") ? (
            <Row gap="md">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} style={styles.skeletonActionTile} />
              ))}
            </Row>
          ) : null}
        </Stack>
      ) : null}

      {sections.includes("stats") ? (
        <Stack gap="lg">
          <Skeleton style={styles.skeletonSectionTitle} />
          <Stack gap="md">
            {[1, 2].map((row) => (
              <Row key={row} gap="md">
                {[1, 2].map((col) => (
                  <Skeleton key={`${row}-${col}`} style={styles.skeletonStatTile} />
                ))}
              </Row>
            ))}
          </Stack>
        </Stack>
      ) : null}

      {sections.includes("activity") ? (
        <Stack gap="md">
          <Skeleton style={styles.skeletonSectionTitle} />
          <View style={styles.skeletonActivityPanel}>
            <Skeleton style={styles.skeletonActivityIllustration} />
            <View style={styles.skeletonActivityCopy}>
              <Skeleton style={styles.skeletonActivityTitle} />
              <Skeleton style={styles.skeletonActivityText} />
              <Skeleton style={styles.skeletonActivityButton} />
            </View>
          </View>
        </Stack>
      ) : null}
    </Stack>
  );
}

const styles = StyleSheet.create({
  skeleton: {
    backgroundColor: "rgba(226, 232, 240, 0.7)",
    borderRadius: radius.sm,
  },
  skeletonActionTile: {
    borderRadius: radius.xl,
    flex: 1,
    height: layout.actionTileMinHeight,
  },
  skeletonActivityButton: {
    borderRadius: radius.lg,
    height: 44,
    marginTop: spacing.sm,
    width: "100%",
  },
  skeletonActivityCopy: {
    flex: 1,
    gap: spacing.sm,
  },
  skeletonActivityIllustration: {
    borderRadius: radius.lg,
    height: 72,
    width: 64,
  },
  skeletonActivityPanel: {
    backgroundColor: colors.surface.card,
    borderRadius: radius.xl,
    flexDirection: "row",
    gap: spacing.lg,
    padding: spacing.lg,
  },
  skeletonActivityText: {
    height: 14,
    width: "90%",
  },
  skeletonActivityTitle: {
    height: 18,
    width: "70%",
  },
  skeletonAnnouncement: {
    borderRadius: radius.xl,
    height: 64,
    width: "100%",
  },
  skeletonAvatar: {
    borderRadius: 20,
    height: 40,
    width: 40,
  },
  skeletonBrand: {
    borderRadius: radius.sm,
    height: 32,
    width: 140,
  },
  skeletonGreeting: {
    height: 28,
    width: 240,
  },
  skeletonHero: {
    borderRadius: radius["2xl"],
    height: 88,
    width: "100%",
  },
  skeletonIconButton: {
    borderRadius: 20,
    height: 40,
    width: 40,
  },
  skeletonLocation: {
    height: 14,
    width: 180,
  },
  skeletonSectionTitle: {
    height: 14,
    width: 112,
  },
  skeletonStatTile: {
    borderRadius: radius.xl,
    flex: 1,
    height: 96,
  },
});
