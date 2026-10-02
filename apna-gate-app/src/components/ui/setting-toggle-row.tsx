import { Pressable, StyleSheet, Switch } from "react-native";

import { Stack } from "@/components/layout";
import { AppText } from "@/components/ui/app-text";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type SettingToggleRowProps = {
  title: string;
  description?: string;
  value: boolean;
  disabled?: boolean;
  onValueChange: (value: boolean) => void;
};

export function SettingToggleRow({
  title,
  description,
  value,
  disabled,
  onValueChange,
}: SettingToggleRowProps) {
  return (
    <Pressable
      accessibilityRole="switch"
      accessibilityState={{ checked: value, disabled }}
      disabled={disabled}
      style={({ pressed }) => [
        styles.row,
        disabled && styles.disabled,
        pressed && !disabled && styles.pressed,
      ]}
      onPress={() => {
        if (!disabled) {
          onValueChange(!value);
        }
      }}
    >
      <Stack gap="xs" style={styles.copy}>
        <AppText variant="body" color="primary" style={styles.title}>
          {title}
        </AppText>
        {description ? (
          <AppText variant="bodySmall" color="secondary">
            {description}
          </AppText>
        ) : null}
      </Stack>
      <Switch
        disabled={disabled}
        value={value}
        onValueChange={onValueChange}
        trackColor={{ false: "#cbd5e1", true: colors.accent.selection }}
        thumbColor={value ? colors.guard.teal : colors.surface.screen}
      />
    </Pressable>
  );
}

const styles = StyleSheet.create({
  copy: {
    flex: 1,
  },
  disabled: {
    opacity: 0.55,
  },
  pressed: {
    opacity: 0.9,
  },
  row: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.lg,
    justifyContent: "space-between",
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.lg,
  },
  title: {
    fontWeight: "600",
  },
});
