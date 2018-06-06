import { call, put } from 'redux-saga/effects';
import crowdsaleInstance from '../../services/Crowdsale';
import { saveICO } from '../actions/ico';

function* getICO() {
  const ICO = yield call(crowdsaleInstance.instance.ICO);
  yield put(saveICO(ICO));
}

export default getICO;
