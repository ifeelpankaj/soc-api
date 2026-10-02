import { Linking, Pressable, StyleSheet, Text, View } from "react-native";
import { Image } from "expo-image";
import { SymbolView } from "expo-symbols";

import { HubCard } from "@/features/hub/components/ui";
import { useHubAttachmentQuery } from "@/lib/api/hub-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

function HubAttachmentItem({
  societyId,
  uploadId,
}: {
  societyId: number;
  uploadId: number;
}) {
  const query = useHubAttachmentQuery({ societyId, uploadId });
  const attachment = query.data;
  const isPdf = attachment?.mime_type?.includes("pdf");

  if (query.isLoading) {
    return <View style={[styles.tile, styles.placeholder]} />;
  }
  if (query.isError || !attachment) {
    return (
      <View style={[styles.tile, styles.placeholder]}>
        <Text style={styles.unavailable}>Unavailable</Text>
      </View>
    );
  }

  if (isPdf) {
    return (
      <Pressable
        accessibilityRole="link"
        style={({ pressed }) => [styles.pdfTile, pressed && styles.pressed]}
        onPress={() => void Linking.openURL(attachment.url)}
      >
        <SymbolView
          name={{ ios: "doc.fill", android: "description", web: "description" }}
          size={28}
          tintColor={colors.brand.navy}
        />
        <Text numberOfLines={2} style={styles.pdfName}>
          {attachment.filename ?? "Document"}
        </Text>
      </Pressable>
    );
  }

  return (
    <Image
      contentFit="cover"
      source={{ uri: attachment.url }}
      style={styles.tile}
      transition={200}
    />
  );
}

type HubAttachmentGalleryProps = {
  societyId: number;
  attachmentIds?: number[];
};

export function HubAttachmentGallery({
  societyId,
  attachmentIds,
}: HubAttachmentGalleryProps) {
  const ids = attachmentIds?.filter((id) => id > 0) ?? [];
  if (ids.length === 0) {
    return null;
  }
  return (
    <HubCard style={styles.gallery}>
      <View style={styles.grid}>
        {ids.map((id) => (
          <HubAttachmentItem key={id} societyId={societyId} uploadId={id} />
        ))}
      </View>
    </HubCard>
  );
}

const styles = StyleSheet.create({
  gallery: {
    marginTop: spacing.md,
    padding: spacing.sm,
  },
  grid: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.sm,
  },
  pdfName: {
    color: colors.brand.navy,
    fontSize: 12,
    fontWeight: "600",
    marginTop: spacing.xs,
    textAlign: "center",
  },
  pdfTile: {
    alignItems: "center",
    backgroundColor: colors.surface.muted,
    borderColor: colors.dashboard.cardBorder,
    borderRadius: radius.md,
    borderWidth: StyleSheet.hairlineWidth,
    height: 120,
    justifyContent: "center",
    padding: spacing.sm,
    width: "48%",
  },
  placeholder: {
    backgroundColor: colors.surface.secondary,
  },
  pressed: {
    opacity: 0.92,
  },
  tile: {
    borderRadius: radius.md,
    height: 120,
    width: "48%",
  },
  unavailable: {
    color: colors.text.muted,
    fontSize: 12,
  },
});
