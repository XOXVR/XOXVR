import React, { Component } from 'react';
import { Row, Col, Input, Button } from 'reactstrap';
import moment from 'moment';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import '../styles/Crowdsale.scss';
import saleTokens from '../redux/actions/saleTokens';
import AdminPanel from './AdminPanel';
import BuyPanel from './BuyPanel';

const dateFormat = 'D.MM.YYYY, HH:mm:ss';

class Crowdsale extends Component {
  state = {
    eth: 0
  }

  render() {
    const { preICO ,ICO, postICO, crowdsale, saleTokens } = this.props;
    const nowDate = moment().unix(Date.now);
    return (
      <Row style={{ marginTop: 50, textAlign: 'left' }}>
        <Col>
          <Row><h3>Crowdsale</h3></Row>
          <Row>
            <Col md={{ size: 6 }}>
              <Row>Owner: {crowdsale.owner}</Row>
              <Row>Multisig: {crowdsale.multisig}</Row>
              <Row>Address: {crowdsale.addressCrowdsale}</Row>
              <Row>Trusted wallet: {crowdsale.deputyAddress}</Row>
              <Row>My wallet: {crowdsale.my_wallet}</Row>
              <Row>Tokens sold: {(crowdsale.salesTokens / crowdsale.decimals).toLocaleString()} XOX</Row>
              <Row>Token cost: {(crowdsale.rate / crowdsale.decimals).toFixed(7).toLocaleString()} ETH</Row>
              <Row>Softcap: {(crowdsale.softcap / crowdsale.decimals).toLocaleString()} XOX</Row>
            </Col>
          </Row>
          <hr className="my-2" />

          <Row>
            <Col>
              <Row><h5>PreICO</h5></Row>
              <Row>Tokens amount: {(preICO.totalSupply / crowdsale.decimals).toLocaleString()} XOX</Row>
              <Row>Tokens left: {(preICO.remainderTokens / crowdsale.decimals).toLocaleString()} XOX</Row>
              <Row>Min investments: {(preICO.minInvestment / crowdsale.decimals).toLocaleString()} ETH</Row>
              <Row>Start: { moment.unix(preICO.start).format(dateFormat) }</Row>
              <Row>End: {moment.unix(preICO.end).format(dateFormat)}</Row>
              <Row>Bonuses, if investments less ten: {preICO.bonusLessTen} %</Row>
              <Row>Bonuses, if investments more ten: {preICO.bonusMoreTen} %</Row>
            </Col>
          </Row>
          <hr className="my-2" />

          <Row>
            <Col>
              <Row><h5>ICO</h5></Row>
              <Row>Tokens amount: {(ICO.totalSupply / crowdsale.decimals).toLocaleString()} XOX</Row>
              <Row>Tokens left: {(ICO.remainderTokens / crowdsale.decimals).toLocaleString()} XOX</Row>
              <Row>Min investments: {(ICO.minInvestment / crowdsale.decimals).toLocaleString()} ETH</Row>
              <Row>Start: { moment.unix(ICO.start).format(dateFormat) }</Row>
              <Row>End: {moment.unix(ICO.end).format(dateFormat)}</Row>
              <Row>Bonuses, if investments less ten: {0} %</Row>
              <Row>Bonuses, if investments more ten: {ICO.bonusMoreTen} %</Row>
            </Col>
          </Row>
          <hr className="my-2" />

          <Row>
            <Col>
              <Row><h5>PostICO</h5></Row>
              <Row>Tokens sold: {(postICO.totalSupply / crowdsale.decimals).toLocaleString()} XOX</Row>
              <Row>Tokens left: {(postICO.remainderTokens / crowdsale.decimals).toLocaleString()} XOX</Row>
              <Row>Start: { moment.unix(postICO.start).format(dateFormat) }</Row>
            </Col>
          </Row>
          <hr className="my-2" />

          <Row style={{ marginTop: 15 }}><h5>Buy tokens</h5></Row>
          <Row className="funcRow" style={{ marginBottom: 15, flexWrap: 'nowrap' }}>
            <Col md={{ size: 6 }}>
              <Input
                onChange={e => this.setState({ eth: e.target.value })}
                placeholder="Eth value"
                disabled={nowDate < preICO.start || nowDate > ICO.end}
                style={{ marginRight: 15 }}
              />
              <Button
                className="funcButton"
                disabled={nowDate < preICO.start || nowDate > ICO.end}
                color="danger"
                onClick={() => saleTokens(this.state.eth)}
              >Sale Tokens
              </Button>
            </Col>
          </Row>
          <hr className="my-2" />
          <BuyPanel />
          <hr className="my-2" />
          <AdminPanel />
        </Col>
      </Row>
    );
  }
}

export default connect(
  state => ({
    preICO: state.preICO,
    ICO: state.ICO,
    postICO: state.postICO,
    crowdsale: state.crowdsale
  }),
  dispatch => bindActionCreators({
    saleTokens
  }, dispatch)
)(Crowdsale);
