import { StyleSheet } from "react-native";

import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

export const hubTheme = {
  screenPaddingBottom: spacing.xl,
  card: {
    backgroundColor: colors.surface.card,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.lg,
    borderWidth: StyleSheet.hairlineWidth,
    padding: spacing.lg,
    ...shadows.sm,
  },
  tones: {
    announcement: {
      accent: colors.dashboard.actionPurple,
      soft: colors.dashboard.actionPurpleSoft,
    },
    community: {
      accent: colors.dashboard.actionBlue,
      soft: colors.dashboard.actionBlueSoft,
    },
  },
  typography: {
    pageTitle: {
      color: colors.brand.navy,
      fontSize: 28,
      fontWeight: "700" as const,
      letterSpacing: -0.4,
    },
    pageSubtitle: {
      color: colors.text.muted,
      fontSize: 14,
      lineHeight: 20,
    },
    sectionEyebrow: {
      color: colors.text.secondary,
      fontSize: 13,
      fontWeight: "700" as const,
      letterSpacing: 0.6,
      textTransform: "uppercase" as const,
    },
    postTitle: {
      color: colors.brand.navy,
      fontSize: 20,
      fontWeight: "700" as const,
      letterSpacing: -0.2,
    },
    postBody: {
      color: colors.brand.navy,
      fontSize: 15,
      lineHeight: 24,
    },
    meta: {
      color: colors.text.muted,
      fontSize: 12,
    },
    fieldLabel: {
      color: colors.text.secondary,
      fontSize: 13,
      fontWeight: "600" as const,
      marginBottom: spacing.xs,
    },
  },
  pressOpacity: 0.92,
  minTouchTarget: 44,
};
