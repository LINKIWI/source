import React from 'react';
import { background } from 'client/app/util/style/color';

/**
 * Simple horizontal separator bar.
 */
const Separator = () => (
  <div
    style={{
      backgroundColor: background.dark.GAMMA,
      height: '1.5px',
    }}
  />
);

export default Separator;
