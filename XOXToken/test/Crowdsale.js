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
  totalSupply: 0,
  remainderTokens: 1,
  end: 5,
  start: 9,
  minInvestment: 2,
  bonusLessTen: 6,
  bonusMoreTen: 7,
  id: 8
};

contract('Crowdsale', function ([owner, ...accounts]) {
  before(async function () {
    // Advance to the next block to correctly read time in the solidity "now"
    // function interpreted by testrpc
    await advanceBlock();
  });

  beforeEach(async function () {
    this.startPreICO = latestTime() + duration.weeks(1);
    this.endPreICO = this.startPreICO + duration.days(24);
    this.startICO = this.endPreICO + duration.seconds(1);
    this.endICO = this.startICO + duration.days(24);
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
      (balance.toNumber() / 1E18).should.be.bignumber.equal(5);
      let balanceCrowd = await this.token.balanceOf(this.crowdsale.address);
      balanceCrowd = balanceCrowd.toNumber();
      const newVar = balanceCrowd / 1E18;
      newVar.toFixed().should.be.bignumber.equal(899963477);
    });

    // it('PreICO Stage (buy all tokens)', async function () {
    //   await increaseTimeTo(this.startPreICO);
    //   const rate = await this.crowdsale.rate();
    //   let preICO = await this.crowdsale.preICO();
    //   const totalSupply = preICO[Stage.totalSupply];
    //   const expectedWei = totalSupply.div(1E18).times(rate);
    //   await this.crowdsale.saleTokens({ from: accounts[0], value: expectedWei });
    //   preICO = await this.crowdsale.preICO();
    //   preICO[Stage.remainderTokens].should.be.bignumber.equal(0);
    //   const investments = await this.crowdsale.investments(accounts[0]);
    //   investments.should.be.bignumber.equal(expectedWei);
    // });

    // it('PreICO buying all tokens, then ICO Stage (bet in ICO equally 0.4 eth)', async function () {
    //   await increaseTimeTo(this.startPreICO);
    //   const preICO = await this.crowdsale.preICO[Stage.remainderTokens];
    //   await this.crowdsale.saleTokens({ from: accounts[0], value: ether(10) });
    //   await this.crowdsale.saleTokens({ from: accounts[0], value: ether(9) });
    //   preICO.should.be.bignumber.equal(0);
    //   await increaseTimeTo(this.startICO);
    //   await this.crowdsale.saleTokens({ from: accounts[0], value: ether(0.4) }).should.be.rejected;
    // });

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
    });

    it('ICO stage', async function () {
      const ICO = await this.crowdsale.ICO();
      ICO[Stage.remainderTokens].should.be.bignumber.equal(ether(54E6));
      ICO[Stage.minInvestment].should.be.bignumber.equal(ether(0.5));
      ICO[Stage.bonusLessTen].should.be.bignumber.equal(0);
      ICO[Stage.bonusMoreTen].should.be.bignumber.equal(140);
      ICO[Stage.start].should.be.bignumber.equal(this.startICO);
      ICO[Stage.end].should.be.bignumber.equal(this.endICO);
    });

    it('postICO stage', async function () {
      const postICO = await this.crowdsale.postICO();
      postICO[Stage.remainderTokens].should.be.bignumber.equal(ether(840E6));
      postICO[Stage.start].should.be.bignumber.equal(this.startPostICO);
    });

    it('value tokens for 10 eth', async function () {
      await increaseTimeTo(this.startPreICO);
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(10) });
      let balanceBuyer = await this.token.balanceOf(accounts[0]);
      balanceBuyer = balanceBuyer.toNumber();
      const newVar = balanceBuyer / 1E18;
      newVar.toFixed().should.be.bignumber.equal(73046);
    });

    it('update stage', async function () {
      const preICO = await this.crowdsale.preICO();
      const ICO = await this.crowdsale.ICO();
      const postICO = await this.crowdsale.postICO();
      const idPreICO = preICO[Stage.id];
      const idICO = ICO[Stage.id];
      const idPostICO = postICO[Stage.id];
      await increaseTimeTo(this.startPreICO);
      idPreICO.should.be.bignumber.equal(0);
      await increaseTimeTo(this.startICO);
      idICO.should.be.bignumber.equal(1);
      await increaseTimeTo(this.startPostICO);
      idPostICO.should.be.bignumber.equal(2);
    });

    it('multisig', async function () {
      await this.crowdsale.setMultisig(ZERO_ADDRESS).should.be.rejected;
      await this.crowdsale.setMultisig(accounts[0]).should.be.fulfilled;
    });

    it('backend', async function () {
      const postICO = await this.crowdsale.postICO();
      const remainderTokensBefore = postICO[Stage.remainderTokens];
      await increaseTimeTo(this.startPostICO);
      let balanceBuyerBefore = await this.token.balanceOf(accounts[0]);
      balanceBuyerBefore = balanceBuyerBefore.toNumber();
      const tokensBefore = balanceBuyerBefore / 1E18;
      await this.crowdsale.getEthForBackend(ether(10), accounts[0]);
      let balanceBuyer = await this.token.balanceOf(accounts[0]);
      balanceBuyer = balanceBuyer.toNumber();
      const tokensAfter = balanceBuyer / 1E18;
      tokensBefore.should.be.bignumber.equal(0);
      tokensAfter.toFixed().should.be.bignumber.equal(73046);

      const postI = await this.crowdsale.postICO();
      const remainderTokensAfter = postI[Stage.remainderTokens];
      const tokensAvailable = (remainderTokensAfter.toNumber() / 1E18).toFixed();
      const tokensSold = tokensAfter.toFixed();
      (+tokensAvailable + +tokensSold).should.be.bignumber.equal(840000000);
    });

    it('withDrawal', async function () {
      const multisig = await this.crowdsale.multisig();
      const balanceBefore = web3.eth.getBalance(multisig);
      await increaseTimeTo(this.startPreICO);
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(50) });
      await this.crowdsale.saleTokens({ from: accounts[1], value: ether(50) });
      await this.crowdsale.saleTokens({ from: accounts[2], value: ether(50) });
      await this.crowdsale.saleTokens({ from: accounts[3], value: ether(50) });

      const balanceAfter = web3.eth.getBalance(multisig);
      const saleTokens = await this.crowdsale.salesTokens();
      const salesTokens = (saleTokens.toNumber() / 1E18).toFixed();

      const ethBalanceBefore = (balanceBefore.toNumber() / 1E18).toFixed();
      const ethBalanceAfter = (balanceAfter.toNumber() / 1E18).toFixed();

      ethBalanceBefore.should.be.bignumber.equal(97);
      ethBalanceAfter.should.be.bignumber.equal(297);
      salesTokens.should.be.bignumber.equal(1460920);

      await increaseTimeTo(this.startICO);

      await this.crowdsale.saleTokens({ from: accounts[4], value: ether(50) });
      await this.crowdsale.saleTokens({ from: accounts[5], value: ether(50) });
      await this.crowdsale.saleTokens({ from: accounts[6], value: ether(50) });
      await this.crowdsale.saleTokens({ from: accounts[7], value: ether(50) });

      const balanceAfterTwo = web3.eth.getBalance(multisig);
      const ethBalanceAfterTwo = (balanceAfterTwo.toNumber() / 1E18).toFixed();

      ethBalanceAfterTwo.should.be.bignumber.equal(497);

      let balanceBuyerTwo = await this.token.balanceOf(accounts[5]);
      balanceBuyerTwo = balanceBuyerTwo.toNumber();
      const tokensBalanceTwoBuyer = (balanceBuyerTwo / 1E18).toFixed();
      tokensBalanceTwoBuyer.should.be.bignumber.equal(365230);
    });

    it('buy tokens from different accounts', async function () {
      await increaseTimeTo(this.startPreICO);
      await this.crowdsale.saleTokens({ from: accounts[0], value: ether(10) });
      await this.crowdsale.saleTokens({ from: accounts[1], value: ether(10) });
      const balance = web3.eth.getBalance(this.crowdsale.address);
      (balance.toNumber() / 1E18).should.be.bignumber.equal(20);

      let balanceBuyerOne = await this.token.balanceOf(accounts[0]);
      balanceBuyerOne = balanceBuyerOne.toNumber();
      const tokensBalanceOneBuyer = (balanceBuyerOne / 1E18).toFixed();
      tokensBalanceOneBuyer.should.be.bignumber.equal(73046);

      let balanceBuyerTwo = await this.token.balanceOf(accounts[1]);
      balanceBuyerTwo = balanceBuyerTwo.toNumber();
      const tokensBalanceTwoBuyer = (balanceBuyerTwo / 1E18).toFixed();
      tokensBalanceTwoBuyer.should.be.bignumber.equal(73046);
    });

    // it('refund', async function () {
    //   await increaseTimeTo(this.startPostICO);
      
    // });
  });
});
