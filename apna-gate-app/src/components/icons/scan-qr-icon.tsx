import Svg, { Path, Rect } from "react-native-svg";

import type { SvgIconProps } from "@/components/icons/svgs/icon-props";

export function ScanQrIcon({
  size = 24,
  color = "#FFFFFF",
  strokeWidth = 2,
}: SvgIconProps) {
  return (
    <Svg fill="none" height={size} viewBox="0 0 24 24" width={size}>
      <Path
        d="M4 8V6C4 4.89543 4.89543 4 6 4H8M16 4H18C19.1046 4 20 4.89543 20 6V8M20 16V18C20 19.1046 19.1046 20 18 20H16M8 20H6C4.89543 20 4 19.1046 4 18V16"
        stroke={color}
        strokeLinecap="round"
        strokeWidth={strokeWidth}
      />
      <Rect height={6} stroke={color} strokeWidth={strokeWidth} width={6} x={9} y={9} />
    </Svg>
  );
}
