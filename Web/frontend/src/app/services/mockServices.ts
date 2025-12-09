export type MockService = {
  id: string;
  name: string;
  url: string;
  badge: string;
  category?: string;
  gradient: { from: string; to: string };
  actions: string[];
  reactions: string[];
  connected: boolean;
};

const gradients: Array<{ from: string; to: string }> = [
  { from: "#002642", to: "#0b3c5d" },
  { from: "#840032", to: "#a33a60" },
  { from: "#e59500", to: "#f2b344" },
  { from: "#5B834D", to: "#68915a" },
  { from: "#02040f", to: "#1b2640" },
];

export const mockServices: MockService[] = [
  {
    id: "google",
    name: "Google",
    url: "https://www.google.com",
    badge: "G",
    category: "Recherche",
    gradient: gradients[0],
    actions: ["Nouvelle recherche tendance", "Formulaire Google soumis", "Nouveau fichier Drive"],
    reactions: ["Créer un événement Agenda", "Envoyer un email Gmail", "Ajouter une tâche Keep"],
    connected: true,
  },
  {
    id: "slack",
    name: "Slack",
    url: "https://slack.com",
    badge: "Sl",
    category: "Communication",
    gradient: gradients[1],
    actions: ["Message posté dans un canal", "Réaction ajoutée", "Nouveau membre dans l'espace"],
    reactions: ["Poster un message", "Épingler un message", "Créer un rappel"],
    connected: true,
  },
  {
    id: "github",
    name: "GitHub",
    url: "https://github.com",
    badge: "Gh",
    category: "Développeurs",
    gradient: gradients[2],
    actions: ["Nouveau push sur la branche", "Ouverture d'une issue", "Pull request créée"],
    reactions: ["Créer une issue", "Commenter une PR", "Ouvrir une discussion"],
    connected: true,
  },
  {
    id: "discord",
    name: "Discord",
    url: "https://discord.com",
    badge: "Di",
    category: "Communication",
    gradient: gradients[3],
    actions: ["Message posté dans un salon", "Nouvel utilisateur rejoint", "Réaction ajoutée"],
    reactions: ["Envoyer un message bot", "Attribuer un rôle", "Envoyer un DM"],
    connected: false,
  },
  {
    id: "notion",
    name: "Notion",
    url: "https://www.notion.so",
    badge: "No",
    category: "Organisation",
    gradient: gradients[4],
    actions: ["Nouvelle page ajoutée", "Propriété mise à jour", "Commentaire ajouté"],
    reactions: ["Créer une page", "Mettre à jour un statut", "Ajouter un rappel daté"],
    connected: true,
  },
  {
    id: "dropbox",
    name: "Dropbox",
    url: "https://www.dropbox.com",
    badge: "Db",
    category: "Stockage",
    gradient: gradients[1],
    actions: ["Nouveau fichier déposé", "Fichier mis à jour", "Dossier partagé"],
    reactions: ["Créer un lien partagé", "Déplacer un fichier", "Ajouter un commentaire"],
    connected: false,
  },
  {
    id: "jira",
    name: "Jira",
    url: "https://www.atlassian.com/software/jira",
    badge: "Ji",
    category: "Productivité",
    gradient: gradients[0],
    actions: ["Ticket créé", "Statut changé", "Commentaire ajouté"],
    reactions: ["Créer un ticket", "Ajouter un commentaire", "Mettre à jour le statut"],
    connected: true,
  },
  {
    id: "trello",
    name: "Trello",
    url: "https://trello.com",
    badge: "Tr",
    category: "Organisation",
    gradient: gradients[2],
    actions: ["Carte créée", "Échéance proche", "Membre ajouté à une carte"],
    reactions: ["Créer une carte", "Déplacer une carte", "Ajouter une checklist"],
    connected: false,
  },
];
