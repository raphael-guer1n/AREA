export type FieldDefinition = {
  name: string;
  type: "text" | "number" | "date";
  label: string;
  required?: boolean;
  default?: string | number;
  private?: boolean;
};

export type ActionDefinition = {
  id: string;
  title: string;
  label: string;
  type: string;
  fields: FieldDefinition[];
  output_fields?: FieldDefinition[];
};

export type ReactionDefinition = {
  id: string;
  title: string;
  label: string;
  url?: string;
  method?: string;
  fields: FieldDefinition[];
};

export type MockService = {
  id: string;
  name: string;
  url: string;
  badge: string;
  category?: string;
  gradient: { from: string; to: string };
  actions: ActionDefinition[];
  reactions: ReactionDefinition[];
  connected: boolean;
};

const gradients: Array<{ from: string; to: string }> = [
  { from: "#002642", to: "#0b3c5d" },
  { from: "#840032", to: "#a33a60" },
  { from: "#e59500", to: "#f2b344" },
  { from: "#5B834D", to: "#68915a" },
  { from: "#02040f", to: "#1b2640" },
];

export { gradients };

export const mockServices: MockService[] = [
  {
    id: "timer",
    name: "Timer",
    url: "#",
    badge: "Tm",
    category: "Interne",
    gradient: gradients[0],
    actions: [
      {
        id: "cron_action",
        title: "cron action",
        label: "Timer",
        type: "cron",
        fields: [
          {
            name: "delay",
            type: "number",
            label: "Delay (seconds)",
            required: true,
            default: 0,
          },
        ],
        output_fields: [
          {
            name: "delay",
            type: "number",
            label: "Delay (seconds)",
            private: false,
          },
        ],
      },
    ],
    reactions: [],
    connected: true,
  },
  {
    id: "google",
    name: "Google",
    url: "https://www.google.com",
    badge: "G",
    category: "Productivité",
    gradient: gradients[1],
    actions: [],
    reactions: [
      {
        id: "create_event",
        title: "create_event",
        label: "Create Event",
        url: "https://www.googleapis.com/calendar/v3/calendars/{google_calendar}/events",
        method: "POST",
        fields: [
          {
            name: "summary",
            type: "text",
            label: "Event title",
            required: true,
            default: "",
          },
          {
            name: "description",
            type: "text",
            label: "Event description",
            required: false,
            default: "",
          },
          {
            name: "start_time",
            type: "date",
            label: "Start Time",
            required: true,
            default: "",
          },
          {
            name: "end_time",
            type: "date",
            label: "End Time",
            required: true,
            default: "",
          },
          {
            name: "calendar",
            type: "text",
            label: "Calendar",
            required: true,
            default: "primary",
          },
        ],
      },
    ],
    connected: true,
  },
  {
    id: "github",
    name: "GitHub",
    url: "https://github.com",
    badge: "Gh",
    category: "Développeurs",
    gradient: gradients[2],
    actions: [],
    reactions: [],
    connected: true,
  },
];
