import { BillRoute } from "@/features/resident/maintenance/bill-route";
import { BillDetailsScreen } from "@/features/resident/maintenance/bill-details-screen";
export default function Route() {
  return <BillRoute screen={BillDetailsScreen} />;
}
