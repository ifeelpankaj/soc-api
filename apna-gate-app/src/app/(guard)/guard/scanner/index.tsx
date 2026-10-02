import { useCallback, useRef, useState } from "react";
import {
  ActivityIndicator,
  BackHandler,
  Platform,
  Pressable,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { scanFromURLAsync } from "expo-camera";
import * as ImagePicker from "expo-image-picker";
import { SymbolView } from "expo-symbols";
import { useFocusEffect, useRouter } from "expo-router";
import { SafeAreaView } from "react-native-safe-area-context";

import { QrCamera } from "@/components/camera/qr-camera";
import { AppStatusBar } from "@/components/layout/app-status-bar";
import { AppIcon } from "@/components/icons";
import { GuardSocietyGate } from "@/features/guard/components/guard-society-gate";
import { ManualCodeSheet } from "@/features/guard/components/scanner/manual-code-sheet";
import { ScanFrameOverlay } from "@/features/guard/components/scanner/scan-frame-overlay";
import {
  guardCheckInRoute,
  guardHomeRoute,
} from "@/features/guard/guard-routes";
import { extractQrToken } from "@/features/guard/guard-utils";
import { useGuardFeedback } from "@/features/guard/hooks/use-guard-feedback";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import { useBackAction } from "@/lib/navigation/use-back-action";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

export default function GuardScannerScreen() {
  const router = useRouter();
  const handleBack = useBackAction(guardHomeRoute());
  const feedback = useGuardFeedback();
  const { isLoading, isReady, memberships, requiresSelection } =
    useGuardScreen();
  const scanLockedRef = useRef(false);
  const [flashEnabled, setFlashEnabled] = useState(false);
  const [manualVisible, setManualVisible] = useState(false);
  const [galleryLoading, setGalleryLoading] = useState(false);

  useFocusEffect(
    useCallback(() => {
      scanLockedRef.current = false;
      return () => {
        scanLockedRef.current = true;
      };
    }, []),
  );

  const handleHardwareBack = useCallback(() => {
    if (manualVisible) {
      setManualVisible(false);
      return;
    }

    handleBack();
  }, [handleBack, manualVisible]);

  useFocusEffect(
    useCallback(() => {
      if (Platform.OS !== "android") {
        return undefined;
      }

      const subscription = BackHandler.addEventListener("hardwareBackPress", () => {
        handleHardwareBack();
        return true;
      });

      return () => subscription.remove();
    }, [handleHardwareBack]),
  );

  const proceedWithScan = useCallback(
    (data: string) => {
      if (scanLockedRef.current) {
        return;
      }

      const token = extractQrToken(data);

      if (!token) {
        feedback.showError(
          "Invalid QR",
          "The scanned code does not contain a visitor token.",
          "Invalid QR code",
        );
        return;
      }

      scanLockedRef.current = true;
      router.push(guardCheckInRoute({ source: "qr", token }));
    },
    [feedback, router],
  );

  const handleGalleryPick = async () => {
    if (galleryLoading || scanLockedRef.current) {
      return;
    }

    setGalleryLoading(true);

    try {
      const permission =
        await ImagePicker.requestMediaLibraryPermissionsAsync();

      if (!permission.granted) {
        feedback.showError(
          "Permission required",
          "Allow photo library access to scan a QR code from an image.",
          "Photo library permission is required.",
        );
        return;
      }

      const result = await ImagePicker.launchImageLibraryAsync({
        allowsEditing: false,
        mediaTypes: ["images"],
        quality: 1,
      });

      if (result.canceled || !result.assets[0]?.uri) {
        return;
      }

      const scans = await scanFromURLAsync(result.assets[0].uri, ["qr"]);

      if (!scans.length || !scans[0]?.data) {
        feedback.showError(
          "No QR code found",
          "Choose a clearer image with a visitor QR code.",
          "No QR code was detected in the image.",
        );
        return;
      }

      proceedWithScan(scans[0].data);
    } catch {
      feedback.showError(
        "Scan failed",
        "Unable to read a QR code from the selected image.",
        "Please try another image.",
      );
    } finally {
      setGalleryLoading(false);
    }
  };

  if (!isLoading && (memberships.length === 0 || requiresSelection)) {
    return <GuardSocietyGate />;
  }

  return (
    <SafeAreaView
      edges={["top", "left", "right", "bottom"]}
      style={styles.screen}
    >
      <AppStatusBar />

      <View style={styles.headerCard}>
        <View style={styles.headerTopRow}>
          <Pressable
            accessibilityLabel="Go back"
            accessibilityRole="button"
            hitSlop={8}
            style={({ pressed }) => [
              styles.iconButton,
              pressed && styles.iconButtonPressed,
            ]}
            onPress={handleBack}
          >
            <AppIcon color={colors.brand.navy} name="back" size={18} />
          </Pressable>

          <Pressable
            accessibilityLabel={
              flashEnabled ? "Turn flash off" : "Turn flash on"
            }
            accessibilityRole="button"
            hitSlop={8}
            style={({ pressed }) => [
              styles.iconButton,
              flashEnabled && styles.iconButtonActive,
              pressed && styles.iconButtonPressed,
            ]}
            onPress={() => setFlashEnabled((current) => !current)}
          >
            <SymbolView
              name={{ ios: "bolt.fill", android: "flash_on", web: "flash_on" }}
              size={18}
              tintColor={flashEnabled ? colors.brand.orange : colors.brand.navy}
            />
          </Pressable>
        </View>

        <Text style={styles.title}>Scan Visitor</Text>
        <Text style={styles.subtitle}>
          Point the camera at a visitor QR code to continue check-in.
        </Text>
      </View>

      <View style={styles.cameraSection}>
        {isLoading || !isReady ? (
          <View style={styles.loadingState}>
            <ActivityIndicator color={colors.brand.orange} size="small" />
          </View>
        ) : (
          <QrCamera
            active
            enableTorch={flashEnabled}
            overlay={
              <ScanFrameOverlay
                flashEnabled={flashEnabled}
                onToggleFlash={() => setFlashEnabled((current) => !current)}
              />
            }
            onScanned={proceedWithScan}
          />
        )}
      </View>

      <View style={styles.footerCard}>
        <Text style={styles.helpTitle}>Can&apos;t scan the code?</Text>
        <Text style={styles.helpText}>
          Try adjusting the angle or distance.
        </Text>

        <View style={styles.footerActions}>
          <Pressable
            accessibilityLabel="Enter code manually"
            accessibilityRole="button"
            style={({ pressed }) => [
              styles.footerAction,
              pressed && styles.footerActionPressed,
            ]}
            onPress={() => setManualVisible(true)}
          >
            <SymbolView
              name={{ ios: "keyboard", android: "keyboard", web: "keyboard" }}
              size={18}
              tintColor={colors.brand.orange}
            />
            <Text style={styles.footerActionText}>Enter Code Manually</Text>
          </Pressable>

          <View style={styles.footerDivider} />

          <Pressable
            accessibilityLabel="Use from gallery"
            accessibilityRole="button"
            disabled={galleryLoading}
            style={({ pressed }) => [
              styles.footerAction,
              (pressed || galleryLoading) && styles.footerActionPressed,
            ]}
            onPress={() => {
              void handleGalleryPick();
            }}
          >
            {galleryLoading ? (
              <ActivityIndicator color={colors.brand.orange} size="small" />
            ) : (
              <SymbolView
                name={{ ios: "photo", android: "image", web: "image" }}
                size={18}
                tintColor={colors.brand.orange}
              />
            )}
            <Text style={styles.footerActionText}>Use from Gallery</Text>
          </Pressable>
        </View>
      </View>

      <ManualCodeSheet
        visible={manualVisible}
        onClose={() => setManualVisible(false)}
        onSubmit={(code) => {
          setManualVisible(false);
          proceedWithScan(code);
        }}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  cameraSection: {
    backgroundColor: "#0a0a0a",
    flex: 1,
    overflow: "hidden",
  },
  footerAction: {
    alignItems: "center",
    flex: 1,
    gap: spacing.sm,
    justifyContent: "center",
    minHeight: 88,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.md,
  },
  footerActionPressed: {
    opacity: 0.75,
  },
  footerActionText: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "700",
    textAlign: "center",
  },
  footerActions: {
    borderColor: colors.border.default,
    borderRadius: radius.xl,
    borderWidth: StyleSheet.hairlineWidth,
    flexDirection: "row",
    marginTop: spacing.lg,
    overflow: "hidden",
  },
  footerCard: {
    backgroundColor: colors.surface.card,
    borderTopLeftRadius: radius["2xl"],
    borderTopRightRadius: radius["2xl"],
    marginTop: -radius["2xl"],
    paddingBottom: spacing["2xl"],
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.xl,
    zIndex: 2,
  },
  footerDivider: {
    backgroundColor: colors.border.default,
    width: StyleSheet.hairlineWidth,
  },
  headerCard: {
    backgroundColor: colors.surface.card,
    borderBottomLeftRadius: radius["2xl"],
    borderBottomRightRadius: radius["2xl"],
    gap: spacing.sm,
    paddingBottom: spacing.lg,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.sm,
    zIndex: 2,
  },
  headerTopRow: {
    alignItems: "center",
    flexDirection: "row",
    justifyContent: "space-between",
  },
  helpText: {
    color: colors.guard.textMuted,
    fontSize: 14,
    marginTop: spacing.xs,
    textAlign: "center",
  },
  helpTitle: {
    color: colors.brand.navy,
    fontSize: 16,
    fontWeight: "700",
    textAlign: "center",
  },
  iconButton: {
    alignItems: "center",
    backgroundColor: colors.surface.muted,
    borderRadius: 999,
    height: 44,
    justifyContent: "center",
    width: 44,
  },
  iconButtonActive: {
    backgroundColor: colors.brand.orangeSoft,
  },
  iconButtonPressed: {
    opacity: 0.8,
  },
  loadingState: {
    alignItems: "center",
    flex: 1,
    justifyContent: "center",
  },
  screen: {
    backgroundColor: colors.surface.card,
    flex: 1,
  },
  subtitle: {
    color: colors.guard.textMuted,
    fontSize: 13,
    lineHeight: 18,
    paddingHorizontal: spacing.md,
    textAlign: "center",
  },
  title: {
    color: colors.brand.navy,
    fontSize: 22,
    fontWeight: "800",
    letterSpacing: -0.3,
    textAlign: "center",
  },
});
