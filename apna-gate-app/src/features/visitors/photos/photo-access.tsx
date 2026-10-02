import { createContext, useContext, type PropsWithChildren } from "react";
import type { PhotoContext } from "./photo-utils";

type PhotoAccess = { societyId?: number; context: PhotoContext };
const Context = createContext<PhotoAccess | null>(null);

export function PhotoAccessProvider({
  children,
  societyId,
  context,
}: PropsWithChildren<PhotoAccess>) {
  return (
    <Context.Provider value={{ societyId, context }}>
      {children}
    </Context.Provider>
  );
}

export function usePhotoAccess() {
  return useContext(Context);
}
