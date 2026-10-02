import Svg, { Path, Rect } from "react-native-svg";

import type { SvgIconProps } from "@/components/icons/svgs/icon-props";

export function CalendarSvg({
  size = 24,
  color = "#211714",
  strokeWidth = 2,
}: SvgIconProps) {
  return (
    <Svg fill="none" height={size} viewBox="0 0 24 24" width={size}>
      <Rect height={16} rx={2} stroke={color} strokeWidth={strokeWidth} width={18} x={3} y={5} />
      <Path d="M3 9H21M8 3V7M16 3V7" stroke={color} strokeLinecap="round" strokeWidth={strokeWidth} />
    </Svg>
  );
}
