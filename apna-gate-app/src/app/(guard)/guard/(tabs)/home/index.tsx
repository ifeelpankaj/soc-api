import { GuardCommandCenter } from "@/features/guard/components/guard-command-center";
import { useAndroidRootBackExit } from "@/lib/navigation/use-android-root-back-exit";

export default function GuardHomeScreen() {
  useAndroidRootBackExit();

  return <GuardCommandCenter />;
}
