const Crowdsale = artifacts.require('./Crowdsale.sol');

module.exports = (deployer) => {
  /* constructor(uint _startPreICO, uint _startICO, uint _startPostICO, address _reserved) */
  deployer.deploy(Crowdsale, 1529452800, 1531267201, 1533859201, '0x2fD1CCe6BD2e563ABb3E3DbE6fA7334aF7332eB0');
};
