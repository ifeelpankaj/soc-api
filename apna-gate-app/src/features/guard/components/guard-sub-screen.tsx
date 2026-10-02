import type { PropsWithChildren, ReactNode } from "react";

import { GuardSocietyGate } from "@/features/guard/components/guard-society-gate";
import { guardHomeRoute } from "@/features/guard/guard-routes";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import { RoleSubScreen } from "@/features/shared/role-sub-screen";

type GuardSubScreenProps = PropsWithChildren<{
  footer?: ReactNode;
  headerExtra?: ReactNode;
  headerTrailing?: ReactNode;
  title: string;
}>;

export function GuardSubScreen({ children, footer, headerExtra, headerTrailing, title }: GuardSubScreenProps) {
  const { isLoading, isReady, memberships, requiresSelection } = useGuardScreen();

  return (
    <RoleSubScreen
      footer={footer}
      fallbackHomeRoute={guardHomeRoute()}
      gate={<GuardSocietyGate />}
      headerExtra={headerExtra}
      headerTrailing={headerTrailing}
      isLoading={isLoading}
      isReady={isReady}
      showGate={!isLoading && (memberships.length === 0 || requiresSelection)}
      title={title}
    >
      {children}
    </RoleSubScreen>
  );
}
