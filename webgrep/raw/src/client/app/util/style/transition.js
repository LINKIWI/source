import { durations, timing } from 'react-elemental';

export const transition = {
  all: {
    DEFAULT: `all ${durations.beta} ${timing.default}`,
    FAST: `all ${durations.alpha} ${timing.default}`,
  },
  margin: {
    DEFAULT: `margin ${durations.beta} ${timing.default}`,
  },
  opacity: {
    DEFAULT: `opacity ${durations.beta} ${timing.default}`,
  },
};

export default undefined;
