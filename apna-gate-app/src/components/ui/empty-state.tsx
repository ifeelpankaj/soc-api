import { StyleSheet, View } from "react-native";
import { SymbolView } from "expo-symbols";

import { Stack } from "@/components/layout/stack";
import { AppText } from "@/components/ui/app-text";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

import { Button } from "./button";

type EmptyStateProps = {
  title: string;
  message: string;
  actionLabel?: string;
  onAction?: () => void;
};

export function EmptyState({ title, message, actionLabel, onAction }: EmptyStateProps) {
  return (
    <View style={styles.container}>
      <View style={styles.icon}>
        <SymbolView name={{ ios: "tray", android: "inbox", web: "inbox" }} size={26} tintColor={colors.text.muted} />
      </View>
      <Stack gap="sm">
        <AppText variant="title" color="primary" style={styles.centered}>
          {title}
        </AppText>
        <AppText variant="body" color="secondary" style={styles.centered}>
          {message}
        </AppText>
      </Stack>
      {actionLabel && onAction ? <Button variant="secondary" title={actionLabel} onPress={onAction} /> : null}
    </View>
  );
}

const styles = StyleSheet.create({
  centered: {
    textAlign: "center",
  },
  container: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: 1,
    gap: spacing.lg,
    padding: spacing["2xl"],
  },
  icon: {
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.guard.sectionBg,
    borderRadius: radius["2xl"],
    height: 48,
    width: 48,
  },
});
