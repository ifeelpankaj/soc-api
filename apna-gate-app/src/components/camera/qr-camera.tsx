import { CameraView, useCameraPermissions } from "expo-camera";
import type { ReactNode } from "react";
import { useRef } from "react";
import { Pressable, StyleSheet, Text, View } from "react-native";

import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type QrCameraProps = {
  active?: boolean;
  enableTorch?: boolean;
  overlay?: ReactNode;
  onScanned?: (data: string) => void;
};

export function QrCamera({
  active = true,
  enableTorch = false,
  overlay,
  onScanned,
}: QrCameraProps) {
  const [permission, requestPermission] = useCameraPermissions();
  const hasScannedRef = useRef(false);

  if (!permission) {
    return (
      <View style={styles.permissionContainer}>
        <Text style={styles.permissionText}>Checking camera permission</Text>
      </View>
    );
  }

  if (!permission.granted) {
    return (
      <View style={[styles.permissionContainer, styles.permissionDenied]}>
        <Text style={styles.permissionText}>Camera access is required to scan QR codes.</Text>
        <Pressable style={styles.allowButton} onPress={requestPermission}>
          <Text style={styles.allowButtonText}>Allow camera</Text>
        </Pressable>
      </View>
    );
  }

  return (
    <View style={styles.wrapper}>
      <CameraView
        active={active}
        barcodeScannerSettings={{
          barcodeTypes: ["qr"],
        }}
        enableTorch={enableTorch}
        facing="back"
        style={styles.camera}
        onBarcodeScanned={({ data }) => {
          if (hasScannedRef.current || !active) {
            return;
          }

          hasScannedRef.current = true;
          onScanned?.(data);
          setTimeout(() => {
            hasScannedRef.current = false;
          }, 1800);
        }}
      />
      {overlay ? <View pointerEvents="box-none" style={styles.overlay}>{overlay}</View> : null}
    </View>
  );
}

const styles = StyleSheet.create({
  allowButton: {
    backgroundColor: colors.brand.orange,
    borderRadius: radius.md,
    paddingHorizontal: layout.footerPaddingHorizontal,
    paddingVertical: spacing.md,
  },
  allowButtonText: {
    ...typography.button,
    color: colors.text.inverse,
  },
  camera: {
    flex: 1,
  },
  overlay: {
    ...StyleSheet.absoluteFill,
  },
  permissionContainer: {
    alignItems: "center",
    backgroundColor: "#0a0a0a",
    flex: 1,
    justifyContent: "center",
    paddingHorizontal: spacing["2xl"],
  },
  permissionDenied: {
    gap: spacing.lg,
  },
  permissionText: {
    ...typography.body,
    color: colors.text.inverse,
    textAlign: "center",
  },
  wrapper: {
    flex: 1,
    overflow: "hidden",
  },
});
