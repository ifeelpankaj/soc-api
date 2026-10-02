import { type ReactNode } from "react";
import { ScrollView, type ScrollViewProps } from "react-native";

import { layout } from "@/theme/layout";

type ScrollContainerProps = ScrollViewProps & {
  paddingBottom?: number;
  paddingHorizontal?: number;
  paddingTop?: number;
  children?: ReactNode;
};

export function ScrollContainer({
  paddingBottom = layout.screenPaddingBottom,
  paddingHorizontal = layout.screenPaddingHorizontal,
  paddingTop = layout.screenPaddingTop,
  contentContainerStyle,
  children,
  ...props
}: ScrollContainerProps) {
  return (
    <ScrollView
      contentContainerStyle={[
        {
          paddingBottom,
          paddingHorizontal,
          paddingTop,
        },
        contentContainerStyle,
      ]}
      nestedScrollEnabled
      {...props}
    >
      {children}
    </ScrollView>
  );
}
