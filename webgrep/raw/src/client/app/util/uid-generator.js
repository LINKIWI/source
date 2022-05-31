/**
 * Simple, stateful UID generator that allocates indefinitely incrementing IDs.
 */
export default class UIDGenerator {
  id = 0;

  /**
   * Statefully retrieve the next UID.
   *
   * @return {Number} Next unique identifier, as an integer.
   */
  next() {
    return this.id++;  // eslint-disable-line no-plusplus
  }

  /**
   * Reset the UID generator.
   */
  reset() {
    this.id = 0;
  }
}
