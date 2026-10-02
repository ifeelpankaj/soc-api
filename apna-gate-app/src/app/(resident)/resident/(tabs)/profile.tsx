import { ProfileScreen } from "@/features/profile/components/profile-screen";
import { useResidentProfileModel } from "@/features/resident/hooks/use-resident-profile-model";

export default function ResidentProfileScreen() {
  const model = useResidentProfileModel();
  return <ProfileScreen {...model} />;
}
