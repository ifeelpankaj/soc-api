import type {
  ModelsBootstrapData,
  ModelsFlatResidentResponse,
  ModelsGlobalRole,
  ModelsSocietyMemberResponse,
  ModelsSocietyMemberRole,
} from "@/lib/api/generated-api";

const ADMIN_MEMBERSHIP_ROLES = new Set<ModelsSocietyMemberRole>(["owner", "admin"]);
const MOBILE_MEMBERSHIP_ROLES = new Set<ModelsSocietyMemberRole>(["staff", "resident"]);

export function requiresAdminPortal(globalRole?: ModelsGlobalRole) {
  return !!globalRole && globalRole !== "user";
}

export function getAdminMembership(
  memberships?: ModelsSocietyMemberResponse[],
) {
  return memberships?.find(
    (membership) =>
      membership.status === "active" &&
      membership.role &&
      ADMIN_MEMBERSHIP_ROLES.has(membership.role),
  );
}

export function getMobileResidences(
  bootstrap?: ModelsBootstrapData | null,
): ModelsFlatResidentResponse[] {
  return (
    bootstrap?.residences?.filter((residence) => residence.status === "active") ??
    []
  );
}

function getResidenceSocietyIds(residences: ModelsFlatResidentResponse[]) {
  return new Set(
    residences
      .map((residence) => residence.society_id)
      .filter((societyId): societyId is number => societyId != null),
  );
}

export function getMobileMemberships(
  bootstrap?: ModelsBootstrapData | null,
): ModelsSocietyMemberResponse[] {
  const residences = getMobileResidences(bootstrap);
  const residenceSocietyIds = getResidenceSocietyIds(residences);

  return (
    bootstrap?.memberships?.filter((membership) => {
      if (membership.status !== "active" || !membership.role) {
        return false;
      }

      if (!MOBILE_MEMBERSHIP_ROLES.has(membership.role)) {
        return false;
      }

      // Resident society membership is redundant when flat-level residence exists.
      if (
        membership.role === "resident" &&
        membership.society_id != null &&
        residenceSocietyIds.has(membership.society_id)
      ) {
        return false;
      }

      return true;
    }) ?? []
  );
}

export function formatGlobalRole(globalRole?: ModelsGlobalRole) {
  switch (globalRole) {
    case "developer":
      return "Developer";
    case "super_admin":
      return "Super admin";
    default:
      return "Admin";
  }
}

export function formatMembershipRole(role?: ModelsSocietyMemberRole) {
  switch (role) {
    case "owner":
      return "Society owner";
    case "admin":
      return "Society admin";
    case "staff":
      return "Guard";
    case "resident":
      return "Resident";
    default:
      return "Member";
  }
}

export function formatStatus(status?: string | null) {
  if (!status) {
    return "Unknown";
  }

  return status
    .split("_")
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}
