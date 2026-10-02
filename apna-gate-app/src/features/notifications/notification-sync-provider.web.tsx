import { type PropsWithChildren, useCallback, useEffect, useRef, useState, useSyncExternalStore } from "react";
import { useAuth } from "@/features/auth/use-auth";
import { useToast } from "@/components/ui";
import { getMobileMemberships, getMobileResidences } from "@/features/auth/mobile-access";
import { useGetV1BootstrapQuery, usePostV1MeDeviceTokensMutation, useDeleteV1MeDeviceTokensMutation } from "@/lib/api/generated-api";
import { useGetV1MeNotificationsQuery, useGetV1MeNotificationsUnreadCountQuery } from "@/lib/api/notification-api-extensions";
import { enhancedApi, invalidateVisitorNotificationTags } from "@/lib/api/enhanced-api";
import { useAppDispatch } from "@/redux/hooks";
import { setDevicePushToken } from "@/redux/notificationSlice";
import { BrowserPushContext, type BrowserPushState } from "@/features/web/notification-context";
import { browserPushAvailability, browserPushGeneration, listenBrowserPush, registerBrowserPush } from "@/features/web/push-session.web";
import { NotificationTracker } from "@/features/web/notification-tracker";
import { hasWebPushConfig } from "@/features/web/firebase-config";
import { needsWebPushInstallation } from "@/features/web/device";
import { dismissAlertBanner } from "@/features/web/notification-banner-session";

function subscribeVisible(callback: () => void) {
  document.addEventListener("visibilitychange", callback);
  window.addEventListener("online", callback); window.addEventListener("offline", callback);
  return () => { document.removeEventListener("visibilitychange", callback); window.removeEventListener("online", callback); window.removeEventListener("offline", callback); };
}
const visible = () => document.visibilityState === "visible" && navigator.onLine;

export function NotificationSyncProvider({ children }: PropsWithChildren) {
  const { status, user } = useAuth();
  const dispatch = useAppDispatch();
  const { showToast } = useToast();
  const active = useSyncExternalStore(subscribeVisible, visible, () => false);
  const bootstrap = useGetV1BootstrapQuery(undefined, { skip: status !== "authenticated" });
  const eligible = status === "authenticated" && getMobileResidences(bootstrap.data?.data).length + getMobileMemberships(bootstrap.data?.data).length > 0;
  const enabled = eligible && active;
  // Re-subscribing on return/reconnect immediately refreshes even cached data.
  const list = useGetV1MeNotificationsQuery({limit:20}, {skip:!enabled, pollingInterval:enabled ? 30000 : 0, refetchOnMountOrArgChange:true});
  const count = useGetV1MeNotificationsUnreadCountQuery(undefined, {skip:!enabled, pollingInterval:enabled ? 30000 : 0, refetchOnMountOrArgChange:true});
  const tracker = useRef(new NotificationTracker());
  const [pushState,setPushState] = useState<BrowserPushState>("loading");
  const [register] = usePostV1MeDeviceTokensMutation();
  const [unregister] = useDeleteV1MeDeviceTokensMutation();
  const busy = useRef(false);
  const liveUser = useRef(user?.id);
  useEffect(() => { liveUser.current = eligible ? user?.id : undefined; },[eligible,user?.id]);

  const activate = useCallback(async () => {
    const id = liveUser.current;
    if (!id || busy.current) return;
    busy.current = true; setPushState("enabling");
    const epoch = browserPushGeneration();
    try {
      const token = await registerBrowserPush(id,epoch, async (token,installationID) => {
        await register({modelsRegisterDeviceTokenRequest:{token,device_id:installationID,platform:"web"}}).unwrap();
      }, async (token) => { await unregister({modelsUnregisterDeviceTokenRequest:{token}}).unwrap(); });
      if (liveUser.current === id && epoch === browserPushGeneration() && token) { dispatch(setDevicePushToken(token)); setPushState("enabled"); dismissAlertBanner(id); }
    } catch { if (liveUser.current === id) setPushState("error"); }
    finally { busy.current = false; }
  }, [dispatch,register,unregister]);

  const enable = useCallback(() => {
    if (busy.current || !liveUser.current) return;
    if (needsWebPushInstallation()) { setPushState("install"); return; }
    if (typeof Notification === "undefined") { setPushState("unsupported"); return; }
    if (!hasWebPushConfig()) { setPushState("setup"); return; }
    if (Notification.permission === "denied") { setPushState("denied"); return; }
    // Permission must be requested synchronously from this user gesture.
    if (Notification.permission === "default") {
      busy.current = true; setPushState("enabling");
      void Notification.requestPermission().then((permission) => {
        busy.current = false;
        if (permission === "granted") void activate(); else setPushState(permission);
      }).catch(() => { busy.current = false; setPushState("error"); });
    } else void activate();
  },[activate]);

  useEffect(() => { tracker.current = new NotificationTracker(); },[user?.id,status]);
  useEffect(() => {
    if (!eligible || !active) return;
    let cancelled = false;
    void browserPushAvailability().then((state) => {
      if (cancelled) return;
      setPushState(state);
      if (state === "granted") void activate();
    }).catch(() => { if (!cancelled) setPushState("error"); });
    return () => { cancelled = true; };
  },[eligible,active,user?.id,activate]);

  useEffect(() => {
    if (!enabled || !list.data?.data?.items) return;
    const fresh = tracker.current.poll(list.data.data.items.map((item) => item.id));
    if (!fresh.length) return;
    const updates = list.data.data.items.filter((item) => fresh.includes(item.id));
    updates.forEach((item) => dispatch(enhancedApi.util.invalidateTags(invalidateVisitorNotificationTags(item.type))));
    const unread = updates.filter((item) => !item.read_at);
    if (unread.length) showToast({title: unread.length === 1 ? unread[0].title : "New notifications",message:unread.length === 1 ? unread[0].body : `${unread.length} updates in your inbox.`,variant:"info"});
  },[enabled,list.data,dispatch,showToast]);

  const refresh = useRef(() => {});
  useEffect(() => {
    refresh.current = () => { if (enabled) { void list.refetch(); void count.refetch(); } };
  },[enabled,list,count]);
  useEffect(() => {
    if (!eligible) return;
    let disposed = false; let remove = () => {};
    void listenBrowserPush((payload) => {
      if (disposed || payload.data?.user_id !== String(liveUser.current)) return;
      refresh.current();
      if (payload.data?.notification_id && tracker.current.push(payload.data.notification_id)) {
        dispatch(enhancedApi.util.invalidateTags(invalidateVisitorNotificationTags(payload.data.type)));
        if (visible()) showToast({title:payload.data.title || "Apna Gate",message:payload.data.body || "New update in your inbox.",variant:"info"});
      }
    }).then((unsubscribe) => { if (disposed) unsubscribe(); else remove = unsubscribe; }).catch(() => undefined);
    return () => { disposed = true; remove(); };
  },[eligible,user?.id,dispatch,showToast]);

  return <BrowserPushContext.Provider value={eligible ? {state:pushState,enable} : null}>{children}</BrowserPushContext.Provider>;
}
