export interface TcrSessionInfo {
  baseDir: string;
  workDir: string;
  language: string;
  toolchain: string;
  vcsName: string;
  vcsSession: string;
  commitOnFail: boolean;
  gitAutoPush: boolean;
  messageSuffix: string;
}
