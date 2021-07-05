import PropTypes from 'prop-types';
import React from 'react';
import { connect } from 'react-redux';
import { spacing, Spacing } from 'react-elemental';
import Toast from 'client/app/react/components/ui/toast';

/**
 * Container for displaying active toasts in the global store.
 */
const ToastContainer = ({ toasts }) => (
  <div
    style={{
      bottom: spacing.huge,
      left: 0,
      margin: 'auto',
      position: 'fixed',
      right: 0,
      maxWidth: '450px',
      zIndex: 100,
    }}
  >
    <Spacing size="huge" left right>
      {toasts.map((toast, idx) => (
        <Spacing key={toast.toastID} bottom={idx < toasts.length - 1}>
          <Toast duration={toast.duration}>
            {toast.text}
          </Toast>
        </Spacing>
      ))}
    </Spacing>
  </div>
);

ToastContainer.propTypes = {
  // HOC props
  toasts: PropTypes.arrayOf(PropTypes.shape({
    toastID: PropTypes.number.isRequired,
    text: PropTypes.string.isRequired,
  })).isRequired,
};

const mapStateToProps = ({ toasts }) => ({
  toasts: toasts.active,
});

export default connect(mapStateToProps)(ToastContainer);
