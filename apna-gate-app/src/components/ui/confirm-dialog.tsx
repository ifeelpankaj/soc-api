import { Modal, Pressable, StyleSheet, Text, View } from "react-native";

import { Button } from "@/components/ui/button";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type ConfirmDialogProps = {
  cancelLabel?: string;
  confirmLabel?: string;
  loading?: boolean;
  message: string;
  onCancel: () => void;
  onConfirm: () => void;
  title: string;
  visible: boolean;
};

export function ConfirmDialog({
  cancelLabel = "Cancel",
  confirmLabel = "Confirm",
  loading,
  message,
  onCancel,
  onConfirm,
  title,
  visible,
}: ConfirmDialogProps) {
  return (
    <Modal animationType="fade" transparent visible={visible} onRequestClose={onCancel}>
      <View style={styles.backdrop}>
        <View style={styles.card}>
          <Text style={styles.title}>{title}</Text>
          <Text style={styles.message}>{message}</Text>
          <View style={styles.actions}>
            <Pressable style={styles.cancelButton} onPress={onCancel}>
              <Text style={styles.cancelText}>{cancelLabel}</Text>
            </Pressable>
            <View style={styles.confirmSlot}>
              <Button compact fullWidth loading={loading} title={confirmLabel} onPress={onConfirm} />
            </View>
          </View>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  actions: {
    flexDirection: "row",
    gap: spacing.md,
    marginTop: spacing.lg,
  },
  backdrop: {
    alignItems: "center",
    backgroundColor: "rgba(15, 23, 42, 0.45)",
    flex: 1,
    justifyContent: "center",
    padding: spacing.xl,
  },
  cancelButton: {
    alignItems: "center",
    borderColor: colors.guard.border,
    borderRadius: radius.lg,
    borderWidth: 1,
    flex: 1,
    justifyContent: "center",
    minHeight: 44,
    paddingHorizontal: spacing.md,
  },
  cancelText: {
    color: colors.text.secondary,
    fontSize: 15,
    fontWeight: "600",
  },
  card: {
    backgroundColor: colors.surface.card,
    borderRadius: radius.xl,
    maxWidth: 420,
    padding: spacing.xl,
    width: "100%",
  },
  confirmSlot: {
    flex: 1.4,
  },
  message: {
    color: colors.text.secondary,
    fontSize: 15,
    lineHeight: 22,
    marginTop: spacing.sm,
  },
  title: {
    color: colors.text.primary,
    fontSize: 20,
    fontWeight: "700",
  },
});
