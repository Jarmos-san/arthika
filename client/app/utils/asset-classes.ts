/** @description Folio code shown when the asset class name is blank. */
const FALLBACK_FOLIO = "--";

/** @description Minimum number of characters for a word to contribute to the folio code. */
const MIN_WORD_LENGTH = 1;

/**
 * @description Derives a ledger-style folio code from an asset class name. Single-word names
 * use their first two letters ("Crypto" → "CR"); multi-word names use the
 * initials of the first two words ("Fixed income" → "FI"). Falls back to "--"
 * for blank names.
 *
 * @param {string} name - The asset class name.
 *
 * @returns {string} The uppercased folio code.
 */
const getFolioCode = (name: string): string => {
  const words = name
    .trim()
    .split(/\s+/u)
    .filter((word) => word.length >= MIN_WORD_LENGTH);

  const [first, second] = words;

  if (first === undefined) {
    return FALLBACK_FOLIO;
  }

  if (second === undefined) {
    const [initial, next] = first;
    return `${initial ?? ""}${next ?? ""}`.toUpperCase();
  }

  const [firstInitial] = first;
  const [secondInitial] = second;

  return `${firstInitial ?? ""}${secondInitial ?? ""}`.toUpperCase();
};

export default getFolioCode;
