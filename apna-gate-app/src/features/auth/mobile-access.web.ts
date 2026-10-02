import { getMobileMemberships as allMemberships } from "./mobile-access-shared";
import type { ModelsBootstrapData } from "@/lib/api/generated-api";
export { requiresAdminPortal, getAdminMembership, getMobileResidences, formatGlobalRole, formatMembershipRole, formatStatus } from "./mobile-access-shared";
export function getMobileMemberships(bootstrap?: ModelsBootstrapData | null) {
  return allMemberships(bootstrap).filter((item) => item.role === "resident");
}
