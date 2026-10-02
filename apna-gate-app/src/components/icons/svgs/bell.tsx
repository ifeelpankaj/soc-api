import Svg, { Path } from "react-native-svg";

import type { SvgIconProps } from "@/components/icons/svgs/icon-props";

export function BellSvg({
  size = 24,
  color = "#211714",
  strokeWidth = 2,
}: SvgIconProps) {
  return (
    <Svg fill="none" height={size} viewBox="0 0 24 24" width={size}>
      <Path
        d="M18 16V11C18 7.68629 15.3137 5 12 5C8.68629 5 6 7.68629 6 11V16L5 18H19L18 16Z"
        stroke={color}
        strokeLinejoin="round"
        strokeWidth={strokeWidth}
      />
      <Path
        d="M10 20C10.5523 21.1046 11.8954 22 13.5 22C15.1046 22 16.4477 21.1046 17 20"
        stroke={color}
        strokeLinecap="round"
        strokeWidth={strokeWidth}
      />
    </Svg>
  );
}
