import PropTypes from 'prop-types';
import React from 'react';
import { colors, Button, Spacing } from 'react-elemental';
import { MdSearch } from 'react-icons/md';
import { compose, withForm, withForwardedRef } from '@linkiwi/hoc';
import RepositorySelector from 'client/app/react/components/search/query/repository-selector';
import TextField from 'client/app/react/components/ui/text-field';
import { NODE_ID } from 'client/app/util/constants/dom';

/**
 * Modal providing repository filter options.
 */
const RepositoryFilter = ({
  forwardedRef,
  form,
  repositories,
  selectedRepos,
  onHide,
  onRepositoryToggle,
  handleFormChange,
}) => (
  <div style={{ display: 'flex', flexDirection: 'column', maxHeight: '500px' }}>
    <Spacing bottom>
      <TextField
        ref={forwardedRef}
        id={NODE_ID.SEARCH_REPOSITORY_FILTER_FIELD}
        before={<MdSearch style={{ color: colors.primary, fontSize: '20px' }} />}
        placeholder="Filter repositories…"
        autoComplete="off"
        value={form.filterRepoQuery}
        onChange={handleFormChange('filterRepoQuery')}
      />
    </Spacing>

    <Spacing style={{ flexGrow: 1, overflow: 'auto' }} bottom>
      {repositories
        .filter((repo) =>
          !form.filterRepoQuery.length ||
          repo.name.includes(form.filterRepoQuery) ||
          repo.remote.includes(form.filterRepoQuery))
        .map((repo) => (
          <RepositorySelector
            key={repo.name}
            name={repo.name}
            description={repo.remote}
            isSelected={selectedRepos.has(repo.name)}
            onClick={() => onRepositoryToggle(repo.name)}
          />
        ))}
    </Spacing>

    <div style={{ alignItems: 'center', display: 'flex', justifyContent: 'flex-end' }}>
      <Spacing size="small" right>
        <Button
          text="Reset"
          onClick={() => {
            handleFormChange('filterRepoQuery')({ target: { value: '' } });
            onRepositoryToggle();
          }}
          style={{ border: 0 }}
          secondary
        />
      </Spacing>

      <Button
        text="OK"
        onClick={onHide}
        style={{ borderRadius: '3px' }}
      />
    </div>
  </div>
);

RepositoryFilter.propTypes = {
  repositories: PropTypes.arrayOf(PropTypes.shape({
    name: PropTypes.string.isRequired,
    remote: PropTypes.string.isRequired,
  }).isRequired).isRequired,
  selectedRepos: PropTypes.instanceOf(Set).isRequired,
  onHide: PropTypes.func.isRequired,
  onRepositoryToggle: PropTypes.func.isRequired,
  // HOC props
  forwardedRef: PropTypes.oneOfType([
    PropTypes.shape({ current: PropTypes.instanceOf(Element) }),
    PropTypes.func,
  ]),
  form: PropTypes.shape({
    filterRepoQuery: PropTypes.string.isRequired,
  }).isRequired,
  handleFormChange: PropTypes.func.isRequired,
};

RepositoryFilter.defaultProps = {
  forwardedRef: null,
};

export default compose(
  withForwardedRef,
  withForm({
    initial: () => ({ filterRepoQuery: '' }),
  }),
)(RepositoryFilter);
