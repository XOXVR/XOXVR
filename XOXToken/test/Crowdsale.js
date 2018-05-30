/* global artifacts contract before beforeEach describe it expect web3 */
/* eslint func-names: 0 */
/* eslint no-await-in-loop: 0 */
/* eslint prefer-arrow-callback: 0 */
/* eslint no-unused-vars: 0 */

import { advanceBlock } from './helpers/advanceToBlock';
import ether from './helpers/ether';
import { increaseTimeTo, duration } from './helpers/increaseTime';
import latestTime from './helpers/latestTime';

const { BigNumber } = web3;
const should = require('chai')
  .use(require('chai-as-promised'))
  .use(require('chai-bignumber')(BigNumber))
  .should();

const Crowdsale = artifacts.require('./Crowdsale.sol');
const Token = artifacts.require('./XOXToken.sol');
const ZERO_ADDRESS = '0x00000000000000000000000000000000';

const Stage = {
  id: 0,
  totalSupply: 1,
  remainderTokens: 2,
  minInvestment: 3,
  end: 5,
  bonusLessTen: 6,
  bonusMoreTen: 7,
  start: 8,
};

contract('Crowdsale', function ([owner, ...accounts]) {
  before(async function () {
    // Advance to the next block to correctly read time in the solidity "now"
    // function interpreted by testrpc
    await advanceBlock();
  });

  beforeEach(async function () {
    this.startPreICO = latestTime() + duration.weeks(1);
    this.endPreICO = this.startPreICO + duration.days(28);
    this.startICO = this.endPreICO + duration.seconds(1);
    this.endICO = this.startICO + duration.days(28);
    this.startPostICO = this.endICO + duration.seconds(1);

    /* constructor(uint _startPreICO, uint _startICO, uint _startPostICO, address _reserved) */
    this.crowdsale = await Crowdsale.new(this.startPreICO, this.startICO, this.startPostICO, '0x0E484E54F7C9b7e97C3F8c1AB873910f6cBB4F80');
    const tokenAddress = await this.crowdsale.token();
    this.token = Token.at(tokenAddress);
  });

  describe('Sale tokens', function () {
    it('PreICO Stage (usually)', async function () {
      await increaseTimeTo(this.startPreICO);
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(10) }).should.be.fulfilled;
    });

    it('Investors', async function () {
      await increaseTimeTo(this.startPreICO);
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(1) });
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(4) });
      const numberInvestors = await this.crowdsale.returnNumberInvestors();
      numberInvestors.should.be.bignumber.equal(1);
      const balance = web3.eth.getBalance(this.crowdsale.address);
      balance.should.be.bignumber.equal(ether(5));
      const balanceCrowd = await this.token.balanceOf(this.crowdsale.address);
      balanceCrowd.div(1E18).toNumber().toFixed(6).should.be.bignumber.equal(899948867.786706);
    });

    it('PreICO Stage (bet less than 1 eth)', async function () {
      await increaseTimeTo(this.startPreICO);
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(10) });
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(0.5) }).should.be.rejected;
    });

    it('PreICO stage', async function () {
      const preICO = await this.crowdsale.preICO();
      preICO[Stage.remainderTokens].should.be.bignumber.equal(ether(6E6));
      preICO[Stage.minInvestment].should.be.bignumber.equal(ether(1));
      preICO[Stage.bonusLessTen].should.be.bignumber.equal(140);
      preICO[Stage.bonusMoreTen].should.be.bignumber.equal(160);
      preICO[Stage.start].should.be.bignumber.equal(this.startPreICO);
      preICO[Stage.end].should.be.bignumber.equal(this.endPreICO);
      preICO[Stage.id].should.be.bignumber.equal(0);
    });

    it('ICO stage', async function () {
      const ICO = await this.crowdsale.ICO();
      ICO[Stage.remainderTokens].should.be.bignumber.equal(ether(54E6));
      ICO[Stage.minInvestment].should.be.bignumber.equal(ether(0.5));
      ICO[Stage.bonusLessTen].should.be.bignumber.equal(0);
      ICO[Stage.bonusMoreTen].should.be.bignumber.equal(140);
      ICO[Stage.start].should.be.bignumber.equal(this.startICO);
      ICO[Stage.end].should.be.bignumber.equal(this.endICO);
      ICO[Stage.id].should.be.bignumber.equal(1);
    });

    it('postICO stage', async function () {
      const postICO = await this.crowdsale.postICO();
      postICO[Stage.remainderTokens].should.be.bignumber.equal(ether(840E6));
      postICO[Stage.start].should.be.bignumber.equal(this.startPostICO);
      postICO[Stage.id].should.be.bignumber.equal(2);
    });

    it('value tokens for 10 eth', async function () {
      await increaseTimeTo(this.startPreICO);
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(10) });
      const balanceBuyer = await this.token.balanceOf(accounts[0]);
      balanceBuyer.div(1E18).toFixed(9).should.be.bignumber.equal(116873.630387143);
    });

    it('multisig', async function () {
      await this.crowdsale.setMultisig(ZERO_ADDRESS).should.be.rejected;
      await this.crowdsale.setMultisig(accounts[5]);
      const multisig = await this.crowdsale.multisig();
      multisig.should.be.equal(accounts[5]);
    });

    it('backend', async function () {
      await increaseTimeTo(this.startPostICO);
      await this.crowdsale.getValueForBackend(ether(10), accounts[0]);
      const balanceBuyer = await this.token.balanceOf(accounts[0]);
      balanceBuyer.div(1E18).toFixed(10).should.be.bignumber.equal(73046.0189919649);
      const postICO = await this.crowdsale.postICO();
      const remainderTokensAfter = postICO[Stage.remainderTokens];
      const totalSold = balanceBuyer.div(1E18);
      const availableTokens = remainderTokensAfter.div(1E18);
      totalSold.add(availableTokens).toFixed(10).should.be.bignumber.equal(839999999.9999999999);
    });

    it('withDrawal', async function () {
      const multisig = await this.crowdsale.multisig();
      await increaseTimeTo(this.startPreICO);
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(50) });
      const balanceLessSoftCap = web3.eth.getBalance(multisig);
      balanceLessSoftCap.div(1E18).toNumber().toFixed().should.be.bignumber.equal(97);

      await this.crowdsale.saleTokens({ from: accounts[1], value: ether(50) });
      await this.crowdsale.saleTokens({ from: accounts[2], value: ether(50) });
      await this.crowdsale.saleTokens({ from: accounts[3], value: ether(50) });

      const balanceAfter = web3.eth.getBalance(multisig);
      const salesTokens = await this.crowdsale.salesTokens();
      balanceAfter.div(1E18).toNumber().toFixed().should.be.bignumber.equal(297);
      salesTokens.div(1E18).toFixed(5).should.be.bignumber.equal(2337472.60774);

      await increaseTimeTo(this.startICO);

      await this.crowdsale.saleTokens({ from: accounts[4], value: ether(50) });
      await this.crowdsale.saleTokens({ from: accounts[5], value: ether(50) });
      const balanceAfterTwo = web3.eth.getBalance(multisig);

      balanceAfterTwo.div(1E18).toNumber().toFixed().should.be.bignumber.equal(397);
    });

    it('buy tokens from different accounts', async function () {
      await increaseTimeTo(this.startPreICO);
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(10) });
      await this.crowdsale.saleTokens({ from: accounts[1], value: ether(10) });
      const balance = web3.eth.getBalance(this.crowdsale.address);
      balance.should.be.bignumber.equal(ether(20));

      const balanceBuyerOne = await this.token.balanceOf(accounts[0]);
      balanceBuyerOne.div(1E18).toFixed(9).should.be.bignumber.equal(116873.630387143);

      const balanceBuyerTwo = await this.token.balanceOf(accounts[1]);
      balanceBuyerTwo.div(1E18).toFixed(9).should.be.bignumber.equal(116873.630387143);
    });

    it('refund', async function () {
      await increaseTimeTo(this.startICO);
      await this.crowdsale.saleTokens({ from: accounts[8], value: ether(10) });
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(2) });
      await increaseTimeTo(this.startPostICO);
      await this.crowdsale.refund();
      const balanceBuyerOne = web3.eth.getBalance(accounts[8]);
      balanceBuyerOne.div(1E18).toNumber().toFixed().should.be.bignumber.equal(100);

      const balanceBuyerTwo = web3.eth.getBalance(accounts[0]);
      balanceBuyerTwo.div(1E18).toNumber().toFixed().should.be.bignumber.equal(5);
    });

    it('burn', async function () {
      await increaseTimeTo(this.startPostICO);
      await this.crowdsale.burnTokens();
      const balance = await this.token.balanceOf(this.crowdsale.address);
      balance.should.be.bignumber.equal(840E24);
    });

    it('bonus programm', async function () {
      await increaseTimeTo(this.startPreICO);
      await this.crowdsale.saleTokens({ from: accounts[8], value: ether(5) });
      let balanceBuyerOne5ETH = await this.token.balanceOf(accounts[8]);
      balanceBuyerOne5ETH = balanceBuyerOne5ETH.div(1E18).toFixed(10);
      balanceBuyerOne5ETH.should.be.bignumber.equal(51132.2132943754);
      await this.crowdsale.saleTokens({ from: accounts[8], value: ether(10) });
      let balanceSum = await this.token.balanceOf(accounts[8]);
      balanceSum = balanceSum.div(1E18).toFixed(9);
      balanceSum.should.be.bignumber.equal(168005.843681519);
      const preICO = await this.crowdsale.preICO();
      preICO[Stage.remainderTokens].div(1E18).toNumber().toFixed(4).should.be.bignumber.equal(5831994.1563);
      await increaseTimeTo(this.startICO);
      await this.crowdsale.saleTokens({ from: accounts[6], value: ether(5) });
      const balanceBuyerTwo5ETH = await this.token.balanceOf(accounts[6]);
      balanceBuyerTwo5ETH.div(1E18).toFixed(9).should.be.bignumber.equal(36523.009495982);
      await this.crowdsale.saleTokens({ from: accounts[6], value: ether(10) });
      const balanceSum2 = await this.token.balanceOf(accounts[6]);
      balanceSum2.div(1E18).toFixed(9).should.be.bignumber.equal(138787.436084733);
      const ICO = await this.crowdsale.ICO();
      ICO[Stage.remainderTokens].div(1E18).toNumber().toFixed(4).should.be.bignumber.equal(53861212.5639);
    });
  });
});
