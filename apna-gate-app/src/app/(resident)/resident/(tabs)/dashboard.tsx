import { ResidentCommandCenter } from "@/features/resident/components/resident-command-center";
import { useAndroidRootBackExit } from "@/lib/navigation/use-android-root-back-exit";

export default function ResidentDashboardScreen() {
  useAndroidRootBackExit();

  return <ResidentCommandCenter />;
}
