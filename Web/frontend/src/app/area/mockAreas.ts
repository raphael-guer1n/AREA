export type MockArea = {
  id: string;
  name: string;
  action: { label: string; colorClass?: string };
  reaction: { label: string; colorClass?: string };
  active: boolean;
  gradient?: { from: string; to: string };
  lastRun?: string;
};

const gradients: Array<{ from: string; to: string }> = [
  { from: "#002642", to: "#0b3c5d" },
  { from: "#840032", to: "#a33a60" },
  { from: "#e59500", to: "#f2b344" },
  { from: "#5B834D", to: "#68915aff" },
  { from: "#02040f", to: "#1b2640" },
];

export const mockAreas: MockArea[] = [
  {
    id: "1",
    name: "Gmail vers Slack",
    action: { label: "Gmail", colorClass: "bg-white/20 text-white" },
    reaction: { label: "Slack", colorClass: "bg-white/20 text-white" },
    active: true,
    gradient: gradients[0],
    lastRun: "Il y a 5 min",
  },
  {
    id: "2",
    name: "Sauvegarder les photos Instagram sur Drive",
    action: { label: "Instagram", colorClass: "bg-white/20 text-white" },
    reaction: { label: "Drive", colorClass: "bg-white/20 text-white" },
    active: true,
    gradient: gradients[1],
    lastRun: "Il y a 18 min",
  },
  {
    id: "3",
    name: "Créer une tâche Jira depuis GitHub",
    action: { label: "GitHub", colorClass: "bg-white/20 text-white" },
    reaction: { label: "Jira", colorClass: "bg-white/20 text-white" },
    active: false,
    gradient: gradients[2],
    lastRun: "En pause",
  },
  {
    id: "4",
    name: "Envoyer un SMS d'alerte météo",
    action: { label: "Météo", colorClass: "bg-white/20 text-white" },
    reaction: { label: "SMS", colorClass: "bg-white/20 text-white" },
    active: true,
    gradient: gradients[3],
    lastRun: "Il y a 42 min",
  },
  {
    id: "5",
    name: "Créer un rappel calendrier pour chaque facture",
    action: { label: "Facture", colorClass: "bg-white/20 text-white" },
    reaction: { label: "Calendrier", colorClass: "bg-white/20 text-white" },
    active: false,
    gradient: gradients[4],
    lastRun: "Brouillon",
  },
  {
    id: "6",
    name: "Poster sur Twitter quand un article Medium sort",
    action: { label: "Medium", colorClass: "bg-white/20 text-white" },
    reaction: { label: "X", colorClass: "bg-white/20 text-white" },
    active: true,
    gradient: gradients[0],
    lastRun: "Il y a 2 h",
  },
];
