import { StyleSheet, Text, View } from "react-native";

import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type StatusPillProps = {
  status?: string | null;
};

function getStatusStyle(status?: string | null) {
  switch (status) {
    case "active":
    case "approved":
    case "verified":
      return { backgroundColor: "#ecfdf5", color: "#047857" };
    case "pending":
    case "trial":
      return { backgroundColor: "#fffbeb", color: "#b45309" };
    case "suspended":
    case "blocked":
    case "rejected":
      return { backgroundColor: "#fff1f2", color: "#be123c" };
    default:
      return { backgroundColor: "#f1f5f9", color: "#334155" };
  }
}

function formatStatus(status?: string | null) {
  if (!status) {
    return "Unknown";
  }

  return status
    .split("_")
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

export function StatusPill({ status }: StatusPillProps) {
  const statusStyle = getStatusStyle(status);

  return (
    <View style={styles.wrapper}>
      <Text style={[styles.pill, { backgroundColor: statusStyle.backgroundColor, color: statusStyle.color }]}>
        {formatStatus(status)}
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  pill: {
    ...typography.caption,
    borderRadius: radius["2xl"],
    fontWeight: "600",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.xs,
  },
  wrapper: {
    alignSelf: "flex-start",
    borderRadius: radius["2xl"],
  },
});
