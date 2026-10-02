import { type ReactNode } from "react";
import { StyleSheet, View, type ViewProps, type ViewStyle } from "react-native";

import { spacing, type spacing as SpacingType } from "@/theme/spacing";

type SpacingKey = keyof typeof SpacingType;

type RowProps = ViewProps & {
  gap?: SpacingKey | number;
  align?: ViewStyle["alignItems"];
  justify?: ViewStyle["justifyContent"];
  wrap?: boolean;
  children?: ReactNode;
};

export function Row({
  gap = "sm",
  align = "center",
  justify = "space-between",
  wrap = false,
  style,
  children,
  ...props
}: RowProps) {
  const gapValue = typeof gap === "number" ? gap : spacing[gap];

  return (
    <View
      style={[
        styles.row,
        {
          gap: gapValue,
          alignItems: align,
          justifyContent: justify,
          flexWrap: wrap ? "wrap" : "nowrap",
        },
        style,
      ]}
      {...props}
    >
      {children}
    </View>
  );
}

const styles = StyleSheet.create({
  row: {
    flexDirection: "row",
  },
});
