import { useState } from "react";
import { StyleSheet, Text, TextInput, type TextInputProps, View } from "react-native";

import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type InputProps = TextInputProps & {
  label: string;
  error?: string;
};

type InputFocusEvent = Parameters<NonNullable<TextInputProps["onFocus"]>>[0];
type InputBlurEvent = Parameters<NonNullable<TextInputProps["onBlur"]>>[0];

export function Input({ label, error, editable = true, style, onBlur, onFocus, ...props }: InputProps) {
  const [isFocused, setIsFocused] = useState(false);
  const isActive = isFocused && editable;

  const handleFocus = (event: InputFocusEvent) => {
    setIsFocused(true);
    onFocus?.(event);
  };

  const handleBlur = (event: InputBlurEvent) => {
    setIsFocused(false);
    onBlur?.(event);
  };

  return (
    <View style={styles.wrapper}>
      <Text
        style={[
          styles.label,
          error ? styles.labelError : isActive ? styles.labelActive : styles.labelDefault,
        ]}
      >
        {label}
      </Text>
      <TextInput
        cursorColor={colors.guard.teal}
        editable={editable}
        onBlur={handleBlur}
        onFocus={handleFocus}
        placeholderTextColor={colors.text.placeholder}
        selectionColor={colors.guard.tealSoft}
        style={[
          styles.input,
          error ? styles.inputError : isActive ? styles.inputActive : styles.inputDefault,
          !editable && styles.inputDisabled,
          style,
        ]}
        {...props}
      />
      {error ? <Text style={styles.error}>{error}</Text> : null}
    </View>
  );
}

const styles = StyleSheet.create({
  error: {
    ...typography.bodySmall,
    color: colors.status.error,
    fontWeight: "500",
  },
  input: {
    ...typography.body,
    backgroundColor: colors.surface.input,
    borderRadius: radius.md,
    borderWidth: 1,
    color: colors.text.primary,
    minHeight: layout.inputHeight,
    paddingBottom: spacing.sm,
    paddingHorizontal: spacing.md,
    paddingTop: spacing.lg,
  },
  inputActive: {
    borderColor: colors.guard.teal,
  },
  inputDefault: {
    borderColor: colors.border.input,
  },
  inputDisabled: {
    backgroundColor: colors.guard.sectionBg,
    color: colors.text.muted,
  },
  inputError: {
    borderColor: colors.status.error,
  },
  label: {
    ...typography.bodySmall,
    backgroundColor: colors.surface.input,
    fontWeight: "500",
    left: spacing.md,
    paddingHorizontal: spacing.xs,
    position: "absolute",
    top: 0,
    zIndex: 10,
  },
  labelActive: {
    color: colors.guard.teal,
  },
  labelDefault: {
    color: colors.text.primary,
  },
  labelError: {
    color: colors.status.error,
  },
  wrapper: {
    gap: spacing.sm,
    paddingTop: spacing.sm,
    width: "100%",
  },
});
