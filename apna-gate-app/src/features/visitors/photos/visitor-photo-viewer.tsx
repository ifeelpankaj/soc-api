import { useCallback, useState } from "react";
import { Modal, Pressable, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useFocusEffect } from "expo-router";
import { useAppSelector } from "@/redux/hooks";
import { usePhotoAccess } from "./photo-access";
import { VisitorPhoto } from "./visitor-photo";

type Entry = {
  id?: number;
  visitor?: { photo_url?: string; full_name?: string };
};

export function VisitorPhotoPreview({ entry }: { entry: Entry }) {
  const access = usePhotoAccess();
  const userId = useAppSelector((state) => state.auth.user?.id);
  const [focused, setFocused] = useState(false);
  useFocusEffect(useCallback(() => {
    setFocused(true);
    return () => setFocused(false);
  }, []));
  const identity = JSON.stringify([
    userId, access?.societyId, access?.context, entry.id,
    entry.visitor?.photo_url,
  ]);
  // Remounting resets the open state when access, photo, or route changes.
  return (
    <Preview key={identity + String(focused)} entry={entry}
      enabled={Boolean(focused && userId && access?.societyId)} />
  );
}

function Preview({ entry, enabled }: { entry: Entry; enabled: boolean }) {
  const [open, setOpen] = useState(false);
  return (
    <>
      <Pressable
        accessibilityRole="button"
        accessibilityLabel="View full visitor photo"
        disabled={!enabled}
        onPress={() => setOpen(true)}
        style={styles.fill}
      >
        <VisitorPhoto entry={entry} display="banner" variant="detail" />
        <View pointerEvents="none" style={styles.expand}>
          <Text style={styles.icon}>⛶</Text>
        </View>
      </Pressable>
      {open && enabled ? (
        <Modal visible presentationStyle="fullScreen" animationType="fade"
          onRequestClose={() => setOpen(false)}>
          <SafeAreaView style={styles.viewer}>
            <View style={styles.header}>
              <Pressable accessibilityRole="button" accessibilityLabel="Close photo viewer"
                onPress={() => setOpen(false)} style={styles.close}>
                <Text style={styles.icon}>×</Text>
              </Pressable>
            </View>
            <View style={styles.image}>
              <VisitorPhoto entry={entry} variant="original" display="fullscreen" />
            </View>
          </SafeAreaView>
        </Modal>
      ) : null}
    </>
  );
}

const styles = StyleSheet.create({
  fill: { width: "100%", height: "100%" },
  expand: {
    position: "absolute", right: 12, bottom: 12, borderRadius: 22,
    backgroundColor: "rgba(0,0,0,0.6)", width: 44, height: 44,
    alignItems: "center", justifyContent: "center",
  },
  icon: { color: "#FFFFFF", fontSize: 28 },
  viewer: { flex: 1, backgroundColor: "#000000" },
  header: { alignItems: "flex-end", paddingHorizontal: 12 },
  close: { width: 48, height: 48, alignItems: "center", justifyContent: "center" },
  image: { flex: 1 },
});
