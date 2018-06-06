import contract from 'truffle-contract';

import TokenContract from '../../build/contracts/XOXToken.json';
import web3Instance from './getWeb3';
import crowdsaleInstance from './Crowdsale';

class Token {
  async getInstance() {
    const tokenAddress = await crowdsaleInstance.instance.token();
    const XOXToken = await contract({ abi: TokenContract.abi });
    await XOXToken.setProvider(web3Instance.currentProvider);
    this.instance = await XOXToken.at(tokenAddress);
  }
}

const tokenInstance = new Token();
export { tokenInstance as default };
