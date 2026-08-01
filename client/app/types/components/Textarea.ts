/**
 * @description Labeled textarea with v-model support. All other attributes (e.g. `id`,
 * `placeholder`, `rows`, `class`) fall through to the native `<textarea>`
 * element.
 */
export default interface Props {
  /** @description Visible text displayed above the textarea. */
  label: string;
}
