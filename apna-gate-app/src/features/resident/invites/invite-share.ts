import { Linking, Platform, Share } from "react-native";

import type { InviteShareContent } from "@/features/resident/invites/invite-share-content";

export {
  buildMemberInviteShareContent,
  buildVisitorInviteShareContent,
  type InviteShareContent,
  type InviteShareResult,
  type MemberInviteShareInput,
  type VisitorInviteShareInput,
} from "@/features/resident/invites/invite-share-content";

export async function shareInviteNative(content: InviteShareContent) {
  await Share.share({
    message: content.message,
    title: content.title,
    url: Platform.OS === "ios" ? content.url : undefined,
  });
}

export async function shareInviteBySms(content: InviteShareContent) {
  const separator = Platform.OS === "ios" ? "&" : "?";
  const smsUrl = `sms:${separator}body=${encodeURIComponent(content.message)}`;
  const canOpen = await Linking.canOpenURL(smsUrl);
  if (!canOpen) {
    await shareInviteNative(content);
    return;
  }
  await Linking.openURL(smsUrl);
}

export async function shareInviteOnWhatsApp(content: InviteShareContent) {
  const nativeUrl = `whatsapp://send?text=${encodeURIComponent(content.message)}`;
  const webUrl = `https://wa.me/?text=${encodeURIComponent(content.message)}`;

  try {
    const canOpen = await Linking.canOpenURL(nativeUrl);
    await Linking.openURL(canOpen ? nativeUrl : webUrl);
  } catch {
    await shareInviteNative(content);
  }
}

export async function shareInviteOnTelegram(content: InviteShareContent) {
  const telegramUrl = `https://t.me/share/url?url=${encodeURIComponent(content.url)}&text=${encodeURIComponent(content.message)}`;

  try {
    const canOpen = await Linking.canOpenURL(telegramUrl);
    if (!canOpen) {
      await shareInviteNative(content);
      return;
    }

    await Linking.openURL(telegramUrl);
  } catch {
    await shareInviteNative(content);
  }
}

export async function copyInviteLink(content: InviteShareContent) {
  if (typeof navigator !== "undefined" && navigator.clipboard) {
    await navigator.clipboard.writeText(content.url);
    return true;
  }
  await Share.share({ message: content.url, title: `${content.title} link` });
  return false;
}
