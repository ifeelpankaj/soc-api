import { type PropsWithChildren, useSyncExternalStore } from "react";
import { ScrollView, StyleSheet, Text, View } from "react-native";
import { Button } from "@/components/ui/button";
import { LoadingState } from "@/components/ui/loading-state";
import { colors } from "@/theme/colors";
import { getGooglePlayURL, resolveWebAccess } from "./web-access";

const subscribe = () => () => {};
const getSnapshot = () => resolveWebAccess(
  process.env.EXPO_PUBLIC_APPLE_WEB_ENABLED,
  typeof navigator === "undefined" ? undefined : navigator,
);
const getServerSnapshot = () => resolveWebAccess(process.env.EXPO_PUBLIC_APPLE_WEB_ENABLED);

export function WebAccessGate({ children }: PropsWithChildren) {
  const access = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);

  if (access === "allowed") return children;
  if (access === "loading") return <LoadingState message="Loading Apna Gate" />;

  const downloadURL = getGooglePlayURL();

  return (
    <ScrollView contentContainerStyle={styles.screen}>
      <View style={styles.card}>
        <Text accessibilityRole="header" style={styles.title}>Apna Gate</Text>
        <Text style={styles.message}>
          {access === "android"
            ? "Use the Apna Gate Android app for resident and guard access."
            : "Use Apna Gate on your Android phone. You can also open the resident web app on an iPhone, iPad, or Mac."}
        </Text>
        {downloadURL ? (
          <Button
            title={access === "android" ? "Get the app" : "Get the Android app"}
            onPress={() => window.location.assign(downloadURL)}
          />
        ) : (
          <Text style={styles.message}>The download link is currently unavailable.</Text>
        )}
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  screen: {
    flexGrow: 1,
    justifyContent: "center",
    padding: 24,
    backgroundColor: colors.surface.screen,
  },
  card: {
    maxWidth: 480,
    alignSelf: "center",
    width: "100%",
    gap: 20,
    padding: 24,
    borderRadius: 24,
    backgroundColor: colors.surface.card,
  },
  title: { fontSize: 30, fontWeight: "700", color: colors.text.primary },
  message: { fontSize: 17, lineHeight: 26, color: colors.text.secondary },
});
