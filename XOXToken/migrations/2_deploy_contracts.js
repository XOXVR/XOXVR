const Crowdsale = artifacts.require('./Crowdsale.sol');

/* constructor(uint _startPreICO, uint _startICO, uint _startPostICO, address _reserved) */
module.exports = (deployer) => {
  deployer.deploy(Crowdsale, 1529202069, 1531794069, 1534472469, "0x0E484E54F7C9b7e97C3F8c1AB873910f6cBB4F80");
};
