import { useEffect, useMemo, useRef, useState } from "react";
import { ScrollView, StyleSheet, Text, View } from "react-native";
import { useRouter } from "expo-router";

import { Stack } from "@/components/layout";
import { Button } from "@/components/ui";
import { ScreenBackHeader } from "@/components/layout/screen-back-header";
import {
  HubAddAttachmentTile,
  HubAttachmentTile,
  HubChip,
  HubChipWrap,
  HubPageHeader,
  HubTextField,
} from "@/features/hub/components/ui";
import { HubPostSkeleton } from "@/features/hub/components/ui/hub-feed-skeleton";
import { hubTheme } from "@/features/hub/hub-theme";
import {
  defaultCategoryId,
  HUB_MAX_ATTACHMENTS,
  hubUploadFormData,
  normalizeAttachmentIds,
  pickHubPhotos,
} from "@/features/hub/hub-upload";
import { hubPublishErrorMessage, hubUploadErrorMessage } from "@/features/hub/hub-api-errors";
import { residentHubCommunityRoute, residentHubPostRoute } from "@/features/hub/hub-routes";
import { useAppFeedback } from "@/features/shared/use-app-feedback";
import { useResident } from "@/features/resident/resident-context";
import {
  findHubChannel,
  useHubCategoriesQuery,
  useHubChannelsQuery,
  useHubCreatePostMutation,
  useHubDeleteUploadMutation,
  useHubPostQuery,
  useHubUpdatePostMutation,
  useHubAttachmentQuery,
  useHubUploadAttachmentMutation,
} from "@/lib/api/hub-api";
import { spacing } from "@/theme/spacing";

type AttachmentSlot = {
  key: string;
  localUri: string;
  uploadId?: number;
  uploading?: boolean;
  error?: string;
  persisted?: boolean;
};

type HubCreatePostScreenProps = {
  showBackHeader?: boolean;
  editPostId?: number;
};

