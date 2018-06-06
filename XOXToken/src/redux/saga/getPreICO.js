import { call, put } from 'redux-saga/effects';
import crowdsaleInstance from '../../services/Crowdsale';
import { savePreICO } from '../actions/ico';

function* getPreICO() {
  const preICO = yield call(crowdsaleInstance.instance.preICO);
  yield put(savePreICO(preICO));
}

export default getPreICO;
