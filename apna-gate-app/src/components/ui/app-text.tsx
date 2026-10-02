import { Text, type TextProps, type TextStyle } from "react-native";

import { colors } from "@/theme/colors";
import { typography } from "@/theme/typography";

type TypographyVariant = keyof typeof typography;

type ColorToken =
  | "primary"
  | "secondary"
  | "muted"
  | "placeholder"
  | "inverse"
  | "brand"
  | "operational"
  | "success"
  | "warning"
  | "error"
  | "announcement";

const colorMap: Record<ColorToken, string> = {
  primary: colors.text.primary,
  secondary: colors.text.secondary,
  muted: colors.text.muted,
  placeholder: colors.text.placeholder,
  inverse: colors.text.inverse,
  brand: colors.brand.orange,
  operational: colors.operational.primary,
  success: colors.status.success,
  warning: colors.status.warning,
  error: colors.status.error,
  announcement: colors.announcement.text,
};

type AppTextProps = TextProps & {
  variant?: TypographyVariant;
  color?: ColorToken | string;
};

export function AppText({ variant = "body", color = "primary", style, ...props }: AppTextProps) {
  const colorValue = color in colorMap ? colorMap[color as ColorToken] : color;

  return (
    <Text style={[typography[variant], { color: colorValue }, style]} {...props} />
  );
}

export type { TypographyVariant, ColorToken, TextStyle };
