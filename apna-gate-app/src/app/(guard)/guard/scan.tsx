import { Redirect, useLocalSearchParams } from "expo-router";

import {
  firstParam,
  guardCheckInRoute,
  guardScannerRoute,
} from "@/features/guard/guard-routes";

export default function LegacyGuardScanRoute() {
  const params = useLocalSearchParams<{ token?: string | string[] }>();
  const token = firstParam(params.token)?.trim();

  if (token) {
    return <Redirect href={guardCheckInRoute({ source: "qr", token })} />;
  }

  return <Redirect href={guardScannerRoute()} />;
}
