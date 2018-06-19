import React, { Component } from 'react';
import { Row, Col, Input, Button } from 'reactstrap';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import moment from 'moment';

import buyTokens from '../redux/actions/buy';
import { generateWallet } from '../redux/actions/generateWallet';
import '../styles/Crowdsale.scss';


class BuyPanel extends Component {
  state = {
    walletEthFiat: '',
    fiatSelect: '',
    cents: 0,
    EthAddress: '',
    BtcAddress: ''
  }

  showInput = () => {
    const { postICO } = this.props;
    const { fiatSelect } = this.state;
    const nowDate = moment().unix(Date.now);
    if (fiatSelect === 'tokens') {
      return (
        <Input
          onChange={e => this.setState({ cents: e.target.value * 10 })}
          placeholder="Tokens value"
          style={{ width: 200, marginRight: 10 }}
          disabled={nowDate < postICO.start}
        />
      );
    }
    return (
      <Input
        onChange={(e) => {
          if (e.target.value.indexOf(',') !== -1) return alert('Need to put a point!');
          this.setState({ cents: e.target.value * 100 });
        }}
        placeholder="Fiat value (USD)"
        style={{ width: 100, marginRight: 10 }}
        disabled={nowDate < postICO.start}
      />
    );
  }

  render() {
    const { buyTokens, postICO, generateWallet, btcWallet } = this.props;
    const {
      walletEthFiat, cents, EthAddress, BtcAddress
    } = this.state;
    const nowDate = moment().unix(Date.now);
    return (
      <Row style={{ marginTop: 15, textAlign: 'left' }}><h3 style={{ textAlign: 'center' }}>Buy Panel</h3>
        <Row style={{ marginTop: 5, fontWeight: 'bold' }}><p>Fiat</p></Row>
        <Row className="funcRow" style={{ marginBottom: 15 }}>
          <select defaultValue="fiat" size="3" disabled={nowDate < postICO.start} onChange={e => this.setState({ fiatSelect: e.target.value })}>
            <option disabled>Select tokens or fiat</option>
            <option value="fiat">Fiat</option>
            <option value="tokens">Tokens</option>
          </select>
          <Col md={{ size: 6 }}>
            <Input
              disabled={nowDate < postICO.start}
              onChange={e => this.setState({ walletEthFiat: e.target.value })}
              placeholder="Address of wallet"
              style={{ width: 400, marginRight: 10 }}
            />
            {this.showInput()}
            <Button
              className="funcButton"
              disabled={nowDate < postICO.start}
              color="info"
              onClick={() => { const data = { walletEthFiat, cents }; buyTokens(data); }}>Buy</Button>
          </Col>
        </Row>
        <hr className="my-2" />
        <Row style={{ marginTop: 40, fontWeight: 'bold' }}><p>BTC</p></Row>
        <Row className="showBTCWallet" style={{ marginBottom: 15, textAlign: 'center' }}>{btcWallet.status}</Row>
        <Row className="funcRow" style={{ marginBottom: 15 }}>
          <Col md={{ size: 6 }}>
            <Input
              disabled={nowDate < postICO.start}
              onChange={e => this.setState({ EthAddress: e.target.value })}
              placeholder="Address of ETH wallet"
              style={{ width: 370, marginRight: 10 }}
            />
            <Input
              disabled={nowDate < postICO.start}
              onChange={e => this.setState({ BtcAddress: e.target.value })}
              placeholder="Address of BTC wallet"
              style={{ width: 370, marginRight: 20 }}
            />
            <Button className="funcButton" color="info" disabled={nowDate < postICO.start} onClick={() => { const data = { EthAddress, BtcAddress }; generateWallet(data); }}>Get generated address
            </Button>
          </Col>
        </Row>
      </Row>
    );
  }
}

export default connect(
  state => ({
    crowdsale: state.crowdsale,
    postICO: state.postICO,
    btcWallet: state.btcWallet
  }),
  dispatch => bindActionCreators({
    buyTokens,
    generateWallet
  }, dispatch)
)(BuyPanel);
