import PropTypes from 'prop-types';
import React from 'react';
import { colors, sizes, Link, Spacing, Text } from 'react-elemental';
import { MdNavigateNext } from 'react-icons/md';
import { withToggleState } from '@linkiwi/hoc';
import { text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Common layout for an admin control item within a panel.
 */
const AdminControl = ({
  title,
  description,
  linkTitle,
  linkHref,
  children,
  isLinkHover,
  handleLinkMouseEnter,
  handleLinkMouseLeave,
}) => (
  <div style={{ alignItems: 'center', display: 'flex', justifyContent: 'space-between' }}>
    <Spacing
      size="large"
      style={{ width: children ? '60%' : 'auto' }}
      right={!!children}
      padding
    >
      <Spacing size="micro" bottom>
        <Text size="kilo" bold>
          {title}
        </Text>
      </Spacing>

      <Text color={text.dark.BETA} style={{ wordBreak: 'break-word' }} size="kilo">
        {description}
      </Text>

      {linkTitle && linkHref && (
        <Spacing size="tiny" top>
          <Text color={colors.primary} size="kilo">
            <Link
              type="plain"
              activeColor={text.dark.ALPHA}
              href={linkHref}
              onMouseEnter={handleLinkMouseEnter}
              onMouseLeave={handleLinkMouseLeave}
            >
              {linkTitle}

              <Spacing size="micro" style={{ display: 'inherit' }} left inline>
                <MdNavigateNext
                  style={{
                    fontSize: sizes.kilo,
                    marginBottom: '-2px',
                    marginLeft: isLinkHover ? '4px' : 0,
                    transition: transition.all.DEFAULT,
                  }}
                />
              </Spacing>
            </Link>
          </Text>
        </Spacing>
      )}
    </Spacing>

    {children && (
      <div style={{ display: 'flex', justifyContent: 'flex-end', maxWidth: '280px', width: '40%' }}>
        {children}
      </div>
    )}
  </div>
);

AdminControl.propTypes = {
  title: PropTypes.string.isRequired,
  description: PropTypes.string.isRequired,
  linkTitle: PropTypes.string,
  linkHref: PropTypes.string,
  children: PropTypes.node,
  // HOC props
  isLinkHover: PropTypes.bool.isRequired,
  handleLinkMouseEnter: PropTypes.func.isRequired,
  handleLinkMouseLeave: PropTypes.func.isRequired,
};

AdminControl.defaultProps = {
  linkTitle: null,
  linkHref: null,
  children: null,
};

export default withToggleState({
  key: 'isLinkHover',
  enable: 'handleLinkMouseEnter',
  disable: 'handleLinkMouseLeave',
})(AdminControl);
