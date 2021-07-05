import React from 'react';
import { Spinner as ElementalSpinner } from 'react-elemental';
import Delayed from 'client/app/react/components/ui/delayed';
import { background } from 'client/app/util/style/color';

/**
 * Styled loading spinner, with built-in mount delay.
 */
const Spinner = () => (
  <div style={{ height: '30px' }}>
    <Delayed>
      <ElementalSpinner
        ringColor={background.dark.GAMMA}
        style={{ display: 'block', margin: 'auto' }}
      />
    </Delayed>
  </div>
);

export default Spinner;
