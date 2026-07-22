/**
 * @description Labeled text input with v-model support. All other attributes (e.g. `id`,
 * `type`, `placeholder`, `class`) fall through to the native `<input>`
 * element.
 */
export default interface Props {
  /** @description Visible text displayed above the input. */
  label: string;
}
