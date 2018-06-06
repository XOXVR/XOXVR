import contract from 'truffle-contract';

import CrowdsaleContract from '../../build/contracts/Crowdsale.json';
import web3Instance from './getWeb3';

class Crowdsale {
  async getInstance() {
    const crowdsale = await contract(CrowdsaleContract);
    await crowdsale.setProvider(web3Instance.currentProvider);
    this.instance = await crowdsale.deployed();
    return this.instance;
  }
}

const crowdsaleInstance = new Crowdsale();
export { crowdsaleInstance as default };
