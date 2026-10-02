import { colors } from "./colors";
import { layout } from "./layout";
import { radius } from "./radius";
import { shadows } from "./shadows";
import { spacing } from "./spacing";
import { typography } from "./typography";

export { colors, layout, radius, shadows, spacing, typography };

/** @deprecated Use individual token imports from `@/theme` instead. */
export const theme = {
  brand: colors.brand,
  text: colors.text,
  border: colors.border,
  surface: colors.surface,
  operational: colors.operational,
  status: colors.status,
  announcement: colors.announcement,
  guard: {
    ...colors.guard,
    cardShadow: "0 2px 12px rgba(34, 40, 54, 0.06)",
    heroShadow: "0 8px 28px rgba(13, 148, 136, 0.14)",
    ctaShadow: "0 10px 28px rgba(13, 148, 136, 0.32)",
    inputHeight: layout.inputHeight,
    buttonHeight: layout.buttonHeight,
    purposeCardWidth: layout.purposeCardWidth,
    purposeCardHeight: layout.purposeCardHeight,
  },
  selection: {
    ...colors.selection,
    cardRadius: radius.xl,
  },
} as const;
