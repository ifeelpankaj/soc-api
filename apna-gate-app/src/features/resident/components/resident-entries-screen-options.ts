import type { ResidentVisibleEntriesSegment } from "@/features/resident/hooks/resident-entries-query";

export const RESIDENT_SEGMENT_OPTIONS: {
  label: string;
  value: ResidentVisibleEntriesSegment;
}[] = [
  { label: "Expected", value: "expected" },
  { label: "Inside", value: "inside" },
  { label: "All", value: "all" },
];