export function HubCreatePostScreen({
  showBackHeader = true,
  editPostId,
}: HubCreatePostScreenProps) {
  const router = useRouter();
  const feedback = useAppFeedback();
  const { societyId } = useResident();
  const isEdit = Boolean(editPostId && editPostId > 0);

  const channelsQuery = useHubChannelsQuery(societyId ?? 0, { skip: !societyId });
  const communityChannel = findHubChannel(channelsQuery.data, "community");
  const categoriesQuery = useHubCategoriesQuery(societyId ?? 0, { skip: !societyId });
  const postQuery = useHubPostQuery(
    { societyId: societyId ?? 0, postId: editPostId ?? 0 },
    { skip: !societyId || !isEdit },
  );

  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [categoryId, setCategoryId] = useState<number | undefined>();
  const [attachments, setAttachments] = useState<AttachmentSlot[]>([]);
  const [bodyError, setBodyError] = useState<string | undefined>();

  const [createPost, createState] = useHubCreatePostMutation();
  const [updatePost, updateState] = useHubUpdatePostMutation();
  const [uploadAttachment] = useHubUploadAttachmentMutation();
  const [deleteUpload] = useHubDeleteUploadMutation();

  const categories = categoriesQuery.data ?? [];
  const publishing = createState.isLoading || updateState.isLoading;
  const uploadsBusy = attachments.some((slot) => slot.uploading);

  const editInitialized = useRef(false);
  useEffect(() => {
    if (!isEdit || !postQuery.data || editInitialized.current) {
      return;
    }
    editInitialized.current = true;
    const post = postQuery.data;
    setTitle(post.title ?? "");
    setBody(post.body);
    setCategoryId(post.category_id ?? defaultCategoryId(categories));
    setAttachments(
      (post.attachment_ids ?? []).map((uploadId) => ({
        key: `existing-${uploadId}`,
        localUri: "",
        uploadId,
        persisted: true,
      })),
    );
  }, [categories, isEdit, postQuery.data]);

  useEffect(() => {
    if (isEdit || categoryId !== undefined || !categories.length) {
      return;
    }
    setCategoryId(defaultCategoryId(categories));
  }, [categories, categoryId, isEdit]);

  const attachmentCount = attachments.length;
  const canAddMore = attachmentCount < HUB_MAX_ATTACHMENTS;

  const addPhotos = async () => {
    if (!societyId || !canAddMore) {
      return;
    }
    try {
      const picked = await pickHubPhotos(HUB_MAX_ATTACHMENTS - attachmentCount);
      for (const photo of picked) {
        const key = `local-${Date.now()}-${Math.random()}`;
        setAttachments((current) => [
          ...current,
          { key, localUri: photo.uri, uploading: true },
        ]);
        try {
          const form = await hubUploadFormData(photo);
          const uploaded = await uploadAttachment({ societyId, body: form }).unwrap();
          setAttachments((current) =>
            current.map((slot) =>
              slot.key === key
                ? { ...slot, uploadId: uploaded.id, uploading: false }
                : slot,
            ),
          );
        } catch (error) {
          setAttachments((current) =>
            current.map((slot) =>
              slot.key === key
                ? {
                    ...slot,
                    uploading: false,
                    error: hubUploadErrorMessage(error),
                  }
                : slot,
            ),
          );
        }
      }
    } catch (error) {
      feedback.showInfo("Photo not added", hubUploadErrorMessage(error));
    }
  };

  const removeAttachment = (slot: AttachmentSlot) => {
    setAttachments((current) => current.filter((item) => item.key !== slot.key));
    if (societyId && slot.uploadId && !slot.persisted) {
      void deleteUpload({ societyId, uploadId: slot.uploadId });
    }
  };

  const submit = async () => {
    if (!societyId || !communityChannel?.id) {
      feedback.showInfo("Unavailable", "Community channel is not ready yet.");
      return;
    }
    const trimmed = body.trim();
    if (trimmed.length < 4) {
      setBodyError("Write a short message for your neighbours.");
      return;
    }
    if (trimmed.length > 20_000) {
      setBodyError("Message is too long.");
      return;
    }
    setBodyError(undefined);
    if (uploadsBusy) {
      feedback.showInfo("Please wait", "Photos are still uploading.");
      return;
    }
    const ids = normalizeAttachmentIds(
      attachments
        .filter((slot) => slot.uploadId && !slot.error)
        .map((slot) => slot.uploadId!),
    );
    const payload = {
      title: title.trim() || undefined,
      body: trimmed,
      category_id: categoryId,
      attachment_ids: ids.length ? ids : undefined,
    };
    try {
      if (isEdit && editPostId) {
        await updatePost({
          societyId,
          postId: editPostId,
          body: payload,
        }).unwrap();
        feedback.showSuccess("Post updated", "Your changes were saved.");
        router.back();
        return;
      }
      const post = await createPost({
        societyId,
        channelId: communityChannel.id,
        body: payload,
      }).unwrap();
      const postId = post.post_id ?? post.id;
      feedback.showSuccess("Published", "Your post is live in Community.");
      router.replace(
        residentHubPostRoute(societyId, postId, communityChannel.id, "community"),
      );
    } catch (error) {
      feedback.showInfo(
        isEdit ? "Could not save" : "Could not publish",
        hubPublishErrorMessage(error),
      );
    }
  };

  const headerTitle = isEdit ? "Edit post" : "Create post";
  const loadingEdit = isEdit && postQuery.isLoading && !postQuery.data;

  const categorySection = useMemo(
    () =>
      categories.length ? (
        <View>
          <Text style={hubTheme.typography.fieldLabel}>Category</Text>
          <HubChipWrap>
            {categories.map((category) => (
              <HubChip
                key={category.id}
                label={category.name}
                selected={categoryId === category.id}
                onPress={() => setCategoryId(category.id)}
              />
            ))}
          </HubChipWrap>
        </View>
      ) : null,
    [categories, categoryId],
  );

  return (
    <ScrollView contentContainerStyle={styles.content} keyboardShouldPersistTaps="handled">
      <Stack gap="lg" style={styles.screen}>
        {showBackHeader ? (
          <ScreenBackHeader
            fallbackHomeRoute={residentHubCommunityRoute()}
            title={headerTitle}
          />
        ) : (
          <HubPageHeader
            subtitle={
              isEdit
                ? "Update your message or photos"
                : "Share with neighbours in your society"
            }
          />
        )}

        {loadingEdit ? <HubPostSkeleton /> : null}

        {!loadingEdit ? (
          <>
            {categorySection}
            <HubTextField
              label="Title (optional)"
              placeholder="Give your post a title"
              value={title}
              onChangeText={setTitle}
            />
            <HubTextField
              error={bodyError}
              label="Message"
              multiline
              placeholder="Share an update, question, or recommendation…"
              value={body}
              onChangeText={(text) => {
                setBody(text);
                if (bodyError) {
                  setBodyError(undefined);
                }
              }}
            />
            <View>
              <Text style={hubTheme.typography.fieldLabel}>Photos (optional)</Text>
              <View style={styles.attachGrid}>
                {attachments.map((slot) =>
                  slot.localUri ? (
                    <HubAttachmentTile
                      key={slot.key}
                      error={slot.error}
                      uri={slot.localUri}
                      uploading={slot.uploading}
                      onRemove={() => removeAttachment(slot)}
                    />
                  ) : slot.uploadId ? (
                    <ExistingAttachmentThumb
                      key={slot.key}
                      societyId={societyId ?? 0}
                      uploadId={slot.uploadId}
                      onRemove={() => removeAttachment(slot)}
                    />
                  ) : null,
                )}
                {canAddMore ? (
                  <HubAddAttachmentTile disabled={uploadsBusy} onPress={() => void addPhotos()} />
                ) : null}
              </View>
              <Text style={styles.attachHint}>Up to {HUB_MAX_ATTACHMENTS} photos (JPEG or PNG, 5 MB each)</Text>
            </View>
            <Button
              disabled={publishing || uploadsBusy}
              loading={publishing}
              title={isEdit ? "Save changes" : "Publish post"}
              onPress={() => void submit()}
            />
          </>
        ) : null}
      </Stack>
    </ScrollView>
  );
}

function ExistingAttachmentThumb({
  societyId,
  uploadId,
  onRemove,
}: {
  societyId: number;
  uploadId: number;
  onRemove: () => void;
}) {
  const query = useHubAttachmentQuery({ societyId, uploadId }, { skip: !societyId });
  return (
    <HubAttachmentTile
      error={query.isError ? "Unavailable" : undefined}
      uri={query.data?.url ?? "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"}
      uploading={query.isLoading}
      onRemove={onRemove}
    />
  );
}

const styles = StyleSheet.create({
  attachGrid: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.sm,
  },
  attachHint: {
    ...hubTheme.typography.meta,
    marginTop: spacing.xs,
  },
  content: {
    flexGrow: 1,
    paddingBottom: spacing["2xl"],
  },
  screen: {
    flex: 1,
  },
});
