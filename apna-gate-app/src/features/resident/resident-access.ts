export function canManageFlatVisitors(residence?: {
  status?: string | null;
}) {
  return residence?.status === "active";
}

export function canManageFlatMembers(residence?: {
  is_primary?: boolean;
  role?: string | null;
}) {
  return residence?.is_primary === true || residence?.role === "owner";
}
