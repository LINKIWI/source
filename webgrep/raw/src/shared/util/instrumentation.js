/**
 * Utility for recording elapsed durations. The stopwatch is started when the function is invoked.
 *
 * @return {Object} Object with a single method, `elapsed`.
 */
export const stopwatch = () => {
  const start = Date.now();

  return {
    // Lazily report the number of milliseconds elapsed since instantiation.
    elapsed: () => Date.now() - start,
  };
};

export default undefined;
