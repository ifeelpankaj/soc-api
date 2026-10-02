export class NotificationTracker {
  private seen = new Set<string>();
  private ready = false;
  private remember(id: string) {
    this.seen.add(id);
    if (this.seen.size > 1000) this.seen.delete(this.seen.values().next().value!);
  }
  poll(ids: string[]) {
    const fresh = this.ready ? ids.filter((id) => !this.seen.has(id)) : [];
    ids.forEach((id) => this.remember(id)); this.ready = true; return fresh;
  }
  push(id: string) {
    if (!id || this.seen.has(id)) return false;
    this.remember(id); return true;
  }
}

export function mergeNotificationPages<T extends { id: string; created_at: string }>(current: T[], incoming: T[]) {
  const items = new Map(current.map((item) => [item.id, item]));
  incoming.forEach((item) => items.set(item.id, item));
  return [...items.values()].sort((a,b) => b.created_at.localeCompare(a.created_at) || b.id.localeCompare(a.id));
}
