import { useEffect, useRef, useState } from "react";
import { Modal, Pressable, StyleSheet, Text, View } from "react-native";
import { Image } from "expo-image";
import { SafeAreaView } from "react-native-safe-area-context";
import { AppToast, useToast } from "@/components/ui/toast";
import {
  PhotoPicker,
  photoFormData,
  releaseSelectedPhoto,
  type SelectedPhoto,
} from "@/features/visitors/photos/photo-picker";
import {
  logPhotoFailure,
  photoUploadError,
  PhotoOperationError,
} from "@/features/visitors/photos/photo-errors";
import { usePutProfilePhotoMutation } from "@/lib/api/photo-api";
import { colors } from "@/theme/colors";

export function ProfilePhotoSheet({
  visible,
  onClose,
  avatarUrl,
}: {
  visible: boolean;
  onClose: () => void;
  avatarUrl?: string;
}) {
  const [photo, setPhoto] = useState<SelectedPhoto>();
  const photoRef = useRef<SelectedPhoto | undefined>(undefined);
  const [pickerBusy, setPickerBusy] = useState(false);
  const [saving, setSaving] = useState(false);
  const uploadLock = useRef(false);
  const [upload] = usePutProfilePhotoMutation();
  const { showToast } = useToast();
  const busy = saving || pickerBusy;
  const close = () => {
    if (busy) return;
    releaseSelectedPhoto(photo);
    setPhoto(undefined);
    onClose();
  };
  useEffect(() => {
    photoRef.current = photo;
  }, [photo]);
  useEffect(
    () => () => {
      releaseSelectedPhoto(photoRef.current);
    },
    [],
  );
  const save = async () => {
    if (!photo || busy || uploadLock.current) return;
    uploadLock.current = true;
    setSaving(true);
    try {
      await upload({ body: await photoFormData(photo) }).unwrap();
      releaseSelectedPhoto(photo);
      setPhoto(undefined);
      showToast({ title: "Profile photo updated", variant: "success" });
    } catch (error) {
      const failure = photoUploadError(error);
      if (error instanceof PhotoOperationError) logPhotoFailure(failure);
      showToast({
        title: "Couldn't update photo",
        message: failure.message,
        variant: "error",
      });
    } finally {
      uploadLock.current = false;
      setSaving(false);
    }
  };
  return (
    <Modal
      visible={visible}
      presentationStyle="pageSheet"
      animationType="slide"
      onRequestClose={close}
    >
      <SafeAreaView style={styles.screen}>
        <View style={styles.header}>
          <Text style={styles.title}>Profile photo</Text>
          <Pressable
            disabled={busy}
            onPress={close}
          >
            <Text style={styles.link}>Close</Text>
          </Pressable>
        </View>
        <PhotoPicker
          profile
          value={photo}
          onChange={setPhoto}
          disabled={saving}
          onBusyChange={setPickerBusy}
          currentPhoto={
            avatarUrl ? (
              <Image
                source={{ uri: avatarUrl }}
                contentFit="cover"
                style={StyleSheet.absoluteFill}
                cachePolicy="none"
              />
            ) : undefined
          }
        />
        <Pressable
          disabled={!photo || busy}
          style={[styles.save, (!photo || busy) && { opacity: 0.5 }]}
          onPress={() => void save()}
        >
          <Text style={styles.saveText}>
            {saving ? "Saving photo..." : "Save photo"}
          </Text>
        </Pressable>
      </SafeAreaView>
      {visible ? <AppToast /> : null}
    </Modal>
  );
}
const styles = StyleSheet.create({
  screen: { flex: 1, backgroundColor: "white", padding: 18, gap: 22 },
  header: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    paddingVertical: 14,
  },
  title: { fontSize: 20, fontWeight: "700", color: colors.text.primary },
  link: { color: colors.brand.orange, padding: 8 },
  save: {
    borderRadius: 12,
    backgroundColor: colors.brand.orange,
    padding: 16,
    alignItems: "center",
    marginTop: "auto",
    marginBottom: 16,
  },
  saveText: { color: "white", fontWeight: "700", fontSize: 16 },
});
