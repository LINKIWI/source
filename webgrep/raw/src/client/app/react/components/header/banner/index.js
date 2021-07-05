import PropTypes from 'prop-types';
import React from 'react';
import { Spacing, Text } from 'react-elemental';
import BannerCloseButton from 'client/app/react/components/header/banner/close';
import BannerLink from 'client/app/react/components/header/banner/link';
import { text } from 'client/app/util/style/color';

const Banner = ({ title, description, link, href, onHide }) => (
  <Spacing
    size="20px"
    style={{ alignItems: 'center', display: 'flex', justifyContent: 'space-between' }}
    top
    bottom
    padding
  >
    <div>
      <Spacing size="small" right inline>
        <Text color={text.light.ALPHA} bold inline>
          {title}
        </Text>
      </Spacing>

      <Spacing size="small" right={!!link && !!href} inline>
        <Text color={text.light.ALPHA} inline>
          {description}
        </Text>
      </Spacing>

      {link && href && (
        <BannerLink title={link} href={href} />
      )}
    </div>

    {onHide && (
      <Spacing left>
        <BannerCloseButton onClick={onHide} />
      </Spacing>
    )}
  </Spacing>
);

Banner.propTypes = {
  title: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
  link: PropTypes.string,
  href: PropTypes.string,
  onHide: PropTypes.func,
};

Banner.defaultProps = {
  link: null,
  href: null,
  onHide: null,
};

export default Banner;
