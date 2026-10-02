import { Pressable, StyleSheet, Text } from "react-native";
import { useRouter } from "expo-router";

import { DashboardSection } from "@/components/dashboard";
import { Stack } from "@/components/layout";
import { residentHubAnnouncementsRoute } from "@/features/hub/hub-routes";
import { formatHubRelativeTime, hubPostPreview } from "@/features/hub/hub-format";
import type { HubContent } from "@/lib/api/hub-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

type ResidentAnnouncementsSnapshotProps = {
  items: HubContent[];
  loading?: boolean;
  onItemPress: (postId: number) => void;
};

export function ResidentAnnouncementsSnapshot({
  items,
  loading,
  onItemPress,
}: ResidentAnnouncementsSnapshotProps) {
  const router = useRouter();

  if (loading || items.length === 0) {
    return null;
  }

  const preview = items.slice(0, 2);
  const sectionAction = preview.length > 1 ? "View all →" : "Open Hub →";

  return (
    <DashboardSection
      actionLabel={sectionAction}
      title="Society Announcements"
      onAction={() => router.push(residentHubAnnouncementsRoute())}
    >
      <Stack gap="sm">
        {preview.map((post) => {
          const postId = post.post_id ?? post.id;
          return (
            <Pressable
              key={postId}
              accessibilityRole="button"
              style={({ pressed }) => [styles.card, pressed && styles.pressed]}
              onPress={() => onItemPress(postId)}
            >
              <Text numberOfLines={2} style={styles.title}>
                {post.title?.trim() || "Society announcement"}
              </Text>
              <Text numberOfLines={2} style={styles.body}>
                {hubPostPreview(post.body)}
              </Text>
              <Text style={styles.meta}>
                {[post.is_important ? "Important" : null, formatHubRelativeTime(post.created_at)]
                  .filter(Boolean)
                  .join(" · ")}
              </Text>
            </Pressable>
          );
        })}
      </Stack>
    </DashboardSection>
  );
}

const styles = StyleSheet.create({
  body: {
    color: colors.text.secondary,
    fontSize: 13,
    lineHeight: 18,
    marginTop: 4,
  },
  card: {
    backgroundColor: colors.surface.card,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.lg,
    borderWidth: StyleSheet.hairlineWidth,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    ...shadows.sm,
  },
  meta: {
    color: colors.text.muted,
    fontSize: 12,
    marginTop: spacing.sm,
  },
  pressed: {
    opacity: 0.92,
  },
  title: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "600",
  },
});
