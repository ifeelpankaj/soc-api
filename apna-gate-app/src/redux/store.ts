import { configureStore, createListenerMiddleware } from "@reduxjs/toolkit";
import { setupListeners } from "@reduxjs/toolkit/query";
import { Image } from "expo-image";

import appReducer, { clearWorkspaceSelection } from "./appSlice";
import authReducer, { clearAuth } from "./authSlice";
import notificationReducer, { setDevicePushToken } from "./notificationSlice";
import photoUploadsReducer from "./photoUploadsSlice";
import { baseApi, resetRefreshState } from "./queries/baseApi";

import { clearBrowserPushSession } from "@/features/web/push-session";

const sessionListener = createListenerMiddleware();
sessionListener.startListening({
  actionCreator: clearAuth,
  effect: async (_action, api) => {
    resetRefreshState();
    void clearBrowserPushSession();
    api.dispatch(setDevicePushToken(null));
    api.dispatch(clearWorkspaceSelection());
    api.dispatch(baseApi.util.resetApiState());
    await Image.clearMemoryCache().catch(() => undefined);
  },
});

export const store = configureStore({
  reducer: {
    app: appReducer,
    auth: authReducer,
    notifications: notificationReducer,
    photoUploads: photoUploadsReducer,
    [baseApi.reducerPath]: baseApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware()
      .prepend(sessionListener.middleware)
      .concat(baseApi.middleware),
});

setupListeners(store.dispatch);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
