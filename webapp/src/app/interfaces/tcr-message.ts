export enum TcrMessageType {
  SIMPLE = "simple",
  INFO = "info",
  TITLE = "title",
  SUCCESS = "success",
  WARNING = "warning",
  ERROR = "error",
  ROLE = "role",
  TIMER = "timer",
}

export interface TcrMessage {
  type: TcrMessageType;
  severity: string;
  text: string;
  emphasis: boolean;
  timestamp: string;
}
