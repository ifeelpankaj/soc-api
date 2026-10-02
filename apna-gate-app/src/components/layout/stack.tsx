import { type ReactNode } from "react";
import { StyleSheet, View, type ViewProps, type ViewStyle } from "react-native";

import { spacing, type spacing as SpacingType } from "@/theme/spacing";

type SpacingKey = keyof typeof SpacingType;

type StackProps = ViewProps & {
  gap?: SpacingKey | number;
  align?: ViewStyle["alignItems"];
  justify?: ViewStyle["justifyContent"];
  children?: ReactNode;
};

export function Stack({
  gap = "sm",
  align,
  justify,
  style,
  children,
  ...props
}: StackProps) {
  const gapValue = typeof gap === "number" ? gap : spacing[gap];

  return (
    <View style={[styles.stack, { gap: gapValue, alignItems: align, justifyContent: justify }, style]} {...props}>
      {children}
    </View>
  );
}

const styles = StyleSheet.create({
  stack: {
    flexDirection: "column",
  },
});
