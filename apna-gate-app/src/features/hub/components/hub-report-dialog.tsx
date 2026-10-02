import { useState } from "react";
import { Modal, StyleSheet, Text, View } from "react-native";

import { Button } from "@/components/ui";
import { HubTextField } from "@/features/hub/components/ui";
import { hubTheme } from "@/features/hub/hub-theme";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type HubReportDialogProps = {
  visible: boolean;
  loading?: boolean;
  onClose: () => void;
  onSubmit: (reason: string) => void;
};

export function HubReportDialog({
  visible,
  loading,
  onClose,
  onSubmit,
}: HubReportDialogProps) {
  const [reason, setReason] = useState("");

  const submit = () => {
    const trimmed = reason.trim();
    if (trimmed.length < 3) {
      return;
    }
    onSubmit(trimmed);
    setReason("");
  };

  return (
    <Modal animationType="fade" transparent visible={visible} onRequestClose={onClose}>
      <View style={styles.backdrop}>
        <View style={styles.sheet}>
          <Text style={styles.title}>Report post</Text>
          <Text style={styles.message}>
            Tell us briefly what is wrong. Society admins may review this report.
          </Text>
          <HubTextField
            multiline
            placeholder="Reason for report"
            value={reason}
            onChangeText={setReason}
          />
          <View style={styles.actions}>
            <Button fullWidth title="Cancel" variant="ghost" onPress={onClose} />
            <Button
              disabled={reason.trim().length < 3}
              fullWidth
              loading={loading}
              title="Submit report"
              variant="danger"
              onPress={submit}
            />
          </View>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  actions: {
    gap: spacing.sm,
    marginTop: spacing.md,
  },
  backdrop: {
    backgroundColor: "rgba(16, 29, 54, 0.45)",
    flex: 1,
    justifyContent: "flex-end",
    padding: spacing.lg,
  },
  message: {
    ...hubTheme.typography.pageSubtitle,
    marginBottom: spacing.md,
  },
  sheet: {
    backgroundColor: colors.surface.card,
    borderRadius: radius.xl,
    padding: spacing.lg,
  },
  title: {
    ...hubTheme.typography.postTitle,
    fontSize: 18,
    marginBottom: spacing.xs,
  },
});
