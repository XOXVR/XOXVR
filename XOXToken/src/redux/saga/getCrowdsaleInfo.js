import { call, put, all } from 'redux-saga/effects';
import crowdsaleInstance from '../../services/Crowdsale';
import saveCrowdsale from '../actions/crowdsaleData';

function* getCrowdsaleInfo() {
  const decimals = 1E18;
  const myWallet = crowdsaleInstance.instance.contract._eth.accounts[0];
  const addressCrowdsale = crowdsaleInstance.instance.address;
  const [
    multisig, requireOnce, salesTokens, rate, softcap, owner, deputyAddress
  ] = yield all([
    call(crowdsaleInstance.instance.multisig),
    call(crowdsaleInstance.instance.requireOnce), call(crowdsaleInstance.instance.salesTokens), call(crowdsaleInstance.instance.rate),
    call(crowdsaleInstance.instance.SOFT_CAP), call(crowdsaleInstance.instance.owner),
    call(crowdsaleInstance.instance.proxyMultisig)
  ]);
  const crowdsaleData = {
    multisig, decimals, requireOnce, salesTokens, rate, softcap, owner, addressCrowdsale, myWallet, deputyAddress
  };
  yield put(saveCrowdsale(crowdsaleData));
}

export default getCrowdsaleInfo;
