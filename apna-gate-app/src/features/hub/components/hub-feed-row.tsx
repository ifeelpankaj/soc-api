import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";

import { Row } from "@/components/layout";
import { hubTheme } from "@/features/hub/hub-theme";
import { formatHubRelativeTime, hubPostPreview } from "@/features/hub/hub-format";
import { reactionTotal, type HubContent } from "@/lib/api/hub-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type HubFeedRowProps = {
  post: HubContent;
  variant?: "announcement" | "community";
  categoryLabel?: string;
  onPress: () => void;
};

export function HubFeedRow({
  post,
  variant = "announcement",
  categoryLabel,
  onPress,
}: HubFeedRowProps) {
  const likes = reactionTotal(post.reactions);
  const comments = post.reply_count ?? 0;
  const attachmentCount = post.attachment_ids?.length ?? 0;
  const authorLine =
    variant === "community"
      ? [post.author?.name, post.author?.flat_label].filter(Boolean).join(" · ")
      : null;

  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [
        styles.row,
        (post.is_important || post.is_pinned) && styles.rowAccent,
        pressed && styles.pressed,
      ]}
      onPress={onPress}
    >
      <Row align="flex-start" gap="sm">
        <SymbolView
          name={
            variant === "community"
              ? { ios: "person.circle", android: "account_circle", web: "account_circle" }
              : { ios: "megaphone", android: "campaign", web: "campaign" }
          }
          size={18}
          tintColor={
            variant === "community"
              ? colors.dashboard.actionBlue
              : colors.dashboard.actionPurple
          }
        />
        <View style={styles.copy}>
          {authorLine ? <Text style={styles.author}>{authorLine}</Text> : null}
          <Text numberOfLines={2} style={styles.title}>
            {post.title?.trim() || (variant === "community" ? "Community post" : "Announcement")}
          </Text>
          <Text numberOfLines={2} style={styles.body}>
            {hubPostPreview(post.body)}
          </Text>
          <Row align="center" gap="md" style={styles.metaRow}>
            <Text style={styles.meta}>
              {[
                post.is_important ? "Important" : post.is_pinned ? "Pinned" : null,
                categoryLabel,
                formatHubRelativeTime(post.created_at),
              ]
                .filter(Boolean)
                .join(" · ")}
            </Text>
          </Row>
          <Row align="center" gap="md" style={styles.statsRow}>
            <Text style={styles.stat}>{likes} likes</Text>
            <Text style={styles.stat}>{comments} comments</Text>
            {attachmentCount > 0 ? (
              <Row align="center" gap="sm">
                <SymbolView
                  name={{ ios: "photo", android: "image", web: "image" }}
                  size={14}
                  tintColor={colors.text.muted}
                />
                <Text style={styles.stat}>{attachmentCount}</Text>
              </Row>
            ) : null}
          </Row>
        </View>
      </Row>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  author: {
    color: colors.text.secondary,
    fontSize: 12,
    fontWeight: "600",
    marginBottom: 2,
  },
  body: {
    color: colors.text.secondary,
    fontSize: 13,
    lineHeight: 18,
    marginTop: 4,
  },
  copy: {
    flex: 1,
    minWidth: 0,
  },
  meta: {
    color: colors.text.muted,
    flex: 1,
    fontSize: 12,
  },
  metaRow: {
    marginTop: spacing.sm,
  },
  pressed: {
    opacity: hubTheme.pressOpacity,
  },
  row: {
    ...hubTheme.card,
  },
  rowAccent: {
    borderLeftColor: colors.brand.orange,
    borderLeftWidth: 3,
  },
  stat: {
    color: colors.text.muted,
    fontSize: 12,
  },
  statsRow: {
    marginTop: spacing.xs,
  },
  title: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "600",
  },
});
