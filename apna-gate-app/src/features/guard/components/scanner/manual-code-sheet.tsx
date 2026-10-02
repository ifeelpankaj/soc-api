import { useState } from "react";
import {
  KeyboardAvoidingView,
  Modal,
  Platform,
  Pressable,
  StyleSheet,
  Text,
  TextInput,
  View,
} from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";

import { Button } from "@/components/ui";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type ManualCodeSheetProps = {
  onClose: () => void;
  onSubmit: (code: string) => void;
  visible: boolean;
};

export function ManualCodeSheet({ onClose, onSubmit, visible }: ManualCodeSheetProps) {
  const insets = useSafeAreaInsets();
  const [code, setCode] = useState("");

  const handleClose = () => {
    setCode("");
    onClose();
  };

  const handleSubmit = () => {
    const trimmed = code.trim();
    if (!trimmed) {
      return;
    }
    onSubmit(trimmed);
    setCode("");
  };

  return (
    <Modal animationType="slide" transparent visible={visible} onRequestClose={handleClose}>
      <KeyboardAvoidingView
        behavior={Platform.OS === "ios" ? "padding" : undefined}
        style={styles.backdrop}
      >
        <Pressable accessibilityRole="button" style={styles.dismissArea} onPress={handleClose} />
        <View
          style={[
            styles.sheet,
            { paddingBottom: spacing["3xl"] + Math.max(insets.bottom, 0) },
          ]}
        >
          <View style={styles.handle} />
          <Text style={styles.title}>Enter Code Manually</Text>
          <Text style={styles.subtitle}>
            Paste the visitor QR token or the full check-in link below.
          </Text>
          <TextInput
            autoCapitalize="none"
            autoCorrect={false}
            multiline
            placeholder="Visitor token or QR link"
            placeholderTextColor={colors.text.placeholder}
            style={styles.input}
            value={code}
            onChangeText={setCode}
          />
          <View style={styles.actions}>
            <Pressable accessibilityRole="button" style={styles.cancelButton} onPress={handleClose}>
              <Text style={styles.cancelText}>Cancel</Text>
            </Pressable>
            <View style={styles.submitSlot}>
              <Button compact fullWidth title="Continue" disabled={!code.trim()} onPress={handleSubmit} />
            </View>
          </View>
        </View>
      </KeyboardAvoidingView>
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
    flex: 1,
    justifyContent: "flex-end",
  },
  cancelButton: {
    alignItems: "center",
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: 1,
    flex: 1,
    justifyContent: "center",
    minHeight: 44,
  },
  cancelText: {
    color: colors.text.secondary,
    fontSize: 15,
    fontWeight: "600",
  },
  dismissArea: {
    flex: 1,
  },
  handle: {
    alignSelf: "center",
    backgroundColor: colors.border.default,
    borderRadius: 999,
    height: 4,
    marginBottom: spacing.lg,
    width: 40,
  },
  input: {
    backgroundColor: colors.surface.input,
    borderColor: colors.border.input,
    borderRadius: radius.lg,
    borderWidth: 1,
    color: colors.brand.navy,
    fontSize: 15,
    marginTop: spacing.lg,
    minHeight: 96,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    textAlignVertical: "top",
  },
  sheet: {
    backgroundColor: colors.surface.card,
    borderTopLeftRadius: radius["2xl"],
    borderTopRightRadius: radius["2xl"],
    paddingHorizontal: spacing.xl,
    paddingTop: spacing.md,
  },
  submitSlot: {
    flex: 1.4,
  },
  subtitle: {
    color: colors.text.secondary,
    fontSize: 14,
    lineHeight: 20,
    marginTop: spacing.xs,
  },
  title: {
    color: colors.brand.navy,
    fontSize: 20,
    fontWeight: "700",
  },
});
