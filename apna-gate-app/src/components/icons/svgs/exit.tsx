import Svg, { Path } from "react-native-svg";

import type { SvgIconProps } from "@/components/icons/svgs/icon-props";

export function ExitSvg({
  size = 24,
  color = "#211714",
  strokeWidth = 2,
}: SvgIconProps) {
  return (
    <Svg fill="none" height={size} viewBox="0 0 24 24" width={size}>
      <Path
        d="M10 7V5C10 4.44772 10.4477 4 11 4H19C19.5523 4 20 4.44772 20 5V19C20 19.5523 19.5523 20 19 20H11C10.4477 20 10 19.5523 10 19V17"
        stroke={color}
        strokeLinecap="round"
        strokeWidth={strokeWidth}
      />
      <Path
        d="M14 12H4M7 9L4 12L7 15"
        stroke={color}
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={strokeWidth}
      />
    </Svg>
  );
}
