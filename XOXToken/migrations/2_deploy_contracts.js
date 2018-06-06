const Crowdsale = artifacts.require('./Crowdsale.sol');

module.exports = (deployer) => {
  /* constructor(uint _startPreICO, uint _startICO, uint _startPostICO, address _reserved) */
  deployer.deploy(Crowdsale, 1529452800, 1531267201, 1533859201, '0x0E484E54F7C9b7e97C3F8c1AB873910f6cBB4F80');
};
