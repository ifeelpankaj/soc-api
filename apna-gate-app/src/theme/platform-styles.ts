import { Platform, type TextStyle } from "react-native";

/** Removes extra top/bottom font padding on Android so text inputs align correctly. */
export const androidCompactText: TextStyle = Platform.select({
  android: { includeFontPadding: false },
  default: {},
})!;
