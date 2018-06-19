import React, { Component } from 'react';
import { Row, Col } from 'reactstrap';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

class Token extends Component {
  render() {
    const { token } = this.props;
    return (
      <Row style={{ marginTop: 50, textAlign: 'left' }}>
        <Col>
          <Row><h3>Token Info</h3></Row>
          <Row>Name: {token.name}</Row>
          <Row>Symbol: {token.symbol}</Row>
          <Row>Decimals: {token.decimals.toLocaleString()}</Row>
          <Row>Initial supply: {(token.initialSupply / 1E18).toLocaleString()} {token.symbol}</Row>
          <Row>Address: {token.tokenAddress}</Row>
          <Row>My tokens balance: {(token.myTokens / 1E18).toLocaleString()} {token.symbol}</Row>
        </Col>
      </Row>
    );
  }
}

export default connect(
  state => ({
    token: state.token,
  }),
  dispatch => bindActionCreators({
  }, dispatch)
)(Token);
