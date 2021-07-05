import PropTypes from 'prop-types';
import React from 'react';
import { Image, Text } from 'react-elemental';
import { Link } from 'react-router-dom';
import { text } from 'client/app/util/style/color';

const DEFAULT_LOGO_TEXT = '// webgrep';

/**
 * Site logo in the header.
 */
const Logo = ({ src }) => (
  <Link to="/" style={{ textDecoration: 'none' }}>
    {src ? (
      <Image
        alt="webgrep"
        color="transparent"
        src={src}
        height="35px"
        imgStyle={{ filter: 'unset' }}
      />
    ) : (
      <Text color={text.light.ALPHA} style={{ whiteSpace: 'nowrap' }} bold>
        {DEFAULT_LOGO_TEXT}
      </Text>
    )}
  </Link>
);

Logo.propTypes = {
  src: PropTypes.string,
};

Logo.defaultProps = {
  src: null,
};

export default Logo;
