// Serializes writes and invalidates work started before sign-out.
export class SessionWriteQueue {
  private generation = 0;
  private tail: Promise<unknown> = Promise.resolve();
  current() {
    return this.generation;
  }
  private enqueue(work: () => Promise<void>) {
    const result = this.tail.catch(() => undefined).then(work);
    this.tail = result;
    return result;
  }
  save(generation: number, work: () => Promise<void>) {
    return this.enqueue(async () => {
      if (generation === this.generation) await work();
    });
  }
  clear(work: () => Promise<void>) {
    this.generation += 1;
    return this.enqueue(work);
  }
}
