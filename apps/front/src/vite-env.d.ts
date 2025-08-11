/// <reference types="vite/client" />

interface ViteTypeOptions {
  strictImportMetaEnv: unknown;
}

interface ImportMetaEnv {
  readonly VITE_API_URL: string;
  // more env variables...
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
