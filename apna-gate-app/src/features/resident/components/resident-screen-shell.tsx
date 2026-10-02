import type { PropsWithChildren, ReactNode } from "react";

import { ResidentSocietyGate } from "@/features/resident/components/resident-society-gate";
import { useResident } from "@/features/resident/resident-context";
import { RoleScreenShell } from "@/features/shared/role-screen-shell";
import { colors } from "@/theme/colors";

type ResidentScreenShellProps = PropsWithChildren<{
  backgroundColor?: string;
  contentPaddingBottom?: number;
  loadingMessage?: string;
  onRefresh?: () => void;
  refreshing?: boolean;
  footer?: ReactNode;
}>;

export function ResidentScreenShell({
  backgroundColor = colors.guard.screenBg,
  children,
  contentPaddingBottom,
  footer,
  loadingMessage = "Opening resident home",
  onRefresh,
  refreshing = false,
}: ResidentScreenShellProps) {
  const { isLoading, requiresSelection, selectedResidence } = useResident();

  return (
    <RoleScreenShell
      backgroundColor={backgroundColor}
      contentPaddingBottom={contentPaddingBottom}
      footer={footer}
      gate={<ResidentSocietyGate />}
      isLoading={isLoading}
      isReady={!requiresSelection && Boolean(selectedResidence)}
      loadingMessage={loadingMessage}
      onRefresh={onRefresh}
      refreshing={refreshing}
    >
      {children}
    </RoleScreenShell>
  );
}
