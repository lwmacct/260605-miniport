import type { SectionKey } from "./types";

const defaultSection: SectionKey = "overview";
const sectionKeys = new Set<SectionKey>(["overview", "services", "hosts", "dependencies"]);

function isSectionKey(value: string | null): value is SectionKey {
  return value !== null && sectionKeys.has(value as SectionKey);
}

function hashValue(hash: string) {
  return hash.startsWith("#") ? hash.slice(1) : hash;
}

function hashSearchParams(hash: string) {
  const value = hashValue(hash);
  const queryStart = value.indexOf("?");
  if (queryStart >= 0) {
    return new URLSearchParams(value.slice(queryStart + 1));
  }
  if (value.startsWith("?")) {
    return new URLSearchParams(value.slice(1));
  }
  return new URLSearchParams();
}

export function sectionHash(section: SectionKey) {
  return section === defaultSection ? "#/" : `#/${section}`;
}

export function readSectionFromHash(hash = window.location.hash): SectionKey {
  const value = hashValue(hash);
  const path = value.split("?")[0].replace(/^\/+/, "").replace(/\/+$/, "");
  if (isSectionKey(path)) {
    return path;
  }

  const tab = hashSearchParams(hash).get("tab");
  return isSectionKey(tab) ? tab : defaultSection;
}

export function replaceSectionHash(section: SectionKey) {
  const nextHash = sectionHash(section);
  if (window.location.hash !== nextHash) {
    window.history.replaceState(null, "", nextHash);
  }
}

export function navigateSectionHash(section: SectionKey) {
  const nextHash = sectionHash(section);
  if (window.location.hash !== nextHash) {
    window.location.hash = nextHash;
  }
}
