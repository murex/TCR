export enum TcrTimerState {
  OFF = 'off',
  PENDING = 'pending',
  RUNNING = 'running',
  STOPPED = 'stopped',
  TIMEOUT = 'timeout',
}

export interface TcrTimer {
  state: string;
  timeout: string;
  elapsed: string;
  remaining: string;
}
