import Svg, { Circle, Path } from "react-native-svg";

import type { SvgIconProps } from "@/components/icons/svgs/icon-props";

export function SearchSvg({
  size = 24,
  color = "#211714",
  strokeWidth = 2,
}: SvgIconProps) {
  return (
    <Svg fill="none" height={size} viewBox="0 0 24 24" width={size}>
      <Circle cx={11} cy={11} r={7} stroke={color} strokeWidth={strokeWidth} />
      <Path d="M20 20L16.5 16.5" stroke={color} strokeLinecap="round" strokeWidth={strokeWidth} />
    </Svg>
  );
}
