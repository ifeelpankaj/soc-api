import { Pressable } from "react-native";

import { Row, Stack } from "@/components/layout";
import { AppText } from "@/components/ui/app-text";

type SectionTitleProps = {
  title: string;
  subtitle?: string;
  actionLabel?: string;
  onAction?: () => void;
  trailing?: string;
};

export function SectionTitle({
  title,
  subtitle,
  actionLabel,
  onAction,
  trailing,
}: SectionTitleProps) {
  return (
    <Row align="flex-start" gap="md">
      <Stack gap="xs" style={{ flex: 1 }}>
        <AppText variant="label" color="muted">
          {title}
        </AppText>
        {subtitle ? (
          <AppText variant="bodySmall" color="secondary">
            {subtitle}
          </AppText>
        ) : null}
      </Stack>
      {actionLabel && onAction ? (
        <Pressable accessibilityRole="button" onPress={onAction}>
          <AppText variant="bodySmall" color="operational" style={{ fontWeight: "700" }}>
            {actionLabel}
          </AppText>
        </Pressable>
      ) : null}
      {trailing ? (
        <AppText variant="bodySmall" color="operational" style={{ fontWeight: "700" }}>
          {trailing}
        </AppText>
      ) : null}
    </Row>
  );
}

/** @deprecated Use SectionTitle instead */
export const SectionHeader = SectionTitle;
