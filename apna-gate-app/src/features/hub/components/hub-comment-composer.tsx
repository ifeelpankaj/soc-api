import { useState } from "react";
import { StyleSheet, View } from "react-native";

import { Button } from "@/components/ui";
import { hubUploadErrorMessage } from "@/features/hub/hub-api-errors";
import { useAppFeedback } from "@/features/shared/use-app-feedback";
import { HubTextField } from "@/features/hub/components/ui";
import {
  hubUploadFormData,
  pickHubPhotos,
  type HubPendingAttachment,
} from "@/features/hub/hub-upload";
import {
  useHubCreateCommentMutation,
  useHubUploadAttachmentMutation,
} from "@/lib/api/hub-api";
import { spacing } from "@/theme/spacing";
import { HubAddAttachmentTile, HubAttachmentTile } from "@/features/hub/components/ui/hub-attachment-tile";

type HubCommentComposerProps = {
  societyId: number;
  postId: number;
  onPosted?: () => void;
};

export function HubCommentComposer({
  societyId,
  postId,
  onPosted,
}: HubCommentComposerProps) {
  const [body, setBody] = useState("");
  const [attachment, setAttachment] = useState<HubPendingAttachment | null>(null);
  const feedback = useAppFeedback();
  const [createComment, createState] = useHubCreateCommentMutation();
  const [uploadAttachment] = useHubUploadAttachmentMutation();

  const addPhoto = async () => {
    if (attachment) {
      return;
    }
    let localUri = "";
    try {
      const picked = await pickHubPhotos(1);
      const photo = picked[0];
      if (!photo) {
        return;
      }
      localUri = photo.uri;
      const pending: HubPendingAttachment = {
        localUri: photo.uri,
        uploading: true,
      };
      setAttachment(pending);
      const form = await hubUploadFormData(photo);
      const uploaded = await uploadAttachment({ societyId, body: form }).unwrap();
      setAttachment({ localUri: photo.uri, uploadId: uploaded.id });
    } catch (error) {
      setAttachment({
        localUri,
        error: hubUploadErrorMessage(error),
      });
    }
  };

  const submit = async () => {
    const trimmed = body.trim();
    if (trimmed.length < 1) {
      return;
    }
    if (attachment?.uploading) {
      return;
    }
    try {
      await createComment({
        societyId,
        postId,
        body: {
          body: trimmed,
          attachment_ids: attachment?.uploadId ? [attachment.uploadId] : undefined,
        },
      }).unwrap();
      setBody("");
      setAttachment(null);
      onPosted?.();
    } catch (error) {
      feedback.showError("Could not post comment", error, "Try again.");
    }
  };

  const busy = createState.isLoading || attachment?.uploading;

  return (
    <View style={styles.wrap}>
      <HubTextField
        multiline
        placeholder="Write a comment…"
        value={body}
        onChangeText={setBody}
      />
      <View style={styles.attachRow}>
        {!attachment ? (
          <HubAddAttachmentTile disabled={busy} onPress={() => void addPhoto()} />
        ) : (
          <HubAttachmentTile
            error={attachment.error}
            uri={attachment.localUri}
            uploading={attachment.uploading}
            onRemove={() => setAttachment(null)}
          />
        )}
      </View>
      <Button
        disabled={busy || body.trim().length < 1}
        loading={createState.isLoading}
        title="Post comment"
        variant="secondary"
        onPress={() => void submit()}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  attachRow: {
    flexDirection: "row",
    marginVertical: spacing.sm,
  },
  wrap: {
    gap: spacing.sm,
    paddingTop: spacing.md,
  },
});
