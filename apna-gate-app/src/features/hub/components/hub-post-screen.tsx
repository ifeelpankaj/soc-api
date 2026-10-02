import { useEffect, useMemo, useState } from "react";
import { ScrollView, StyleSheet, Text, View } from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";

import { Stack } from "@/components/layout";
import { ScreenBackHeader } from "@/components/layout/screen-back-header";
import { Button, ConfirmDialog, UserAvatar } from "@/components/ui";
import { HubAttachmentGallery } from "@/features/hub/components/hub-attachment-gallery";
import { HubCommentComposer } from "@/features/hub/components/hub-comment-composer";
import { HubReportDialog } from "@/features/hub/components/hub-report-dialog";
import {
  HubCard,
  HubEmptyState,
  HubMetaRow,
  HubReactionBar,
} from "@/features/hub/components/ui";
import { HubPostSkeleton } from "@/features/hub/components/ui/hub-feed-skeleton";
import { formatHubRelativeTime, hubCategoryName } from "@/features/hub/hub-format";
import {
  residentHubEditPostRoute,
  resolveHubPostBackRoute,
  residentHubRoute,
} from "@/features/hub/hub-routes";
import { hubTheme } from "@/features/hub/hub-theme";
import { useAppFeedback } from "@/features/shared/use-app-feedback";
import { useAuth } from "@/features/auth/use-auth";
import {
  type HubReactionType,
  useHubCategoriesQuery,
  useHubCommentsQuery,
  useHubDeletePostMutation,
  useHubMarkReadMutation,
  useHubPostQuery,
  useHubRemoveReactionMutation,
  useHubReportPostMutation,
  useHubToggleReactionMutation,
} from "@/lib/api/hub-api";
import { spacing } from "@/theme/spacing";
import { colors } from "@/theme/colors";

type HubPostScreenProps = {
  showBackHeader?: boolean;
};

