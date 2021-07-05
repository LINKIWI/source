import Color from 'color';
import { colors } from 'react-elemental';

export const background = {
  primary: {
    ALPHA: colors.primary,
    BETA: new Color(colors.primary).fade(0.4).string(),
    DELTA: new Color(colors.primary).fade(0.7).string(),
    GAMMA: new Color(colors.primary).fade(0.8).string(),
    EPSILON: new Color(colors.primary).fade(0.9).string(),
  },
  dark: {
    ALPHA: 'rgba(0, 0, 0, 0.9)',
    BETA: 'rgba(0, 0, 0, 0.8)',
    DELTA: 'rgba(0, 0, 0, 0.1)',
    GAMMA: 'rgba(0, 0, 0, 0.025)',
    EPSILON: 'rgba(0, 0, 0, 0.015)',
    IOTA: 'rgba(0, 0, 0, 0.01)',
  },
};

export const text = {
  highlight: {
    // Portion of a string that has a search match.
    MATCH: 'rgba(250, 210, 0, 0.15)',
    // Indication that the line serves as context.
    CONTEXT: new Color(colors.primary).fade(0.9).string(),
    // Selected portion of a line.
    SELECT: new Color(colors.primary).fade(0.7).string(),
  },
  light: {
    ALPHA: 'rgba(255, 255, 255, 0.95)',
    BETA: 'rgba(255, 255, 255, 0.9)',
    GAMMA: 'rgba(255, 255, 255, 0.75)',
    EPSILON: 'rgba(255, 255, 255, 0.5)',
  },
  dark: {
    ALPHA: 'rgba(0, 0, 0, 0.8)',
    BETA: 'rgba(0, 0, 0, 0.6)',
    GAMMA: 'rgba(0, 0, 0, 0.3)',
    EPSILON: 'rgba(0, 0, 0, 0.2)',
  },
};
