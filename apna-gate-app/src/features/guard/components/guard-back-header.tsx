import type { ComponentProps } from "react";

import { ScreenBackHeader } from "@/components/layout/screen-back-header";
import { guardHomeRoute } from "@/features/guard/guard-routes";

type GuardBackHeaderProps = Omit<ComponentProps<typeof ScreenBackHeader>, "fallbackHomeRoute"> & {
  fallbackHomeRoute?: ComponentProps<typeof ScreenBackHeader>["fallbackHomeRoute"];
};

export function GuardBackHeader({ fallbackHomeRoute, ...props }: GuardBackHeaderProps) {
  return <ScreenBackHeader fallbackHomeRoute={fallbackHomeRoute ?? guardHomeRoute()} {...props} />;
}
