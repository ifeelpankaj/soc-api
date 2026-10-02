import Svg, { Path } from "react-native-svg";

import type { SvgIconProps } from "@/components/icons/svgs/icon-props";

export function FilterSvg({
  size = 24,
  color = "#211714",
  strokeWidth = 2,
}: SvgIconProps) {
  return (
    <Svg fill="none" height={size} viewBox="0 0 24 24" width={size}>
      <Path
        d="M4 6H20M7 12H17M10 18H14"
        stroke={color}
        strokeLinecap="round"
        strokeWidth={strokeWidth}
      />
    </Svg>
  );
}
