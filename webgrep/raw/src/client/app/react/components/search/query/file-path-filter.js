import PropTypes from 'prop-types';
import React from 'react';
import { Button, Spacing, Text } from 'react-elemental';
import { withForwardedRef } from '@linkiwi/hoc';
import TextField from 'client/app/react/components/ui/text-field';
import Toggle from 'client/app/react/components/ui/toggle';
import { NODE_ID } from 'client/app/util/constants/dom';
import { withDefaultPrevented } from 'client/app/util/navigation';
import { text } from 'client/app/util/style/color';

/**
 * Filter results by a file path pattern.
 */
const FilePathFilter = ({ forwardedRef, value, suggestions, onChange, onHide }) => (
  <div style={{ display: 'flex', flexDirection: 'column', maxHeight: '500px' }}>
    <form onSubmit={withDefaultPrevented(onHide)}>
      <TextField
        ref={forwardedRef}
        id={NODE_ID.SEARCH_FILE_PATH_FIELD}
        placeholder="Restrict by file path…"
        autoComplete="off"
        value={value}
        onChange={onChange}
      />

      {suggestions.length > 0 && (
        <Spacing size="20px" top>
          <Spacing size="10px" bottom>
            <Text color={text.dark.BETA} size="lambda" bold uppercase>
              Suggestions
            </Text>
          </Spacing>

          <Spacing
            size="-8px"  // Compensate for bottom margin of individual toggles
            style={{
              alignItems: 'center',
              display: 'flex',
              flexWrap: 'wrap',
            }}
            bottom
          >
            {suggestions.map((suggestion, idx) => (
              <Spacing key={suggestion} size="11px" right={idx < suggestions.length - 1}>
                <Spacing size="8px" bottom>
                  <Toggle
                    onClick={() => {
                      // Simulate populating the file path search filter (and thus executing a new
                      // search) with the selected file extension literal.
                      onChange({ target: { value: suggestion } });
                      onHide();
                    }}
                    secondary
                  >
                    {suggestion}
                  </Toggle>
                </Spacing>
              </Spacing>
            ))}
          </Spacing>
        </Spacing>
      )}

      <Spacing style={{ alignItems: 'center', display: 'flex', justifyContent: 'flex-end' }} top>
        <Spacing size="small" right>
          <Button
            text="Reset"
            onClick={() => onChange({ target: { value: '' } })}
            style={{ border: 0 }}
            secondary
          />
        </Spacing>

        <Button
          type="submit"
          text="OK"
          onClick={onHide}
          style={{ borderRadius: '3px' }}
        />
      </Spacing>
    </form>
  </div>
);

FilePathFilter.propTypes = {
  value: PropTypes.string.isRequired,
  suggestions: PropTypes.arrayOf(PropTypes.string).isRequired,
  onChange: PropTypes.func.isRequired,
  onHide: PropTypes.func.isRequired,
  // HOC props
  forwardedRef: PropTypes.oneOfType([
    PropTypes.shape({ current: PropTypes.instanceOf(Element) }),
    PropTypes.func,
  ]),
};

FilePathFilter.defaultProps = {
  forwardedRef: null,
};

export default withForwardedRef(FilePathFilter);
