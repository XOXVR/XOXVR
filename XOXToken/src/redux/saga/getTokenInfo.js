import { call, put, all } from 'redux-saga/effects';
import tokenInstance from '../../services/Token';
import saveToken from '../actions/tokenInfo';
import crowdsaleInstance from '../../services/Crowdsale';
import web3Instance from '../../services/getWeb3';

function* getTokenInfo() {
  const tokenAddress = yield call(crowdsaleInstance.instance.token);
  const [
    name, symbol, decimals, initialSupply, myTokens
  ] = yield all([
    call(tokenInstance.instance.name), call(tokenInstance.instance.symbol),
    call(tokenInstance.instance.decimals), call(tokenInstance.instance.INITIAL_SUPPLY),
    tokenInstance.instance.balanceOf(web3Instance.eth.accounts[0], { from: web3Instance.eth.accounts[0] })
  ]);
  const token = {
    name, symbol, decimals, tokenAddress, initialSupply, myTokens
  };
  yield put(saveToken(token));
}

export default getTokenInfo;
