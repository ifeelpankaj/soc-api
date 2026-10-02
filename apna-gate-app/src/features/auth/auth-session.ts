import type { Href } from "expo-router";

import type {
  ModelsBootstrapData,
  ModelsSocietyMemberResponse,
  ModelsUserResponse,
} from "@/lib/api/generated-api";
import { setSelectedFlat, setSelectedSociety } from "@/redux/appSlice";
import type { AppDispatch } from "@/redux/store";

import { resolveBootstrapRoute } from "./bootstrap-routing";
import { getMobileMemberships, getMobileResidences } from "./mobile-access";
import {
  loadWorkspace,
  saveWorkspace,
  type StoredWorkspace,
} from "./auth-storage";

export type RestoreSessionResult =
  | {
      status: "authenticated";
      user: ModelsUserResponse;
      bootstrap: ModelsBootstrapData;
      route: Href;
    }
  | { status: "unauthenticated" };

function isGuardMembership(membership: ModelsSocietyMemberResponse) {
  return (
    membership.status === "active" &&
    membership.role === "staff" &&
    typeof membership.society_id === "number"
  );
}

export function hydrateWorkspaceFromBootstrap(
  dispatch: AppDispatch,
  bootstrap: ModelsBootstrapData,
  workspace: StoredWorkspace,
) {
  const residences = getMobileResidences(bootstrap);
  const memberships = getMobileMemberships(bootstrap);

  if (workspace.flatId != null) {
    const residence = residences.find(
      (item) => item.flat_id === workspace.flatId,
    );
    if (residence?.flat_id != null) {
      dispatch(setSelectedFlat(residence.flat_id));
      return;
    }
  }

  if (workspace.societyId != null) {
    const membership = memberships.find(
      (item) => item.society_id === workspace.societyId,
    );
    if (membership?.society_id != null) {
      dispatch(setSelectedSociety(membership.society_id));
      return;
    }
  }

  const singleResidence = residences.length === 1 ? residences[0] : undefined;
  if (singleResidence?.flat_id != null) {
    dispatch(setSelectedFlat(singleResidence.flat_id));
    void saveWorkspace({ flatId: singleResidence.flat_id });
    return;
  }

  const guardMemberships = memberships.filter(isGuardMembership);
  if (guardMemberships.length === 1 && guardMemberships[0].society_id != null) {
    dispatch(setSelectedSociety(guardMemberships[0].society_id));
    void saveWorkspace({ societyId: guardMemberships[0].society_id });
  }
}

export async function resolveAuthenticatedRoute(
  dispatch: AppDispatch,
  bootstrap: ModelsBootstrapData,
): Promise<Href> {
  const workspace = await loadWorkspace();
  hydrateWorkspaceFromBootstrap(dispatch, bootstrap, workspace);
  return resolveBootstrapRoute(bootstrap);
}
