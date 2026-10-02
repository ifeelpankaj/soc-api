import { createContext } from "react";
export type BrowserPushState = "loading" | "install" | "setup" | "unsupported" | "default" | "denied" | "granted" | "enabling" | "error" | "enabled";
export const BrowserPushContext = createContext<{ state: BrowserPushState; enable: () => void } | null>(null);
