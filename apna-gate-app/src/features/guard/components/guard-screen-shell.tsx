import { GuardSocietyGate } from "@/features/guard/components/guard-society-gate";
import { RoleScreenShell } from "@/features/shared/role-screen-shell";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import { colors } from "@/theme/colors";
import type { PropsWithChildren, ReactNode } from "react";

type GuardScreenShellProps = PropsWithChildren<{
  backgroundColor?: string;
  contentPaddingBottom?: number;
  footer?: ReactNode;
  loadingMessage?: string;
  onRefresh?: () => void;
  refreshing?: boolean;
}>;

export function GuardScreenShell({
  backgroundColor = colors.surface.screen,
  children,
  contentPaddingBottom,
  footer,
  loadingMessage = "Opening guard workspace",
  onRefresh,
  refreshing = false,
}: GuardScreenShellProps) {
  const { isLoading, isReady } = useGuardScreen();

  return (
    <RoleScreenShell
      backgroundColor={backgroundColor}
      contentPaddingBottom={contentPaddingBottom}
      footer={footer}
      gate={<GuardSocietyGate />}
      isLoading={false}
      isReady={isReady || isLoading}
      loadingMessage={loadingMessage}
      onRefresh={onRefresh}
      refreshing={refreshing}
    >
      {children}
    </RoleScreenShell>
  );
}
