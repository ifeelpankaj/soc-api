import { useRef, useState, type ReactNode } from "react";
import {
  ActivityIndicator,
  Linking,
  Platform,
  Pressable,
  StyleSheet,
  Text,
  View,
} from "react-native";
import * as ImagePicker from "expo-image-picker";
import { ImageManipulator, SaveFormat } from "expo-image-manipulator";
import { File } from "expo-file-system";
import { Image } from "expo-image";
import { SymbolView } from "expo-symbols";
import { useToast } from "@/components/ui/toast";
import { colors } from "@/theme/colors";
import {
  MAX_PHOTO_BYTES,
  MAX_PHOTO_PIXELS,
  photoPreparationNeeded,
} from "./photo-utils";
import { PhotoOperationError } from "./photo-errors";

export type SelectedPhoto = {
  uri: string;
  name: string;
  type: string;
  file?: Blob;
  managed?: boolean;
};

export function releaseSelectedPhoto(photo?: SelectedPhoto) {
  if (Platform.OS === "web" || !photo?.managed) return;
  try {
    const file = new File(photo.uri);
    if (file.exists) file.delete();
  } catch {
    console.warn("Photo cache cleanup failed");
  }
}

export async function photoFormData(photo: SelectedPhoto) {
  try {
    const body = new FormData();
    if (Platform.OS === "web") {
      body.append(
        "file",
        photo.file ?? (await (await fetch(photo.uri)).blob()),
        photo.name,
      );
    } else {
      const file = new File(photo.uri);
      if (!file.exists || !file.size)
        throw new Error("Prepared photo is missing or unreadable");
      // Expo's native fetch converter accepts File/Blob parts with byte access.
      // React Native's legacy { uri, name, type } descriptor is unsupported by
      // that converter and fails locally before an HTTP request is dispatched.
      body.append("file", file, photo.name);
    }
    return body;
  } catch (error) {
    throw new PhotoOperationError(
      "build_form_data",
      error instanceof Error ? error.message : String(error),
      "The selected photo could not be read. Please choose it again.",
    );
  }
}

async function preparePhoto(
  asset: ImagePicker.ImagePickerAsset,
  profile: boolean,
): Promise<SelectedPhoto> {
  const mime =
    asset.mimeType?.toLowerCase() ??
    (/\.png$/i.test(asset.uri)
      ? "image/png"
      : /\.jpe?g$/i.test(asset.uri)
        ? "image/jpeg"
        : "");
  let size = asset.fileSize ?? asset.file?.size ?? 0;
  try {
    if (!size && Platform.OS !== "web") size = new File(asset.uri).size;
  } catch (error) {
    throw new PhotoOperationError(
      "prepare_file",
      error instanceof Error ? error.message : String(error),
      "The selected photo could not be read. Please choose it again.",
    );
  }
  if (asset.width <= 0 || asset.height <= 0)
    throw new PhotoOperationError(
      "prepare_file",
      "Picker returned invalid image dimensions",
      "The selected photo is not a valid image. Please choose another.",
    );
  if (asset.width * asset.height > MAX_PHOTO_PIXELS)
    throw new PhotoOperationError(
      "prepare_file",
      "Selected photo exceeds the 25-megapixel limit",
      "Choose a photo smaller than 25 megapixels.",
    );
  if (
    Platform.OS === "web" &&
    profile &&
    !photoPreparationNeeded(mime, size, asset.width, asset.height)
  ) {
    return {
      uri: asset.uri,
      type: mime,
      name: mime === "image/png" ? "visitor.png" : "visitor.jpg",
      file: asset.file,
    };
  }
  const context = ImageManipulator.manipulate(asset.uri);
  let rendered: Awaited<ReturnType<typeof context.renderAsync>> | undefined;
  let preparedFile: File | undefined;
  try {
    if (Math.max(asset.width, asset.height) > 1600) {
      context.resize(
        asset.width >= asset.height ? { width: 1600 } : { height: 1600 },
      );
    }
    rendered = await context.renderAsync();
    const result = await rendered.saveAsync({
      format: SaveFormat.JPEG,
      compress: profile ? 0.85 : 0.8,
    });
    const blob =
      Platform.OS === "web"
        ? await (await fetch(result.uri)).blob()
        : undefined;
    preparedFile = Platform.OS === "web" ? undefined : new File(result.uri);
    const outputSize = blob?.size ?? preparedFile?.size ?? 0;
    if (!outputSize || (preparedFile && !preparedFile.exists))
      throw new Error("Prepared photo is missing or unreadable");
    if (outputSize > MAX_PHOTO_BYTES)
      throw new Error("Choose a smaller photo (maximum 5 MiB).");
    return {
      uri: result.uri,
      name: "visitor.jpg",
      type: "image/jpeg",
      file: blob,
      managed: Platform.OS !== "web",
    };
  } catch (error) {
    if (preparedFile?.exists) preparedFile.delete();
    throw new PhotoOperationError(
      "prepare_file",
      error instanceof Error ? error.message : String(error),
      "The selected photo could not be prepared. Please choose it again.",
    );
  } finally {
    rendered?.release();
    context.release();
  }
}

