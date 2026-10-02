import { Platform, StatusBar as RNStatusBar } from "react-native";
import { StatusBar } from "expo-status-bar";

import { colors } from "@/theme/colors";

type AppStatusBarProps = {
  style?: "auto" | "dark" | "light";
};

export function AppStatusBar({ style = "dark" }: AppStatusBarProps) {
  return (
    <>
      <StatusBar style={style} />
      {Platform.OS === "android" ? (
        <RNStatusBar
          backgroundColor={colors.surface.card}
          barStyle={style === "light" ? "light-content" : "dark-content"}
          translucent={false}
        />
      ) : null}
    </>
  );
}
