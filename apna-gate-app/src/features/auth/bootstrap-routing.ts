import type { Href } from "expo-router";

import type { ModelsBootstrapData } from "@/lib/api/generated-api";

import {
  getAdminMembership,
  getMobileMemberships,
  getMobileResidences,
  requiresAdminPortal,
} from "./mobile-access";

const mobileRoutes = new Set([
  "/guard/dashboard",
  "/guard/home",
  "/resident/dashboard",
  "/select-society",
]);

export function resolveBootstrapRoute(
  bootstrap?: ModelsBootstrapData | null,
): Href {
  const user = bootstrap?.user;

  if (
    requiresAdminPortal(user?.global_role) ||
    getAdminMembership(bootstrap?.memberships)
  ) {
    return "/select-society";
  }

  const defaultDashboard = bootstrap?.defaultDashboard;

  if (defaultDashboard?.path && mobileRoutes.has(defaultDashboard.path)) {
    return defaultDashboard.path as Href;
  }

  const residences = getMobileResidences(bootstrap);
  const memberships = getMobileMemberships(bootstrap);

  if (
    defaultDashboard?.kind === "select_society" ||
    residences.length + memberships.length > 1
  ) {
    return "/select-society";
  }

  if (residences.length === 1) {
    return "/resident/dashboard";
  }

  const guardMembership = memberships.find(
    (membership) => membership.role === "staff",
  );

  if (guardMembership) {
    return "/guard/home" as Href;
  }

  const residentMembership = memberships.find(
    (membership) => membership.role === "resident",
  );

  if (residentMembership) {
    return "/resident/dashboard";
  }

  return "/select-society";
}
