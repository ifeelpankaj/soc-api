import { Image } from "expo-image";
import { useState } from "react";
import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";

import { colors } from "@/theme/colors";

type UserAvatarProps = {
  imageUrl?: string | null;
  name?: string | null;
  onPress?: () => void;
  showCameraBadge?: boolean;
  size?: number;
};

function getInitials(name?: string | null) {
  const trimmed = name?.trim();

  if (!trimmed) {
    return "G";
  }

  const parts = trimmed.split(/\s+/).filter(Boolean);

  if (parts.length >= 2) {
    return `${parts[0]?.[0] ?? ""}${parts[1]?.[0] ?? ""}`.toUpperCase();
  }

  return parts[0]?.slice(0, 2).toUpperCase() ?? "G";
}

export function UserAvatar({
  imageUrl,
  name,
  onPress,
  showCameraBadge = false,
  size = 44,
}: UserAvatarProps) {
  const [failedUrl, setFailedUrl] = useState<string>();
  const initials = getInitials(name);
  const fontSize = size >= 44 ? 16 : 13;
  const badgeSize = Math.max(28, Math.round(size * 0.3));

  const content =
    imageUrl && failedUrl !== imageUrl ? (
      <Image
        key={imageUrl}
        cachePolicy="none"
        onError={() => setFailedUrl(imageUrl)}
        accessibilityLabel={name ? `${name} profile photo` : "Profile photo"}
        contentFit="cover"
        source={{ uri: imageUrl }}
        style={{ borderRadius: size / 2, height: size, width: size }}
      />
    ) : (
      <View
        style={[
          styles.initials,
          { borderRadius: size / 2, height: size, width: size },
        ]}
      >
        <Text style={[styles.initialsText, { fontSize }]}>{initials}</Text>
      </View>
    );

  const ring = (
    <View style={styles.ringWrapper}>
      <View style={[styles.ring, { borderRadius: size / 2 + 2, padding: 1.5 }]}>
        {content}
      </View>
      {showCameraBadge ? (
        <View
          style={[
            styles.badge,
            {
              borderRadius: badgeSize / 2,
              height: badgeSize,
              width: badgeSize,
            },
          ]}
        >
          <SymbolView
            name={{
              ios: "camera.fill",
              android: "photo_camera",
              web: "photo_camera",
            }}
            size={Math.round(badgeSize * 0.45)}
            tintColor={colors.text.inverse}
          />
        </View>
      ) : null}
    </View>
  );

  if (!onPress) {
    return ring;
  }

  return (
    <Pressable
      accessibilityLabel={
        showCameraBadge ? "Change profile photo" : "Open profile"
      }
      accessibilityRole="button"
      style={({ pressed }) => ({ opacity: pressed ? 0.85 : 1 })}
      onPress={onPress}
    >
      {ring}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  badge: {
    alignItems: "center",
    backgroundColor: colors.guard.teal,
    borderColor: colors.surface.card,
    borderWidth: 2,
    bottom: 0,
    justifyContent: "center",
    position: "absolute",
    right: 0,
  },
  initials: {
    alignItems: "center",
    backgroundColor: colors.guard.tealSoft,
    justifyContent: "center",
  },
  initialsText: {
    color: colors.guard.teal,
    fontWeight: "700",
  },
  ring: {
    borderColor: "rgba(15, 23, 42, 0.08)",
    borderWidth: 1.5,
  },
  ringWrapper: {
    position: "relative",
  },
});
