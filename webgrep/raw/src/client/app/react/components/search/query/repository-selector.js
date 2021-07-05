import PropTypes from 'prop-types';
import React from 'react';
import { colors, Spacing, Text } from 'react-elemental';
import { MdLens, MdPanoramaFishEye } from 'react-icons/md';
import { withToggleState } from '@linkiwi/hoc';
import { background, text } from 'client/app/util/style/color';
import { transition } from 'client/app/util/style/transition';

/**
 * Button allowing selection and deselection of a repository for filtering.
 */
const RepositorySelector = ({
  name,
  remote,
  isSelected,
  onClick,
  isHighlighted,
  handleHighlightStart,
  handleHighlightEnd,
}) => {
  const buttonStyle = {
    backgroundColor: isHighlighted ? background.dark.EPSILON : 'transparent',
    border: 'none',
    borderRadius: '3px',
    cursor: 'pointer',
    padding: '10px',
    textAlign: 'left',
    transition: transition.all.DEFAULT,
    width: '100%',
  };

  const iconContainerStyle = {
    display: 'inherit',
    flexShrink: 0,
    height: '14px',
    position: 'relative',
    width: '14px',
  };

  const iconStyle = {
    color: 'black',
    height: '100%',
    position: 'absolute',
    width: '100%',
    transition: transition.all.DEFAULT,
  };

  return (
    <button
      style={buttonStyle}
      onClick={onClick}
      onMouseEnter={handleHighlightStart}
      onMouseLeave={handleHighlightEnd}
      onFocus={handleHighlightStart}
      onBlur={handleHighlightEnd}
    >
      <div style={{ alignItems: 'center', display: 'flex' }}>
        <Spacing style={iconContainerStyle} right>
          <MdLens
            style={{
              ...iconStyle,
              color: colors.primary,
              opacity: isSelected ? 0.9 : 0,
            }}
          />
          <MdPanoramaFishEye
            style={{
              ...iconStyle,
              color: text.dark.EPSILON,
              opacity: isSelected ? 0 : 1,
            }}
          />
        </Spacing>

        <div>
          <Spacing size="2px" bottom>
            <Text size="kilo">
              {name}
            </Text>
          </Spacing>
          <Text color={text.dark.GAMMA} style={{ wordBreak: 'break-word' }} size="kilo">
            {remote}
          </Text>
        </div>
      </div>
    </button>
  );
};

RepositorySelector.propTypes = {
  name: PropTypes.string.isRequired,
  remote: PropTypes.string.isRequired,
  isSelected: PropTypes.bool.isRequired,
  onClick: PropTypes.func.isRequired,
  // HOC props
  isHighlighted: PropTypes.bool.isRequired,
  handleHighlightStart: PropTypes.func.isRequired,
  handleHighlightEnd: PropTypes.func.isRequired,
};

export default withToggleState({
  key: 'isHighlighted',
  enable: 'handleHighlightStart',
  disable: 'handleHighlightEnd',
})(RepositorySelector);
