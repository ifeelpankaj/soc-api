import { type ReactNode } from "react";
import { StyleSheet, View, type ViewProps, type ViewStyle } from "react-native";

import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { spacing, type spacing as SpacingType } from "@/theme/spacing";

type SpacingKey = keyof typeof SpacingType;

type ContainerProps = ViewProps & {
  flex?: boolean;
  padding?: boolean | SpacingKey;
  backgroundColor?: string;
  children?: ReactNode;
};

export function Container({
  flex = true,
  padding = false,
  backgroundColor = colors.surface.screen,
  style,
  children,
  ...props
}: ContainerProps) {
  const paddingStyle: ViewStyle | undefined =
    padding === true
      ? { paddingHorizontal: layout.screenPaddingHorizontal }
      : padding
        ? { padding: spacing[padding] }
        : undefined;

  return (
    <View
      style={[flex && styles.flex, { backgroundColor }, paddingStyle, style]}
      {...props}
    >
      {children}
    </View>
  );
}

const styles = StyleSheet.create({
  flex: {
    flex: 1,
  },
});
