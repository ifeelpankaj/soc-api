import { useEffect } from "react";
import { useLocalSearchParams, useRouter, type Href } from "expo-router";
import { Platform } from "react-native";
import { LoadingState } from "@/components/ui/loading-state";
import { useAuth } from "@/features/auth/use-auth";
import { useAppDispatch } from "@/redux/hooks";
import { setSelectedFlat } from "@/redux/appSlice";
import { saveWorkspace } from "@/features/auth/auth-storage";
import { notificationApiExtensions } from "@/lib/api/notification-api-extensions";
import { enhancedApi } from "@/lib/api/enhanced-api";
import { residentApiExtensions } from "@/lib/api/resident-api-extensions";
import { notificationRoute } from "@/features/notifications/notification-routing";
import { notificationWorkspace, validNotificationLink } from "@/features/web/notification-target";
import { clearPendingNotification } from "@/features/web/notification-handoff";

export default function NotificationOpenScreen() {
  const {id,recipient} = useLocalSearchParams<{id:string;recipient:string}>();
  const {status,user} = useAuth();
  const router = useRouter(); const dispatch = useAppDispatch();
  useEffect(() => {
    if (status === "loading") return;
    if (!validNotificationLink(id,recipient)) { router.replace("/notifications"); return; }
    if (status !== "authenticated") {
      if (Platform.OS === "web") {
        try { sessionStorage.setItem("apna_notification_open",JSON.stringify({id,recipient})); } catch { /* Inbox remains available after login. */ }
      }
      router.replace("/login"); return;
    }
    clearPendingNotification();
    let cancelled = false;
    const owned = dispatch(notificationApiExtensions.endpoints.getOwnedNotification.initiate(id,{forceRefetch:true}));
    const access = dispatch(enhancedApi.endpoints.getV1Bootstrap.initiate(undefined,{forceRefetch:true}));
    void (async () => {
      let route: Href = "/notifications";
      try {
        if (recipient !== String(user?.id)) throw new Error("Account changed");
        const [response,bootstrap] = await Promise.all([owned.unwrap(),access.unwrap()]);
        const item = response.data;
        if (!item || !bootstrap.data || !user?.id) throw new Error("Notification unavailable");
        const workspace = notificationWorkspace(item,bootstrap.data,user.id);
        if (!workspace?.flat_id || !workspace.society_id) throw new Error("Residence unavailable");
        const entryID = Number(item.data.entry_id);
        if (item.data.entry_id != null) {
          if (!Number.isSafeInteger(entryID) || entryID <= 0) throw new Error("Invalid entry");
          const entry = dispatch(residentApiExtensions.endpoints.getV1SocietiesBySocietyIdFlatsAndFlatIdVisitorEntriesAndEntryId.initiate({societyId:workspace.society_id,flatId:workspace.flat_id,entryId:entryID},{forceRefetch:true}));
          try {
            const result = await entry.unwrap();
            const current = result.data?.entry;
            if (!current || ["expired","cancelled"].includes(current.status || "")) throw new Error("Request expired");
            if (/pending|approval_requested/.test(item.type) && current.status !== "waiting_approval") throw new Error("Request already resolved");
          } finally { entry.unsubscribe(); }
        }
        const target = notificationRoute(item,"/resident/dashboard", Platform.OS === "web");
        const path = typeof target === "string" ? target.split("?")[0] : target?.pathname;
        if (!path || !/^\/resident\/(visitors|entries|invites)(\/([1-9]\d*|member|visitor))*$/.test(path)) throw new Error("Unsupported destination");
        if (cancelled) return;
        dispatch(setSelectedFlat(workspace.flat_id));
        await saveWorkspace({flatId:workspace.flat_id});
        route = target!;
      } catch { /* Owned, expired and inaccessible notifications safely fall back to the inbox. */ }
      finally { owned.unsubscribe(); access.unsubscribe(); }
      if (!cancelled) router.replace(route);
    })();
    return () => { cancelled = true; owned.unsubscribe(); access.unsubscribe(); };
  },[id,recipient,status,user?.id,dispatch,router]);
  return <LoadingState message="Opening notification…" />;
}
