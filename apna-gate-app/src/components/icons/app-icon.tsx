import { View } from "react-native";

import { BackSvg } from "@/components/icons/svgs/back";
import { BellSvg } from "@/components/icons/svgs/bell";
import { CalendarSvg } from "@/components/icons/svgs/calendar";
import { ExitSvg } from "@/components/icons/svgs/exit";
import { FilterSvg } from "@/components/icons/svgs/filter";
import { HourglassSvg } from "@/components/icons/svgs/hourglass";
import { SearchSvg } from "@/components/icons/svgs/search";

export type AppIconName =
  | "search"
  | "filter"
  | "calendar"
  | "hourglass"
  | "exit"
  | "bell"
  | "back";

export interface AppIconProps {
  name: AppIconName;
  size?: number;
  color?: string;
  strokeWidth?: number;
  accessible?: boolean;
}

export function AppIcon({
  name,
  size = 24,
  color = "#211714",
  strokeWidth = 2,
  accessible = false,
}: AppIconProps) {
  const commonProps = { size, color, strokeWidth };

  let icon = null;

  switch (name) {
    case "search":
      icon = <SearchSvg {...commonProps} />;
      break;
    case "filter":
      icon = <FilterSvg {...commonProps} />;
      break;
    case "calendar":
      icon = <CalendarSvg {...commonProps} />;
      break;
    case "hourglass":
      icon = <HourglassSvg {...commonProps} />;
      break;
    case "exit":
      icon = <ExitSvg {...commonProps} />;
      break;
    case "bell":
      icon = <BellSvg {...commonProps} />;
      break;
    case "back":
      icon = <BackSvg {...commonProps} />;
      break;
    default:
      icon = null;
  }

  return (
    <View
      accessibilityElementsHidden={!accessible}
      importantForAccessibility={accessible ? "yes" : "no"}
    >
      {icon}
    </View>
  );
}
