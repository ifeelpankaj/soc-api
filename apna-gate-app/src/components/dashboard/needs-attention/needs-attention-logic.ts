import type { AttentionItem } from "./types";

export const NEEDS_ATTENTION_COLLAPSE_THRESHOLD = 3;
export const NEEDS_ATTENTION_VISIBLE_COLLAPSED = 2;

export function formatAttentionCount(count: number) {
  if (count === 1) {
    return "1 action";
  }
  return `${count} actions`;
}

export function shouldShowNeedsAttention(items: AttentionItem[]) {
  return items.length > 0;
}

export function shouldShowExpandControl(items: AttentionItem[]) {
  return items.length >= NEEDS_ATTENTION_COLLAPSE_THRESHOLD;
}

export function visibleAttentionItems(
  items: AttentionItem[],
  expanded: boolean,
): AttentionItem[] {
  if (items.length < NEEDS_ATTENTION_COLLAPSE_THRESHOLD || expanded) {
    return items;
  }
  return items.slice(0, NEEDS_ATTENTION_VISIBLE_COLLAPSED);
}

export function hiddenAttentionCount(items: AttentionItem[], expanded: boolean) {
  if (items.length < NEEDS_ATTENTION_COLLAPSE_THRESHOLD || expanded) {
    return 0;
  }
  return items.length - NEEDS_ATTENTION_VISIBLE_COLLAPSED;
}
