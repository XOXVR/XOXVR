import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';

import preICO from './preICO';
import ICO from './ico';
import postICO from './postICO';
import crowdsale from './crowdsaleData';
import token from './tokenInfo';
import btcWallet from './btcWallet';
import statistics from './statistics';

export default combineReducers({
  routing: routerReducer,
  preICO,
  ICO,
  postICO,
  crowdsale,
  token,
  btcWallet,
  statistics
});
