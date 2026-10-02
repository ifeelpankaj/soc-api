import { useState } from "react";
import { ActivityIndicator, StyleSheet, Text, View } from "react-native";

import { Stack } from "@/components/layout";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";
import { ResidentScreenShell } from "@/features/resident/components/resident-screen-shell";
import { useResident } from "@/features/resident/resident-context";
import { useResidentAnnouncementChannelsQuery, useResidentAnnouncementsQuery } from "@/lib/api/announcement-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

export default function AnnouncementsScreen() {
  const { societyId } = useResident();
  return <ResidentScreenShell>{societyId ? <AnnouncementFeed key={societyId} societyId={societyId} /> : null}</ResidentScreenShell>;
}

function AnnouncementFeed({ societyId }: { societyId: number }) {
  const [cursors, setCursors] = useState<(string | undefined)[]>([undefined]);
  const channels = useResidentAnnouncementChannelsQuery(societyId, { refetchOnMountOrArgChange: true });
  const channel = channels.currentData?.data?.find((item) => item.type === "announcement");
  const posts = useResidentAnnouncementsQuery(
    { societyId, channelId: channel?.id ?? 0, cursor: cursors[cursors.length - 1] },
    { skip: !channel, refetchOnMountOrArgChange: true },
  );
  const page = posts.currentData?.data;
  const items = Array.from(new Map([...(page?.pinned ?? []), ...(page?.items ?? [])].map((post) => [post.id, post])).values());
  const loading = channels.isFetching || (Boolean(channel) && posts.isFetching);
  const error = channels.isError || (Boolean(channel) && posts.isError);

  return (
    <Stack gap="lg">
      <Text style={styles.heading}>Announcements</Text>
      <Text style={styles.subtitle}>News and notices from your society</Text>
      {loading ? <ActivityIndicator color={colors.brand.orange} /> : error ? (
        <EmptyState title="Unable to load announcements" message="Please try again." actionLabel="Retry" onAction={() => { void channels.refetch(); if (channel) void posts.refetch(); }} />
      ) : items.length === 0 ? (
        <EmptyState title="No announcements yet" message="Society announcements will appear here." />
      ) : items.map((post) => (
        <View key={post.id} style={styles.card}>
          {post.is_pinned || post.is_important ? <Text style={styles.badge}>{[post.is_pinned ? "Pinned" : "", post.is_important ? "Important" : ""].filter(Boolean).join(" · ")}</Text> : null}
          {post.title ? <Text style={styles.title}>{post.title}</Text> : null}
          <Text style={styles.body}>{post.body}</Text>
          <Text style={styles.subtitle}>{post.author?.name ?? "Society"} · {new Date(post.created_at).toLocaleDateString()}</Text>
        </View>
      ))}
      {cursors.length > 1 ? <Button title="Previous" variant="secondary" disabled={loading} onPress={() => setCursors((current) => current.slice(0, -1))} /> : null}
      {page?.next_cursor ? <Button title="Next announcements" variant="secondary" disabled={loading} onPress={() => setCursors((current) => [...current, page.next_cursor])} /> : null}
      <Button title="Refresh" variant="secondary" disabled={loading} onPress={() => { void channels.refetch(); if (channel) void posts.refetch(); }} />
    </Stack>
  );
}

const styles = StyleSheet.create({
  heading: { color: colors.brand.navy, fontSize: 24, fontWeight: "700" },
  subtitle: { color: colors.text.muted, fontSize: 12 },
  card: { backgroundColor: colors.surface.card, borderRadius: radius.lg, padding: spacing.lg, gap: spacing.sm },
  title: { color: colors.brand.navy, fontSize: 17, fontWeight: "700" },
  body: { color: colors.brand.navy, fontSize: 14, lineHeight: 22 },
  badge: { color: colors.brand.orange, fontSize: 12, fontWeight: "700" },
});
