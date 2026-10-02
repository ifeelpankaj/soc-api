export function getTimeGreeting(name?: string | null, fallbackName = "there") {
  const hour = new Date().getHours();
  const firstName = name?.trim().split(/\s+/)[0] ?? fallbackName;
  const timeGreeting =
    hour < 12 ? "Good morning" : hour < 17 ? "Good afternoon" : "Good evening";

  return `${timeGreeting}, ${firstName} 👋`;
}
