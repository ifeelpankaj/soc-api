import type { PropsWithChildren, ReactNode } from "react";

import { ScreenShell } from "@/components/layout";
import { colors } from "@/theme/colors";

type RoleScreenShellProps = PropsWithChildren<{
  backgroundColor?: string;
  contentPaddingBottom?: number;
  footer?: ReactNode;
  gate: ReactNode;
  isLoading?: boolean;
  isReady: boolean;
  loadingMessage?: string;
  onRefresh?: () => void;
  refreshing?: boolean;
}>;

export function RoleScreenShell({
  backgroundColor = colors.surface.screen,
  children,
  contentPaddingBottom,
  footer,
  gate,
  isLoading = false,
  isReady,
  loadingMessage = "Loading workspace",
  onRefresh,
  refreshing = false,
}: RoleScreenShellProps) {
  return (
    <ScreenShell
      backgroundColor={backgroundColor}
      contentPaddingBottom={contentPaddingBottom}
      footer={footer}
      gate={gate}
      isLoading={isLoading}
      isReady={isReady}
      loadingMessage={loadingMessage}
      onRefresh={onRefresh}
      refreshing={refreshing}
    >
      {children}
    </ScreenShell>
  );
}
