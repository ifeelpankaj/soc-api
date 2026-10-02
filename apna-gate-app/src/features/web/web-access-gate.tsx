import type { PropsWithChildren } from "react";

// Native Android and iOS never evaluate browser detection or web configuration.
export function WebAccessGate({ children }: PropsWithChildren) {
  return children;
}
