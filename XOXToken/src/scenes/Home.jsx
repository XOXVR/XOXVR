import React, { Component } from 'react';
import { Container } from 'reactstrap';
import web3Instance from '../services/getWeb3';

import Token from '../components/Token';
import Crowdsale from '../components/Crowdsale';

class Home extends Component {
  render() {
    if (!web3Instance) return (<div className="wait">Web3 not initialized. Wait pls...</div>);
    return (
      <div className="App">
        <header className="App-header">
          <h1 className="App-title">XOX token crowdsale</h1>
        </header>
        <Container style={{ paddingLeft: 30, paddingRight: 30 }}>
          <Token />
          <Crowdsale />
        </Container>
        <footer> © by Allmax 2018</footer>
      </div>
    );
  }
}

export default Home;
