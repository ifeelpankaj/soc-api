import { Stack } from "@/components/layout";
import { AppText } from "@/components/ui/app-text";

type TitleProps = {
  eyebrow?: string;
  title: string;
  subtitle?: string;
};

export function Title({ eyebrow, title, subtitle }: TitleProps) {
  return (
    <Stack gap="sm">
      {eyebrow ? (
        <AppText variant="eyebrow" color="brand">
          {eyebrow}
        </AppText>
      ) : null}
      <AppText variant="titleLarge" color="primary">
        {title}
      </AppText>
      {subtitle ? (
        <AppText variant="subtitle" color="secondary">
          {subtitle}
        </AppText>
      ) : null}
    </Stack>
  );
}

/** @deprecated Use Title instead */
export const ScreenHeader = Title;
