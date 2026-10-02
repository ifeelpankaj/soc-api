import type { Href } from "expo-router";
import type { ModelsBootstrapData } from "@/lib/api/generated-api";
import { getMobileMemberships, getMobileResidences } from "./mobile-access.web";
export function resolveBootstrapRoute(bootstrap?: ModelsBootstrapData | null): Href {
  const count = getMobileResidences(bootstrap).length + getMobileMemberships(bootstrap).length;
  return count === 1 ? "/resident/dashboard" : "/select-society";
}
