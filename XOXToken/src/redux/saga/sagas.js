import { call, takeLatest } from 'redux-saga/effects';
import { LOCATION_CHANGE } from 'react-router-redux';

import web3Instance from '../../services/getWeb3';
import stripe from '../../services/Stripe';
import crowdsaleInstance from '../../services/Crowdsale';
import tokenInstance from '../../services/Token';
import getPreICO from './getPreICO';
import getICO from './getICO';
import getPostICO from './getPostICO';
import getCrowdsaleInfo from './getCrowdsaleInfo';
import getTokenInfo from './getTokenInfo';
import getBTCAddress from './getBTCWallet';
import getStatistics from './getStatistics';


const regularExpression = /^([0-9]+[.])?[0-9]+$/;

function* loading() {
  try {
    yield crowdsaleInstance.getInstance();
    yield getPreICO();
    yield getICO();
    yield getPostICO();
    yield getCrowdsaleInfo();
    yield tokenInstance.getInstance();
    yield getTokenInfo();
  } catch (err) {
    console.log('Err, when connect to Crowdsale', err);
  }
}

function* saleTokens({ payload }) {
  if (payload === 0 || !regularExpression.test(payload)) {
    alert('This is not a valid number!');
  } else {
    yield call(
      crowdsaleInstance.instance.saleTokens,
      { from: web3Instance.eth.accounts[0], value: +payload * 1E18, gas: 200000 }
    );
    yield loading();
  }
}

function* setMultisig({ payload }) {
  const isAddress = web3Instance.isAddress(payload);
  if (isAddress) {
    yield call(
      crowdsaleInstance.instance.setMultisig,
      payload, { from: web3Instance.eth.accounts[0], gas: 200000 }
    );
    yield crowdsaleInstance.getInstance();
    yield getCrowdsaleInfo();
  } else {
    alert('No valid address!');
  }
}

function* setDeputyMultisig({ payload }) {
  const isAddress = web3Instance.isAddress(payload);
  if (isAddress) {
    yield call(
      crowdsaleInstance.instance.setProxyMultisig,
      payload, { from: web3Instance.eth.accounts[0], gas: 200000 }
    );
    yield crowdsaleInstance.getInstance();
    yield getCrowdsaleInfo();
  } else {
    alert('No valid address!');
  }
}

function* refund() {
  yield crowdsaleInstance.instance.refund({ from: web3Instance.eth.accounts[0], gas: 200000 });
  yield loading();
}

function* burn() {
  yield crowdsaleInstance.instance.burnTokens({ from: web3Instance.eth.accounts[0], gas: 200000 });
  yield crowdsaleInstance.getInstance();
  yield getPreICO();
  yield getICO();
}

function* setDate({ payload }) {
  const [startPreICO, startICO, startPostICO] = payload;
  const res = payload.filter(el => el);
  if (res.length === 3) {
    yield call(
      crowdsaleInstance.instance.setDate,
      startPreICO, startICO, startPostICO,
      { from: web3Instance.eth.accounts[0], gas: 200000 }
    );
    yield crowdsaleInstance.getInstance();
    yield getPreICO();
    yield getICO();
    yield getPostICO();
  } else {
    alert('Not all fields are filled!');
  }
}

function* buyTokens({ payload }) {
  if (payload.walletEthFiat.length === 0 || payload.cents === 0 || !regularExpression.test(payload.cents)) {
    alert('Not all fields are filled!');
  } else {
    const handler = yield stripe(payload);
    handler.open({ description: 'Buy tokens', amount: +payload.cents });
    yield crowdsaleInstance.getInstance();
    yield getPostICO();
    yield getTokenInfo();
    yield getCrowdsaleInfo();
  }
}

function* changeOwner({ payload }) {
  const isAddress = web3Instance.isAddress(payload);
  if (isAddress) {
    yield call(
      crowdsaleInstance.instance.transferOwnership,
      payload, { from: web3Instance.eth.accounts[0], gas: 200000 }
    );
    yield crowdsaleInstance.getInstance();
    yield getCrowdsaleInfo();
  } else {
    alert('No valid address!');
  }
}

function* sagas() {
  yield takeLatest(LOCATION_CHANGE, loading);
  yield takeLatest('SALE_TOKENS', saleTokens);
  yield takeLatest('SET_MULTISIG', setMultisig);
  yield takeLatest('SET_DEPUTY_MULTISIG', setDeputyMultisig);
  yield takeLatest('REFUND', refund);
  yield takeLatest('BURN', burn);
  yield takeLatest('BUY_TOKENS', buyTokens);
  yield takeLatest('SET_DATE', setDate);
  yield takeLatest('NEW_WALLET', getBTCAddress);
  yield takeLatest('NEW_OWNER', changeOwner);
  yield takeLatest('STATISTICS', getStatistics);
}

export default sagas;
