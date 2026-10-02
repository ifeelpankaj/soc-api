import { useGuardProfileModel } from "@/features/guard/hooks/use-guard-profile-model";
import { ProfileScreen } from "@/features/profile";

export default function GuardProfileScreen() {
  const model = useGuardProfileModel();
  return <ProfileScreen {...model} />;
}
