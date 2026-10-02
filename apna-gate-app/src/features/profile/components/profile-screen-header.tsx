import type { Href } from "expo-router";
import { Pressable, StyleSheet, Text } from "react-native";

import { AppIcon } from "@/components/icons";
import { Row } from "@/components/layout";
import { useBackAction } from "@/lib/navigation/use-back-action";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type ProfileScreenHeaderProps = {
  fallbackHomeRoute: Href;
  title: string;
  onBack?: () => void;
  backDisabled?: boolean;
};

export function ProfileScreenHeader({
  fallbackHomeRoute,
  title,
  onBack,
  backDisabled = false,
}: ProfileScreenHeaderProps) {
  const handleBack = useBackAction(fallbackHomeRoute);

  return (
    <Row align="center" gap="md" style={styles.header}>
      <Pressable
        accessibilityLabel="Go back"
        accessibilityRole="button"
        accessibilityState={{ disabled: backDisabled }}
        disabled={backDisabled}
        hitSlop={8}
        style={({ pressed }) => [
          styles.backButton,
          pressed && styles.backButtonPressed,
        ]}
        onPress={onBack ?? handleBack}
      >
        <AppIcon color={colors.guard.text} name="back" size={18} />
      </Pressable>
      <Text style={styles.title}>{title}</Text>
    </Row>
  );
}

const styles = StyleSheet.create({
  backButton: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius["2xl"],
    borderWidth: 1,
    height: 44,
    justifyContent: "center",
    width: 44,
  },
  backButtonPressed: {
    backgroundColor: colors.brand.orangeSoft,
  },
  header: {
    justifyContent: "flex-start",
    paddingBottom: spacing.sm,
    paddingTop: spacing.xs,
  },
  title: {
    color: colors.guard.text,
    flex: 1,
    fontSize: 18,
    fontWeight: "700",
  },
});
