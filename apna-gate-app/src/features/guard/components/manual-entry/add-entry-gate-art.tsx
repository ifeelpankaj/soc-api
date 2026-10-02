import { StyleSheet, View, type ViewStyle } from "react-native";
import Svg, { Path, Rect } from "react-native-svg";

import { colors } from "@/theme/colors";

type AddEntryGateArtProps = {
  style?: ViewStyle;
};

export function AddEntryGateArt({ style }: AddEntryGateArtProps) {
  return (
    <View pointerEvents="none" style={[styles.wrap, style]}>
      <Svg height={72} viewBox="0 0 96 72" width={96}>
        <Rect fill={colors.brand.orangeSoft} height={18} rx={3} width={12} x={8} y="18" />
        <Rect fill={colors.brand.orangeSoft} height={18} rx={3} width={12} x="76" y="18" />
        <Path
          d="M20 36 H76 V54 C76 58 72 62 68 62 H28 C24 62 20 58 20 54 Z"
          fill="rgba(255, 106, 26, 0.14)"
        />
        <Path
          d="M24 36 H72 V52 C72 55 69 58 66 58 H30 C27 58 24 55 24 52 Z"
          fill="rgba(255, 106, 26, 0.22)"
        />
        <Path
          d="M30 36 H66 V48 C66 50 64 52 62 52 H34 C32 52 30 50 30 48 Z"
          fill="rgba(255, 106, 26, 0.3)"
        />
        <Rect fill="rgba(255, 106, 26, 0.18)" height={8} rx={2} width={8} x="44" y="10" />
      </Svg>
    </View>
  );
}

const styles = StyleSheet.create({
  wrap: {
    opacity: 0.85,
  },
});
