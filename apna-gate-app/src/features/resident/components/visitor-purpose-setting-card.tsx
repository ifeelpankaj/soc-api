import { Pressable, StyleSheet, Text, View } from "react-native";

import { Card, PurposeBadge } from "@/components/ui";
import { titleize } from "@/features/guard/guard-utils";
import type { ModelsFlatVisitorSettingsResponse } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

export type VisitorAccessMode = "blocked" | "approve" | "auto";

const ACCESS_OPTIONS: {
  value: VisitorAccessMode;
  label: string;
  hint: string;
}[] = [
  {
    value: "blocked",
    label: "Not allowed",
    hint: "This visitor type cannot come to your flat.",
  },
  {
    value: "approve",
    label: "Approve first",
    hint: "You get notified and must approve before they enter.",
  },
  {
    value: "auto",
    label: "Auto entry",
    hint: "They can check in without waiting for you.",
  },
];

export function settingToAccessMode(setting: ModelsFlatVisitorSettingsResponse): VisitorAccessMode {
  if (setting.is_enabled === false) {
    return "blocked";
  }

  if (setting.approval_required === true) {
    return "approve";
  }

  return "auto";
}

export function accessModeToPatch(mode: VisitorAccessMode): {
  approval_required: boolean;
  is_enabled: boolean;
} {
  switch (mode) {
    case "blocked":
      return { is_enabled: false, approval_required: true };
    case "approve":
      return { is_enabled: true, approval_required: true };
    case "auto":
      return { is_enabled: true, approval_required: false };
  }
}

function accessModeSummary(purpose: string | undefined | null, mode: VisitorAccessMode) {
  const label = purpose ? titleize(purpose) : "Visitor";

  switch (mode) {
    case "blocked":
      return `${label} visits are turned off for your flat.`;
    case "approve":
      return `${label} visitors need your approval before they can enter.`;
    case "auto":
      return `${label} visitors can enter without your approval.`;
  }
}

type VisitorPurposeSettingCardProps = {
  editable: boolean;
  setting: ModelsFlatVisitorSettingsResponse;
  onChange: (mode: VisitorAccessMode) => void;
};

export function VisitorPurposeSettingCard({
  editable,
  onChange,
  setting,
}: VisitorPurposeSettingCardProps) {
  const mode = settingToAccessMode(setting);
  const summary = accessModeSummary(setting.purpose, mode);

  return (
    <Card style={styles.card}>
      <View style={styles.header}>
        <PurposeBadge purpose={setting.purpose} />
      </View>

      <Text style={styles.summary}>{summary}</Text>

      {editable ? (
        <View style={styles.options}>
          {ACCESS_OPTIONS.map((option) => {
            const selected = option.value === mode;

            return (
              <Pressable
                key={option.value}
                accessibilityRole="radio"
                accessibilityState={{ checked: selected }}
                style={[styles.option, selected && styles.optionSelected]}
                onPress={() => onChange(option.value)}
              >
                <View style={[styles.radio, selected && styles.radioSelected]}>
                  {selected ? <View style={styles.radioDot} /> : null}
                </View>
                <View style={styles.optionCopy}>
                  <Text style={[styles.optionLabel, selected && styles.optionLabelSelected]}>
                    {option.label}
                  </Text>
                  <Text style={styles.optionHint}>{option.hint}</Text>
                </View>
              </Pressable>
            );
          })}
        </View>
      ) : (
        <View style={styles.readOnlyBadge}>
          <Text style={styles.readOnlyBadgeText}>
            {ACCESS_OPTIONS.find((option) => option.value === mode)?.label ?? "—"}
          </Text>
        </View>
      )}
    </Card>
  );
}

const styles = StyleSheet.create({
  card: {
    gap: spacing.md,
  },
  header: {
    alignItems: "flex-start",
  },
  option: {
    alignItems: "flex-start",
    backgroundColor: colors.surface.screen,
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  optionCopy: {
    flex: 1,
    gap: 2,
  },
  optionHint: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  optionLabel: {
    ...typography.body,
    color: colors.text.primary,
    fontWeight: "600",
  },
  optionLabelSelected: {
    color: colors.guard.teal,
  },
  optionSelected: {
    backgroundColor: colors.guard.tealSoft,
    borderColor: colors.guard.teal,
  },
  options: {
    gap: spacing.sm,
  },
  radio: {
    alignItems: "center",
    borderColor: colors.border.default,
    borderRadius: 999,
    borderWidth: 2,
    height: 20,
    justifyContent: "center",
    marginTop: 2,
    width: 20,
  },
  radioDot: {
    backgroundColor: colors.guard.teal,
    borderRadius: 999,
    height: 10,
    width: 10,
  },
  radioSelected: {
    borderColor: colors.guard.teal,
  },
  readOnlyBadge: {
    alignSelf: "flex-start",
    backgroundColor: colors.surface.screen,
    borderRadius: radius["2xl"],
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.xs,
  },
  readOnlyBadgeText: {
    ...typography.bodySmall,
    color: colors.text.secondary,
    fontWeight: "600",
  },
  summary: {
    ...typography.body,
    color: colors.text.primary,
    fontWeight: "500",
  },
});