export function PhotoPicker({
  value,
  onChange,
  disabled,
  onBusyChange,
  name,
  purpose,
  currentPhoto,
  profile = false,
}: {
  value?: SelectedPhoto;
  onChange: (photo: SelectedPhoto | undefined) => void;
  disabled?: boolean;
  onBusyChange?: (busy: boolean) => void;
  name?: string;
  purpose?: string;
  currentPhoto?: ReactNode;
  profile?: boolean;
}) {
  const { showToast } = useToast();
  const [busy, setBusy] = useState(false);
  const lock = useRef(false);
  const [settings, setSettings] = useState(false);
  const pick = async (camera: boolean) => {
    if (disabled || lock.current) return;
    lock.current = true;
    setBusy(true);
    onBusyChange?.(true);
    setSettings(false);
    try {
      if (camera) {
        const permission = await ImagePicker.requestCameraPermissionsAsync();
        if (!permission.granted) {
          setSettings(!permission.canAskAgain);
          throw new Error(
            "Camera permission is needed to take a visitor photo.",
          );
        }
      }
      const options: ImagePicker.ImagePickerOptions = {
        mediaTypes: ["images"],
        allowsEditing: false,
        allowsMultipleSelection: false,
      };
      const result = camera
        ? await ImagePicker.launchCameraAsync(options)
        : await ImagePicker.launchImageLibraryAsync(options);
      if (!result.canceled && result.assets[0]) {
        const next = await preparePhoto(result.assets[0], profile);
        releaseSelectedPhoto(value);
        onChange(next);
      }
    } catch (reason) {
      const userMessage =
        reason instanceof PhotoOperationError
          ? reason.userMessage
          : reason instanceof Error
            ? reason.message
            : "Please try another photo.";
      if (reason instanceof PhotoOperationError)
        console.warn("Photo operation failed", {
          stage: reason.stage,
        });
      showToast({
        title: "Couldn't select photo",
        message: userMessage,
        variant: "error",
      });
    } finally {
      lock.current = false;
      setBusy(false);
      onBusyChange?.(false);
    }
  };
  return (
    <View style={styles.container}>
      <View style={styles.preview}>
        {value ? (
          <Image
            source={{ uri: value.uri }}
            style={StyleSheet.absoluteFill}
            contentFit="cover"
            cachePolicy="none"
          />
        ) : (
          (currentPhoto ?? (
            <SymbolView
              name={{ ios: "person.fill", android: "person", web: "person" }}
              size={48}
              tintColor="#AAB1BC"
            />
          ))
        )}
        {!value && !currentPhoto ? (
          <View style={styles.plus}>
            <Text style={styles.plusText}>+</Text>
          </View>
        ) : null}
      </View>
      <View style={styles.copy}>
        <Text numberOfLines={1} style={styles.label}>
          {name || (profile ? "Profile photo" : "Add visitor photo")}
        </Text>
        {purpose ? (
          <Text style={styles.purpose}>{purpose}</Text>
        ) : (
          <Text style={styles.hint}>
            {profile
              ? "Choose a photo for your profile"
              : "Optional — helps identify the visitor"}
          </Text>
        )}
        <View style={styles.actions}>
          <Pressable
            accessibilityRole="button"
            disabled={disabled || busy}
            style={styles.button}
            onPress={() => void pick(true)}
          >
            <SymbolView
              name={{
                ios: "camera.fill",
                android: "photo_camera",
                web: "photo_camera",
              }}
              size={15}
              tintColor={colors.brand.orange}
            />
            <Text style={styles.action}>Use camera</Text>
          </Pressable>
          <Pressable
            accessibilityRole="button"
            disabled={disabled || busy}
            style={styles.button}
            onPress={() => void pick(false)}
          >
            <SymbolView
              name={{ ios: "photo", android: "image", web: "image" }}
              size={15}
              tintColor={colors.brand.orange}
            />
            <Text style={styles.action}>
              {name || profile ? "Choose photo" : "Gallery"}
            </Text>
          </Pressable>
        </View>
        {purpose ? <Text style={styles.hint}>Photo is optional</Text> : null}
        {value ? (
          <Pressable
            disabled={disabled || busy}
            onPress={() => {
              releaseSelectedPhoto(value);
              onChange(undefined);
            }}
          >
            <Text style={styles.hint}>Discard selection</Text>
          </Pressable>
        ) : null}
        {busy ? (
          <ActivityIndicator size="small" color={colors.brand.orange} />
        ) : null}
        {settings ? (
          <Pressable onPress={() => void Linking.openSettings()}>
            <Text style={styles.action}>Open settings</Text>
          </Pressable>
        ) : null}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: "row",
    alignItems: "center",
    gap: 14,
    padding: 10,
    borderWidth: 1,
    borderColor: "#E7EBF1",
    borderRadius: 14,
    backgroundColor: "#FFFFFF",
    shadowColor: "#15213A",
    shadowOpacity: 0.05,
    shadowRadius: 5,
    shadowOffset: { width: 0, height: 2 },
    elevation: 2,
  },
  preview: {
    width: 74,
    height: 78,
    borderRadius: 12,
    overflow: "hidden",
    backgroundColor: "#EEF1F6",
    alignItems: "center",
    justifyContent: "center",
  },
  copy: { flex: 1, gap: 5 },
  label: { fontSize: 13, fontWeight: "700", color: "#13213B" },
  hint: { fontSize: 10, color: "#7A879D", paddingVertical: 2 },
  purpose: {
    color: colors.brand.orange,
    backgroundColor: "#FFF2E9",
    alignSelf: "flex-start",
    borderRadius: 12,
    paddingHorizontal: 10,
    paddingVertical: 3,
    fontSize: 11,
  },
  actions: { flexDirection: "row", flexWrap: "wrap", gap: 7 },
  button: {
    flexDirection: "row",
    alignItems: "center",
    gap: 6,
    paddingVertical: 7,
    paddingHorizontal: 9,
    borderWidth: 1,
    borderColor: "#FFB38A",
    borderRadius: 6,
    backgroundColor: "#FFFDFC",
  },
  action: { color: colors.brand.orange, fontSize: 11, fontWeight: "600" },
  plus: {
    position: "absolute",
    bottom: 5,
    right: 4,
    width: 21,
    height: 21,
    borderRadius: 11,
    backgroundColor: colors.brand.orange,
    borderWidth: 2,
    borderColor: "white",
    alignItems: "center",
    justifyContent: "center",
  },
  plusText: { color: "white", fontWeight: "700", fontSize: 16, lineHeight: 17 },
});
