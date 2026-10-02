export type EntryActionState<Action extends string> = {
  entryId: number;
  action: Action;
};

export function actionForEntry<Action extends string>(
  state: EntryActionState<Action> | undefined,
  entryId?: number,
) {
  return entryId != null && state?.entryId === entryId ? state.action : undefined;
}
