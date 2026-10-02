import { Pressable, StyleSheet, Text } from "react-native";

import { Row, Stack } from "@/components/layout";
import { SectionHeader } from "@/components/ui";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type ProfileQuickActionsProps = {
  onChangePhoto: () => void;
  onEditProfile: () => void;
};

export function ProfileQuickActions({ onChangePhoto, onEditProfile }: ProfileQuickActionsProps) {
  return (
    <Stack gap={10}>
      <SectionHeader title="Quick Actions" />
      <Row align="center" gap="md" justify="flex-start">
        <Pressable
          accessibilityRole="button"
          style={({ pressed }) => [styles.button, pressed && styles.buttonPressed]}
          onPress={onEditProfile}
        >
          <Text style={styles.buttonText}>Edit Profile</Text>
        </Pressable>
        <Pressable
          accessibilityRole="button"
          style={({ pressed }) => [styles.button, pressed && styles.buttonPressed]}
          onPress={onChangePhoto}
        >
          <Text style={styles.buttonText}>Change Photo</Text>
        </Pressable>
      </Row>
    </Stack>
  );
}

const styles = StyleSheet.create({
  button: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius["2xl"],
    borderWidth: 1,
    flex: 1,
    justifyContent: "center",
    minHeight: 48,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.md,
  },
  buttonPressed: {
    opacity: 0.85,
  },
  buttonText: {
    color: colors.text.primary,
    textAlign: "center",
    fontSize: 14,
    fontWeight: "700",
  },
});
