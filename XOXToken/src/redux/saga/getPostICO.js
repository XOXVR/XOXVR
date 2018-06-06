import { call, put } from 'redux-saga/effects';
import crowdsaleInstance from '../../services/Crowdsale';
import { savePostICO } from '../actions/ico';

function* getPostICO() {
  const postICO = yield call(crowdsaleInstance.instance.postICO);
  yield put(savePostICO(postICO));
}

export default getPostICO;
