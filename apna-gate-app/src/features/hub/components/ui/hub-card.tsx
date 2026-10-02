import { StyleSheet, View, type ViewProps } from "react-native";

import { hubTheme } from "@/features/hub/hub-theme";

type HubCardProps = ViewProps & {
  accent?: boolean;
};

export function HubCard({ style, accent, children, ...props }: HubCardProps) {
  return (
    <View
      style={[styles.card, accent && styles.accent, style]}
      {...props}
    >
      {children}
    </View>
  );
}

const styles = StyleSheet.create({
  accent: {
    borderLeftColor: hubTheme.tones.community.accent,
    borderLeftWidth: 3,
  },
  card: {
    ...hubTheme.card,
  },
});
