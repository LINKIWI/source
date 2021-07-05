import 'raf/polyfill';
import Enzyme from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

// Injected by webpack during build; stub out for testing
window.__SSR_INJECTED_GLOBALS__ = {};

Enzyme.configure({ adapter: new Adapter() });
