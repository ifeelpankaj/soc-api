import { StyleSheet, Text, TextInput, View, type TextInputProps } from "react-native";

import { hubTheme } from "@/features/hub/hub-theme";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type HubTextFieldProps = TextInputProps & {
  label?: string;
  error?: string;
  multiline?: boolean;
};

export function HubTextField({
  label,
  error,
  style,
  multiline,
  ...props
}: HubTextFieldProps) {
  return (
    <View>
      {label ? <Text style={hubTheme.typography.fieldLabel}>{label}</Text> : null}
      <View style={[styles.field, multiline && styles.multiline, error && styles.fieldError]}>
        <TextInput
          multiline={multiline}
          placeholderTextColor={colors.text.placeholder}
          style={[styles.input, multiline && styles.inputMultiline, style]}
          textAlignVertical={multiline ? "top" : "auto"}
          {...props}
        />
      </View>
      {error ? <Text style={styles.error}>{error}</Text> : null}
    </View>
  );
}

const styles = StyleSheet.create({
  error: {
    color: colors.status.error,
    fontSize: 12,
    marginTop: spacing.xs,
  },
  field: {
    backgroundColor: colors.surface.card,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.lg,
    borderWidth: StyleSheet.hairlineWidth,
    padding: spacing.md,
  },
  fieldError: {
    borderColor: colors.status.error,
  },
  input: {
    color: colors.brand.navy,
    fontSize: 15,
    lineHeight: 22,
    minHeight: 24,
  },
  inputMultiline: {
    minHeight: 120,
  },
  multiline: {
    minHeight: 140,
  },
});
