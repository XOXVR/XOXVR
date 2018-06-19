import React, { Component } from 'react';
import { Row, Col, Input, Button } from 'reactstrap';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import moment from 'moment';

import { setMultisig, setDeputyMultisig } from '../redux/actions/setMultisig';
import refund from '../redux/actions/refund';
import setDate from '../redux/actions/setDate';
import burn from '../redux/actions/burn';
import { getStatistics } from '../redux/actions/getStatistics';
import setNewOwner from '../redux/actions/changeOwner';
import '../styles/Crowdsale.scss';

const isValidDate = date => date.every(elem => moment.unix(elem).isValid());

class AdminPanel extends Component {
  state = {
    newMultisig: '',
    deputyMultisig: '',
    startPreICO: 0,
    startICO: 0,
    startPostICO: 0,
    owner: '',
    countPagination: 0
  }

  render() {
    const nowDate = moment().unix();
    const {
      crowdsale, setMultisig, setDeputyMultisig, refund, burn, setDate, preICO, postICO, setNewOwner, getStatistics, statistics
    } = this.props;
    const { startPreICO, startICO, startPostICO } = this.state;
    if (crowdsale.owner === crowdsale.myWallet || crowdsale.deputyAddress === crowdsale.myWallet) {
      return (
        <Row style={{ marginTop: 15 }}><h3 style={{ textAlign: 'center' }}>Admin Panel</h3>
          <Row style={{ marginTop: 15 }}><h5>Setting start date</h5></Row>
          <Row className="funcRow" style={{ marginBottom: 15 }}>
            <Col md={{ size: 6 }}>
              <Input
                onChange={e => this.setState({ startPreICO: e.target.value })}
                placeholder="startPreICO UNIX format"
                style={{ width: 150, marginRight: 15, }}
                disabled={nowDate > preICO.start}
              />
              <Input
                onChange={e => this.setState({ startICO: e.target.value })}
                placeholder="startICO UNIX format"
                style={{ width: 150, marginRight: 15 }}
                disabled={nowDate > preICO.start}
              />
              <Input
                onChange={e => this.setState({ startPostICO: e.target.value })}
                placeholder="startPostICO UNIX format"
                style={{ width: 150, marginRight: 10 }}
                disabled={nowDate > preICO.start}
              />
              <Button
                className="funcButton"
                color="info"
                disabled={nowDate > preICO.start}
                onClick={() => {
                  const data = [startPreICO, startICO, startPostICO];
                  if (isValidDate(data)) {
                    setDate(data);
                  } else {
                    alert('One of the dates is incorrect. Please check the input data!');
                  }
                }}
              >Set date
              </Button>
            </Col>
          </Row>

          <Row style={{ marginTop: 15 }}><h5>Owner</h5></Row>
          <Row className="funcRow" style={{ marginBottom: 15 }}>
            <Col md={{ size: 6 }}>
              <Input
                onChange={e => this.setState({ owner: e.target.value })}
                placeholder="Address of new wallet"
                style={{ width: 400, marginRight: 10 }}
              />
              <Button className="funcButton" color="info" onClick={() => setNewOwner(this.state.owner)}>Set new owner</Button>
            </Col>
          </Row>

          <Row style={{ marginTop: 15 }}><h5>Change Wallet</h5></Row>
          <Row className="funcRow" style={{ marginBottom: 15 }}>
            <Col md={{ size: 6 }}>
              <Input
                onChange={e => this.setState({ newMultisig: e.target.value })}
                placeholder="Address of new wallet"
                style={{ width: 400, marginRight: 10 }}
              />
              <Button className="funcButton" color="info" onClick={() => setMultisig(this.state.newMultisig)}>Set multisig</Button>
            </Col>
          </Row>

          <Row style={{ marginTop: 15 }}><h5>Set Deputy Multisig</h5></Row>
          <Row className="funcRow" style={{ marginBottom: 15 }}>
            <Col md={{ size: 6 }}>
              <Input
                onChange={e => this.setState({ deputyMultisig: e.target.value })}
                placeholder="Address of wallet"
                style={{ width: 400, marginRight: 10 }}
              />
              <Button className="funcButton" color="info" onClick={() => setDeputyMultisig(this.state.deputyMultisig)}>Set trusted wallet</Button>
            </Col>
          </Row>
          <Col md={{ size: 6 }}>
            <Button style={{ marginRight: 50 }} className="funcButton" disabled={nowDate < postICO.start} color="info" onClick={() => refund()}>Refund</Button>
            <Button className="funcButton" disabled={nowDate < postICO.start} color="info" onClick={() => burn()}>Burn tokens</Button>
          </Col>

          <Row style={{ marginTop: 15 }}><h5>Statistics</h5></Row>
          <Row className="funcRow" style={{ marginBottom: 15 }}>
            <Col md={{ size: 6 }}>
              <Button
                className="funcButton"
                color="info"
                onClick={() => {
                this.setState({ countPagination: this.state.countPagination + 50 });
                getStatistics(this.state.countPagination);
                }}>Get statistics +{this.state.countPagination} people
              </Button>
              <div className="statistics" style={{ display: 'flex', flexWrap: 'nowrap', width: '100%' }}>
                <div className="invest" style={{ width: '50%' }}>Investors</div>
                <div className="valInvest" style={{ width: '50%' }}>{Object.keys(statistics.objInvest).length || 'No investors yet'}</div>
              </div>
              {Object.entries(statistics.objInvest).map(([key, val]) => {
                return (
                  <div className="statistics" key={key} style={{ display: 'flex', flexWrap: 'nowrap', width: '100%' }}>
                    <div className="invest" style={{ width: '50%' }}>{key}</div>
                    <div className="valInvest" style={{ width: '50%' }}>{val / 1E18}</div>
                  </div>
                );
              })}
            </Col>
          </Row>

        </Row>
      );
    }
    return (<div />);
  }
}

export default connect(
  state => ({
    crowdsale: state.crowdsale,
    preICO: state.preICO,
    postICO: state.postICO,
    statistics: state.statistics
  }),
  dispatch => bindActionCreators({
    setMultisig,
    setDeputyMultisig,
    refund,
    burn,
    setDate,
    setNewOwner,
    getStatistics
  }, dispatch)
)(AdminPanel);
