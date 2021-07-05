/**
 * Logic module superclass containing shared utilities common to all logic modules.
 */
export default class BaseLogic {
  constructor(ctx) {
    this.ctx = ctx;
  }

  /**
   * Convert a gRPC error to a Supercharged error.
   *
   * @param {Object} grpcErr gRPC error object.
   * @return {Object} Supercharged API-compatible error object.
   */
  formatErr(grpcErr) {  // eslint-disable-line class-methods-use-this
    return {
      // Error codes are offset by 100 to differentiate between Supercharged errors and gRPC errors.
      code: grpcErr.code + 100,
      message: `${grpcErr.details} (${grpcErr.code})`,
    };
  }
}
