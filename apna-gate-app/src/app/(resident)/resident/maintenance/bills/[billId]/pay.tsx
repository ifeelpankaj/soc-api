import { BillRoute } from "@/features/resident/maintenance/bill-route";
import { PaymentScreen } from "@/features/resident/maintenance/payment-screen";
export default function Route() {
  return <BillRoute screen={PaymentScreen} />;
}
