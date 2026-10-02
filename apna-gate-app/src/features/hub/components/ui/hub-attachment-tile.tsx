import { ActivityIndicator, Pressable, StyleSheet, Text, View } from "react-native";
import { Image } from "expo-image";
import { SymbolView } from "expo-symbols";

import { hubTheme } from "@/features/hub/hub-theme";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type HubAttachmentTileProps = {
  uri: string;
  uploading?: boolean;
  error?: string;
  onRemove?: () => void;
};

export function HubAttachmentTile({
  uri,
  uploading,
  error,
  onRemove,
}: HubAttachmentTileProps) {
  return (
    <View style={styles.tile}>
      <Image contentFit="cover" source={{ uri }} style={styles.image} />
      {uploading ? (
        <View style={styles.overlay}>
          <ActivityIndicator color={colors.text.inverse} />
        </View>
      ) : null}
      {error ? (
        <View style={[styles.overlay, styles.errorOverlay]}>
          <Text style={styles.errorText} numberOfLines={2}>
            {error}
          </Text>
        </View>
      ) : null}
      {onRemove && !uploading ? (
        <Pressable
          accessibilityLabel="Remove photo"
          accessibilityRole="button"
          hitSlop={8}
          style={styles.remove}
          onPress={onRemove}
        >
          <SymbolView
            name={{ ios: "xmark.circle.fill", android: "cancel", web: "cancel" }}
            size={22}
            tintColor={colors.text.inverse}
          />
        </Pressable>
      ) : null}
    </View>
  );
}

export function HubAddAttachmentTile({ onPress, disabled }: { onPress: () => void; disabled?: boolean }) {
  return (
    <Pressable
      accessibilityLabel="Add photos"
      accessibilityRole="button"
      disabled={disabled}
      style={({ pressed }) => [
        styles.addTile,
        disabled && styles.addDisabled,
        pressed && !disabled && styles.pressed,
      ]}
      onPress={onPress}
    >
      <SymbolView
        name={{ ios: "photo.badge.plus", android: "add_photo_alternate", web: "add_photo_alternate" }}
        size={28}
        tintColor={colors.text.muted}
      />
      <Text style={styles.addLabel}>Add photos</Text>
    </Pressable>
  );
}

const size = 96;

const styles = StyleSheet.create({
  addDisabled: {
    opacity: 0.5,
  },
  addLabel: {
    color: colors.text.muted,
    fontSize: 11,
    fontWeight: "600",
    marginTop: spacing.xs,
  },
  addTile: {
    alignItems: "center",
    backgroundColor: colors.surface.muted,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.md,
    borderStyle: "dashed",
    borderWidth: 1,
    height: size,
    justifyContent: "center",
    width: size,
  },
  errorOverlay: {
    backgroundColor: "rgba(220, 38, 38, 0.75)",
  },
  errorText: {
    color: colors.text.inverse,
    fontSize: 10,
    padding: spacing.xs,
    textAlign: "center",
  },
  image: {
    borderRadius: radius.md,
    height: size,
    width: size,
  },
  overlay: {
    ...StyleSheet.absoluteFill,
    alignItems: "center",
    backgroundColor: "rgba(16, 29, 54, 0.45)",
    borderRadius: radius.md,
    justifyContent: "center",
  },
  pressed: {
    opacity: hubTheme.pressOpacity,
  },
  remove: {
    position: "absolute",
    right: 4,
    top: 4,
  },
  tile: {
    borderRadius: radius.md,
    height: size,
    overflow: "hidden",
    width: size,
  },
});
