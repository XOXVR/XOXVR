import { call, put } from 'redux-saga/effects';
import crowdsaleInstance from '../../services/Crowdsale';
import { saveStatistics } from '../actions/getStatistics';
import web3Instance from '../../services/getWeb3';

function* getStatistics({ payload }) {
  const numIn = (yield crowdsaleInstance.instance.returnNumberInvestors({ from: web3Instance.eth.accounts[0], gas: 200000 })).toNumber();
  const limit = (payload + 50) > numIn ? numIn : payload + 50;
  const objInvest = {};
  for (let i = payload; i < limit; i += 1) {
    const addressInvestors = yield call(
      crowdsaleInstance.instance.investors,
      i, { from: web3Instance.eth.accounts[0], gas: 200000 }
    );
    const valInvestments = yield call(
      crowdsaleInstance.instance.investments,
      addressInvestors, { from: web3Instance.eth.accounts[0], gas: 200000 }
    );
    objInvest[addressInvestors] = valInvestments.toNumber();
  }
  yield put(saveStatistics(objInvest));
}


export default getStatistics;