export function HubPostScreen({ showBackHeader = true }: HubPostScreenProps) {
  const router = useRouter();
  const feedback = useAppFeedback();
  const { user } = useAuth();
  const params = useLocalSearchParams<{
    postId?: string;
    societyId?: string;
    channelId?: string;
    returnTo?: string;
  }>();
  const backRoute = resolveHubPostBackRoute(params.returnTo);
  const societyId = Number(params.societyId);
  const postId = Number(params.postId);
  const channelId = params.channelId ? Number(params.channelId) : undefined;
  const valid = Number.isFinite(societyId) && Number.isFinite(postId) && postId > 0;

  const postQuery = useHubPostQuery(
    { societyId, postId },
    { skip: !valid, refetchOnMountOrArgChange: true },
  );
  const commentsQuery = useHubCommentsQuery(
    { societyId, postId },
    { skip: !valid, refetchOnMountOrArgChange: true },
  );
  const categoriesQuery = useHubCategoriesQuery(societyId, { skip: !valid });
  const [markRead] = useHubMarkReadMutation();
  const [toggleReaction, toggleState] = useHubToggleReactionMutation();
  const [removeReaction, removeState] = useHubRemoveReactionMutation();
  const [deletePost, deleteState] = useHubDeletePostMutation();
  const [reportPost, reportState] = useHubReportPostMutation();

  const [deleteOpen, setDeleteOpen] = useState(false);
  const [reportOpen, setReportOpen] = useState(false);

  const post = postQuery.data;
  const resolvedChannelId = channelId ?? post?.channel_id;
  const isAuthor = Boolean(user?.id && post?.author?.id && user.id === post.author.id);
  const reactionBusy = toggleState.isLoading || removeState.isLoading;

  useEffect(() => {
    if (!valid || !resolvedChannelId || !post?.is_important) {
      return;
    }
    void markRead({ societyId, channelId: resolvedChannelId, postId });
  }, [markRead, post?.is_important, postId, resolvedChannelId, societyId, valid]);

  const commentItems = useMemo(
    () => commentsQuery.data?.items ?? [],
    [commentsQuery.data?.items],
  );

  const onReaction = async (type: HubReactionType) => {
    if (!valid) {
      return;
    }
    const active = post?.my_reactions?.includes(type);
    try {
      if (active) {
        await removeReaction({ societyId, postId, reactionType: type }).unwrap();
      } else {
        await toggleReaction({ societyId, postId, reactionType: type }).unwrap();
      }
    } catch (error) {
      feedback.showError("Could not update reaction", error, "Try again.");
    }
  };

  const confirmDelete = async () => {
    try {
      await deletePost({ societyId, postId }).unwrap();
      setDeleteOpen(false);
      feedback.showSuccess("Post deleted", "Your post was removed.");
      router.replace(backRoute);
    } catch (error) {
      feedback.showError("Could not delete", error, "Try again.");
    }
  };

  const submitReport = async (reason: string) => {
    try {
      await reportPost({ societyId, postId, reason }).unwrap();
      setReportOpen(false);
      feedback.showSuccess("Report sent", "Thank you. Admins may review this post.");
    } catch (error) {
      feedback.showError("Could not report", error, "Try again.");
    }
  };

  if (!valid) {
    return (
      <HubEmptyState
        actionLabel="Back to Hub"
        message="This post link is invalid."
        title="Post unavailable"
        onAction={() => router.replace(residentHubRoute())}
      />
    );
  }

  const categoryLabel = hubCategoryName(categoriesQuery.data, post?.category_id);

  return (
    <>
      <ScrollView contentContainerStyle={styles.content} keyboardShouldPersistTaps="handled">
        {showBackHeader ? (
          <ScreenBackHeader fallbackHomeRoute={backRoute} title="Post" />
        ) : null}

        {postQuery.isLoading && !post ? <HubPostSkeleton /> : null}
        {postQuery.isError ? (
          <HubEmptyState
            actionLabel="Try again"
            message="Check your connection and try again."
            title="Unable to load post"
            onAction={() => void postQuery.refetch()}
          />
        ) : null}

        {post ? (
          <Stack gap="lg">
            <HubCard accent={Boolean(post.is_important || post.is_pinned)}>
              {post.is_important || post.is_pinned ? (
                <Text style={styles.badge}>
                  {[post.is_pinned ? "Pinned" : null, post.is_important ? "Important" : null]
                    .filter(Boolean)
                    .join(" · ")}
                </Text>
              ) : null}
              {post.title ? <Text style={styles.title}>{post.title}</Text> : null}
              <Text style={styles.body}>{post.body}</Text>
              <HubMetaRow
                categoryLabel={categoryLabel}
                parts={[
                  post.author?.name,
                  post.author?.flat_label,
                  formatHubRelativeTime(post.created_at),
                ]}
              />
              <HubAttachmentGallery attachmentIds={post.attachment_ids} societyId={societyId} />
              <HubReactionBar
                counts={post.reactions}
                disabled={reactionBusy}
                mine={post.my_reactions}
                onToggle={(type) => void onReaction(type)}
              />
              <View style={styles.postActions}>
                {isAuthor ? (
                  <>
                    <Button
                      compact
                      fullWidth={false}
                      title="Edit post"
                      variant="secondary"
                      onPress={() =>
                        router.push(
                          residentHubEditPostRoute(
                            postId,
                            Array.isArray(params.returnTo)
                              ? (params.returnTo[0] as "hub")
                              : (params.returnTo as "hub" | undefined),
                          ),
                        )
                      }
                    />
                    <Button
                      compact
                      fullWidth={false}
                      title="Delete"
                      variant="danger"
                      onPress={() => setDeleteOpen(true)}
                    />
                  </>
                ) : (
                  <Button
                    compact
                    fullWidth={false}
                    title="Report post"
                    variant="ghost"
                    onPress={() => setReportOpen(true)}
                  />
                )}
              </View>
            </HubCard>

            <Stack gap="md">
              <Text style={styles.sectionTitle}>Comments</Text>
              {commentsQuery.isLoading && commentItems.length === 0 ? (
                <HubPostSkeleton />
              ) : null}
              {commentItems.length === 0 && !commentsQuery.isLoading ? (
                <Text style={styles.emptyComments}>No comments yet. Start the conversation.</Text>
              ) : null}
              {commentItems.map((comment) => (
                <HubCard key={comment.id} style={styles.commentCard}>
                  <View style={styles.commentHeader}>
                    <UserAvatar name={comment.author?.name ?? "R"} size={32} />
                    <Text style={styles.commentAuthor}>
                      {comment.author?.name ?? "Resident"}
                    </Text>
                  </View>
                  <Text style={styles.commentBody}>{comment.body}</Text>
                  <HubAttachmentGallery
                    attachmentIds={comment.attachment_ids}
                    societyId={societyId}
                  />
                  <Text style={styles.commentMeta}>
                    {formatHubRelativeTime(comment.created_at)}
                  </Text>
                </HubCard>
              ))}
            </Stack>

            {post.comments_enabled !== false ? (
              <HubCommentComposer
                postId={postId}
                societyId={societyId}
                onPosted={() => void commentsQuery.refetch()}
              />
            ) : null}
          </Stack>
        ) : null}
      </ScrollView>

      <ConfirmDialog
        confirmLabel="Delete post"
        loading={deleteState.isLoading}
        message="This cannot be undone."
        title="Delete your post?"
        visible={deleteOpen}
        onCancel={() => setDeleteOpen(false)}
        onConfirm={() => void confirmDelete()}
      />
      <HubReportDialog
        loading={reportState.isLoading}
        visible={reportOpen}
        onClose={() => setReportOpen(false)}
        onSubmit={(reason) => void submitReport(reason)}
      />
    </>
  );
}

const styles = StyleSheet.create({
  badge: {
    color: colors.brand.orange,
    fontSize: 12,
    fontWeight: "700",
    marginBottom: spacing.sm,
  },
  body: {
    ...hubTheme.typography.postBody,
    marginTop: spacing.sm,
  },
  commentAuthor: {
    color: colors.brand.navy,
    fontSize: 13,
    fontWeight: "700",
  },
  commentBody: {
    color: colors.text.secondary,
    fontSize: 14,
    lineHeight: 20,
    marginTop: spacing.sm,
  },
  commentCard: {
    padding: spacing.md,
  },
  commentHeader: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.sm,
  },
  commentMeta: {
    ...hubTheme.typography.meta,
    marginTop: spacing.sm,
  },
  content: {
    gap: spacing.lg,
    paddingBottom: spacing["2xl"],
  },
  emptyComments: {
    color: colors.text.muted,
    fontSize: 14,
  },
  postActions: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.sm,
    marginTop: spacing.md,
  },
  sectionTitle: {
    color: colors.brand.navy,
    fontSize: 16,
    fontWeight: "700",
  },
  title: {
    ...hubTheme.typography.postTitle,
  },
});
